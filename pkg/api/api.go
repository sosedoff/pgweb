package api

import (
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	neturl "net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuvistavie/securerandom"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/shared"
	"github.com/sosedoff/pgweb/pkg/statements"
)

var (
	DbClient   *client.Client
	DbSessions = map[string]*client.Client{}
)

func DB(c *gin.Context) *client.Client {
	if command.Opts.Sessions {
		return DbSessions[getSessionId(c.Request)]
	} else {
		return DbClient
	}
}

func setClient(c *gin.Context, newClient *client.Client) error {
	currentClient := DB(c)
	if currentClient != nil {
		currentClient.Close()
	}

	if !command.Opts.Sessions {
		DbClient = newClient
		return nil
	}

	sessionId := getSessionId(c.Request)
	if sessionId == "" {
		return errors.New("Session ID is required")
	}

	DbSessions[sessionId] = newClient
	return nil
}

func GetHome(c *gin.Context) {
	serveStaticAsset("/index.html", c)
}

func GetAsset(c *gin.Context) {
	serveStaticAsset(c.Params.ByName("path"), c)
}

func GetSessions(c *gin.Context) {
	// In debug mode endpoint will return a lot of sensitive information
	// like full database connection string and all query history.
	if command.Opts.Debug {
		c.JSON(200, DbSessions)
		return
	}

	c.JSON(200, map[string]int{"sessions": len(DbSessions)})
}

func ConnectWithBackend(c *gin.Context) {
	// Setup a new backend client
	backend := Backend{
		Endpoint:    command.Opts.ConnectBackend,
		Token:       command.Opts.ConnectToken,
		PassHeaders: command.Opts.ConnectHeaders,
	}

	// Fetch connection credentials
	cred, err := backend.FetchCredential(c.Param("resource"), c)
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	// Make the new session
	sessionId, err := securerandom.Uuid()
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}
	c.Request.Header.Add("x-session-id", sessionId)

	// Connect to the database
	cl, err := client.NewFromUrl(cred.DatabaseUrl, nil)
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}
	cl.External = true

	// Finalize session seetup
	_, err = cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		c.JSON(400, Error{err.Error()})
		return
	}

	c.Redirect(301, fmt.Sprintf("/%s?session=%s", command.Opts.Prefix, sessionId))
}

func Connect(c *gin.Context) {
	if command.Opts.LockSession {
		c.JSON(400, Error{"Session is locked"})
		return
	}

	var sshInfo *shared.SSHInfo
	url := c.Request.FormValue("url")

	if url == "" {
		c.JSON(400, Error{"Url parameter is required"})
		return
	}

	opts := command.Options{Url: url}
	url, err := connection.FormatUrl(opts)

	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	if c.Request.FormValue("ssh") != "" {
		sshInfo = parseSshInfo(c)
	}

	cl, err := client.NewFromUrl(url, sshInfo)
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	err = cl.Test()
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	info, err := cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		c.JSON(400, Error{err.Error()})
		return
	}

	c.JSON(200, info.Format()[0])
}

func SwitchDb(c *gin.Context) {
	if command.Opts.LockSession {
		c.JSON(400, Error{"Session is locked"})
		return
	}

	name := c.Request.URL.Query().Get("db")
	if name == "" {
		name = c.Request.FormValue("db")
	}
	if name == "" {
		c.JSON(400, Error{"Database name is not provided"})
		return
	}

	conn := DB(c)
	if conn == nil {
		c.JSON(400, Error{"Not connected"})
		return
	}

	// Do not allow switching databases for connections from third-party backends
	if conn.External {
		c.JSON(400, Error{"Session is locked"})
		return
	}

	currentUrl, err := neturl.Parse(conn.ConnectionString)
	if err != nil {
		c.JSON(400, Error{"Unable to parse current connection string"})
		return
	}

	currentUrl.Path = name

	cl, err := client.NewFromUrl(currentUrl.String(), nil)
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	err = cl.Test()
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	info, err := cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		c.JSON(400, Error{err.Error()})
		return
	}

	conn.Close()

	c.JSON(200, info.Format()[0])
}

func Disconnect(c *gin.Context) {
	if command.Opts.LockSession {
		c.JSON(400, Error{"Session is locked"})
		return
	}

	conn := DB(c)

	if conn == nil {
		c.JSON(400, Error{"Not connected"})
		return
	}

	err := conn.Close()
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	c.JSON(200, map[string]bool{"success": true})
}

func GetDatabases(c *gin.Context) {
	conn := DB(c)
	if conn.External {
		c.JSON(403, Error{"Not permitted"})
		return
	}

	names, err := DB(c).Databases()
	serveResult(names, err, c)
}

func GetObjects(c *gin.Context) {
	result, err := DB(c).Objects()
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	objects := client.ObjectsFromResult(result)
	c.JSON(200, objects)
}

func RunQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, NewError(errors.New("Query parameter is missing")))
		return
	}

	HandleQuery(query, c)
}

func ExplainQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, NewError(errors.New("Query parameter is missing")))
		return
	}

	HandleQuery(fmt.Sprintf("EXPLAIN ANALYZE %s", query), c)
}

func GetSchemas(c *gin.Context) {
	res, err := DB(c).Schemas()
	serveResult(res, err, c)
}

func GetTable(c *gin.Context) {
	var res *client.Result
	var err error

	if c.Request.FormValue("type") == "materialized_view" {
		res, err = DB(c).MaterializedView(c.Params.ByName("table"))
	} else {
		res, err = DB(c).Table(c.Params.ByName("table"))
	}

	serveResult(res, err, c)
}

func GetTableRows(c *gin.Context) {
	offset, err := parseIntFormValue(c, "offset", 0)
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	limit, err := parseIntFormValue(c, "limit", 100)
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	opts := client.RowsOptions{
		Limit:      limit,
		Offset:     offset,
		SortColumn: c.Request.FormValue("sort_column"),
		SortOrder:  c.Request.FormValue("sort_order"),
		Where:      c.Request.FormValue("where"),
	}

	res, err := DB(c).TableRows(c.Params.ByName("table"), opts)
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	countRes, err := DB(c).TableRowsCount(c.Params.ByName("table"), opts)
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	numFetch := int64(opts.Limit)
	numOffset := int64(opts.Offset)
	numRows := countRes.Rows[0][0].(int64)
	numPages := numRows / numFetch

	if numPages*numFetch < numRows {
		numPages++
	}

	res.Pagination = &client.Pagination{
		Rows:    numRows,
		Page:    (numOffset / numFetch) + 1,
		Pages:   numPages,
		PerPage: numFetch,
	}

	serveResult(res, err, c)
}

func GetTableInfo(c *gin.Context) {
	res, err := DB(c).TableInfo(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format()[0])
}

func GetHistory(c *gin.Context) {
	c.JSON(200, DB(c).History)
}

func GetConnectionInfo(c *gin.Context) {
	res, err := DB(c).Info()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	info := res.Format()[0]
	info["session_lock"] = command.Opts.LockSession

	c.JSON(200, info)
}

func GetActivity(c *gin.Context) {
	res, err := DB(c).Activity()
	serveResult(res, err, c)
}

func GetTableIndexes(c *gin.Context) {
	res, err := DB(c).TableIndexes(c.Params.ByName("table"))
	serveResult(res, err, c)
}

func GetTableConstraints(c *gin.Context) {
	res, err := DB(c).TableConstraints(c.Params.ByName("table"))
	serveResult(res, err, c)
}

func HandleQuery(query string, c *gin.Context) {
	rawQuery, err := base64.StdEncoding.DecodeString(desanitize64(query))
	if err == nil {
		query = string(rawQuery)
	}

	result, err := DB(c).Query(query)
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	format := getQueryParam(c, "format")
	filename := getQueryParam(c, "filename")

	if filename == "" {
		filename = fmt.Sprintf("pgweb-%v.%v", time.Now().Unix(), format)
	}

	if format != "" {
		c.Writer.Header().Set("Content-disposition", "attachment;filename="+filename)
	}

	switch format {
	case "csv":
		c.Data(200, "text/csv", result.CSV())
	case "json":
		c.Data(200, "application/json", result.JSON())
	case "xml":
		c.XML(200, result)
	default:
		c.JSON(200, result)
	}
}

func GetBookmarks(c *gin.Context) {
	bookmarks, err := bookmarks.ReadAll(bookmarks.Path(command.Opts.BookmarksDir))
	serveResult(bookmarks, err, c)
}

func GetInfo(c *gin.Context) {
	info := map[string]string{
		"version":    command.VERSION,
		"git_sha":    command.GitCommit,
		"build_time": command.BuildTime,
	}

	c.JSON(200, info)
}

// Export database or table data
func DataExport(c *gin.Context) {
	db := DB(c)

	info, err := db.Info()
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	dump := client.Dump{
		Table: strings.TrimSpace(c.Request.FormValue("table")),
	}

	formattedInfo := info.Format()[0]
	filename := formattedInfo["current_database"].(string)
	if dump.Table != "" {
		filename = filename + "_" + dump.Table
	}

	attachment := fmt.Sprintf(`attachment; filename="%s.sql.gz"`, filename)
	c.Header("Content-Disposition", attachment)

	err = dump.Export(db.ConnectionString, c.Writer)
	if err != nil {
		c.JSON(400, Error{err.Error()})
	}
}

func DataImport(c *gin.Context) {
	c.Request.ParseMultipartForm(0)
	table := c.Request.FormValue("table")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Print(err)
		c.JSON(400, err.Error())
	}

	defer file.Close()
	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		log.Print(err)
		c.JSON(400, err.Error())
	}

	data, err := reader.ReadAll()
	if err != nil {
		log.Print(err)
		c.JSON(400, err.Error())
	}

	db := DB(c)
	createQuery := statements.CreateNewTableQuery(table, header)
	insertQuery := statements.GenerateBulkInsertQuery(table, header, len(data))

	_, err = db.NewTable(createQuery)
	if err != nil {
		log.Print(err)
		c.JSON(500, err.Error())
	}

	result, err := db.BulkInsert(insertQuery, statements.Flatten(data))
	if err != nil {
		log.Print(err)
		c.JSON(500, err.Error())
	}

	c.JSON(200, result)
}