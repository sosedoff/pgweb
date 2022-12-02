package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuvistavie/securerandom"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/shared"
	"github.com/sosedoff/pgweb/static"
)

var (
	// DbClient represents the active database connection in a single-session mode
	DbClient *client.Client

	// DbSessions represents the mapping for client connections
	DbSessions *SessionManager
)

// DB returns a database connection from the client context
func DB(c *gin.Context) *client.Client {
	if command.Opts.Sessions {
		return DbSessions.Get(getSessionId(c.Request))
	}
	return DbClient
}

// setClient sets the database client connection for the sessions
func setClient(c *gin.Context, newClient *client.Client) error {
	currentClient := DB(c)
	if currentClient != nil {
		currentClient.Close()
	}

	if !command.Opts.Sessions {
		DbClient = newClient
		return nil
	}

	sid := getSessionId(c.Request)
	if sid == "" {
		return errSessionRequired
	}

	DbSessions.Add(sid, newClient)
	return nil
}

// GetHome renderes the home page
func GetHome(prefix string) http.Handler {
	if prefix != "" {
		prefix = "/" + prefix
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static.Static)))
}

func GetAssets(prefix string) http.Handler {
	if prefix != "" {
		prefix = "/" + prefix + "static/"
	} else {
		prefix = "/static/"
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static.Static)))
}

// GetSessions renders the number of active sessions
func GetSessions(c *gin.Context) {
	// In debug mode endpoint will return a lot of sensitive information
	// like full database connection string and all query history.
	if command.Opts.Debug {
		successResponse(c, DbSessions.Sessions())
		return
	}
	successResponse(c, gin.H{"sessions": DbSessions.Len()})
}

// ConnectWithBackend creates a new connection based on backend resource
func ConnectWithBackend(c *gin.Context) {
	// Setup a new backend client
	backend := Backend{
		Endpoint:    command.Opts.ConnectBackend,
		Token:       command.Opts.ConnectToken,
		PassHeaders: strings.Split(command.Opts.ConnectHeaders, ","),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Fetch connection credentials
	cred, err := backend.FetchCredential(ctx, c.Param("resource"), c)
	if err != nil {
		badRequest(c, err)
		return
	}

	// Make the new session
	sid, err := securerandom.Uuid()
	if err != nil {
		badRequest(c, err)
		return
	}
	c.Request.Header.Add("x-session-id", sid)

	// Connect to the database
	cl, err := client.NewFromUrl(cred.DatabaseURL, nil)
	if err != nil {
		badRequest(c, err)
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
		badRequest(c, err)
		return
	}

	redirectURI := fmt.Sprintf("/%s?session=%s", command.Opts.Prefix, sid)
	c.Redirect(302, redirectURI)
}

// Connect creates a new client connection
func Connect(c *gin.Context) {
	if command.Opts.LockSession {
		badRequest(c, errSessionLocked)
		return
	}

	var sshInfo *shared.SSHInfo
	url := c.Request.FormValue("url")

	if url == "" {
		badRequest(c, errURLRequired)
		return
	}

	opts := command.Options{URL: url}
	url, err := connection.FormatURL(opts)

	if err != nil {
		badRequest(c, err)
		return
	}

	if c.Request.FormValue("ssh") != "" {
		sshInfo = parseSshInfo(c)
	}

	cl, err := client.NewFromUrl(url, sshInfo)
	if err != nil {
		badRequest(c, err)
		return
	}

	err = cl.Test()
	if err != nil {
		badRequest(c, err)
		return
	}

	info, err := cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		badRequest(c, err)
		return
	}

	successResponse(c, info.Format()[0])
}

// SwitchDb perform database switch for the client connection
func SwitchDb(c *gin.Context) {
	if command.Opts.LockSession {
		badRequest(c, errSessionLocked)
		return
	}

	name := c.Request.URL.Query().Get("db")
	if name == "" {
		name = c.Request.FormValue("db")
	}
	if name == "" {
		badRequest(c, errDatabaseNameRequired)
		return
	}

	conn := DB(c)
	if conn == nil {
		badRequest(c, errNotConnected)
		return
	}

	// Do not allow switching databases for connections from third-party backends
	if conn.External {
		badRequest(c, errSessionLocked)
		return
	}

	currentURL, err := neturl.Parse(conn.ConnectionString)
	if err != nil {
		badRequest(c, errInvalidConnString)
		return
	}
	currentURL.Path = name

	cl, err := client.NewFromUrl(currentURL.String(), nil)
	if err != nil {
		badRequest(c, err)
		return
	}

	err = cl.Test()
	if err != nil {
		badRequest(c, err)
		return
	}

	info, err := cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		badRequest(c, err)
		return
	}

	conn.Close()

	successResponse(c, info.Format()[0])
}

// Disconnect closes the current database connection
func Disconnect(c *gin.Context) {
	if command.Opts.LockSession {
		badRequest(c, errSessionLocked)
		return
	}

	conn := DB(c)
	if conn == nil {
		badRequest(c, errNotConnected)
		return
	}

	err := conn.Close()
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, gin.H{"success": true})
}

// RunQuery executes the query
func RunQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		badRequest(c, errQueryRequired)
		return
	}

	HandleQuery(query, c)
}

// ExplainQuery renders query explain plan
func ExplainQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		badRequest(c, errQueryRequired)
		return
	}

	HandleQuery(fmt.Sprintf("EXPLAIN %s", query), c)
}

// AnalyzeQuery renders query explain plan and analyze profile
func AnalyzeQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		badRequest(c, errQueryRequired)
		return
	}

	HandleQuery(fmt.Sprintf("EXPLAIN ANALYZE %s", query), c)
}

// GetDatabases renders a list of all databases on the server
func GetDatabases(c *gin.Context) {
	if command.Opts.LockSession {
		serveResult(c, []string{}, nil)
		return
	}
	conn := DB(c)
	if conn.External {
		errorResponse(c, 403, errNotPermitted)
		return
	}

	names, err := DB(c).Databases()
	serveResult(c, names, err)
}

// GetObjects renders a list of database objects
func GetObjects(c *gin.Context) {
	result, err := DB(c).Objects()
	if err != nil {
		badRequest(c, err)
		return
	}
	successResponse(c, client.ObjectsFromResult(result))
}

// GetSchemas renders list of available schemas
func GetSchemas(c *gin.Context) {
	res, err := DB(c).Schemas()
	serveResult(c, res, err)
}

// GetTable renders table information
func GetTable(c *gin.Context) {
	var res *client.Result
	var err error

	if c.Request.FormValue("type") == client.ObjTypeMaterializedView {
		res, err = DB(c).MaterializedView(c.Params.ByName("table"))
	} else {
		res, err = DB(c).Table(c.Params.ByName("table"))
	}

	serveResult(c, res, err)
}

// GetTableRows renders table rows
func GetTableRows(c *gin.Context) {
	offset, err := parseIntFormValue(c, "offset", 0)
	if err != nil {
		badRequest(c, err)
		return
	}

	limit, err := parseIntFormValue(c, "limit", 100)
	if err != nil {
		badRequest(c, err)
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
		badRequest(c, err)
		return
	}

	countRes, err := DB(c).TableRowsCount(c.Params.ByName("table"), opts)
	if err != nil {
		badRequest(c, err)
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

	serveResult(c, res, err)
}

// GetTableInfo renders a selected table information
func GetTableInfo(c *gin.Context) {
	res, err := DB(c).TableInfo(c.Params.ByName("table"))
	if err == nil {
		successResponse(c, res.Format()[0])
	} else {
		badRequest(c, err)
	}
}

// GetHistory renders a list of recent queries
func GetHistory(c *gin.Context) {
	successResponse(c, DB(c).History)
}

// GetConnectionInfo renders information about current connection
func GetConnectionInfo(c *gin.Context) {
	res, err := DB(c).Info()

	if err != nil {
		badRequest(c, err)
		return
	}

	info := res.Format()[0]
	info["session_lock"] = command.Opts.LockSession

	successResponse(c, info)
}

// GetActivity renders a list of running queries
func GetActivity(c *gin.Context) {
	res, err := DB(c).Activity()
	serveResult(c, res, err)
}

// GetTableIndexes renders a list of database table indexes
func GetTableIndexes(c *gin.Context) {
	res, err := DB(c).TableIndexes(c.Params.ByName("table"))
	serveResult(c, res, err)
}

// GetTableConstraints renders a list of database constraints
func GetTableConstraints(c *gin.Context) {
	res, err := DB(c).TableConstraints(c.Params.ByName("table"))
	serveResult(c, res, err)
}

// HandleQuery runs the database query
func HandleQuery(query string, c *gin.Context) {
	rawQuery, err := base64.StdEncoding.DecodeString(desanitize64(query))
	if err == nil {
		query = string(rawQuery)
	}

	result, err := DB(c).Query(query)
	if err != nil {
		badRequest(c, err)
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

// GetBookmarks renders the list of available bookmarks
func GetBookmarks(c *gin.Context) {
	bookmarks, err := bookmarks.ReadAll(bookmarks.Path(command.Opts.BookmarksDir))
	serveResult(c, bookmarks, err)
}

// GetInfo renders the pgweb system information
func GetInfo(c *gin.Context) {
	successResponse(c, command.Info)
}

// DataExport performs database table export
func DataExport(c *gin.Context) {
	db := DB(c)

	info, err := db.Info()
	if err != nil {
		badRequest(c, err)
		return
	}

	dump := client.Dump{
		Table: strings.TrimSpace(c.Request.FormValue("table")),
	}

	// If pg_dump is not available the following code will not show an error in browser.
	// This is due to the headers being written first.
	if !dump.CanExport() {
		badRequest(c, errPgDumpNotFound)
		return
	}

	formattedInfo := info.Format()[0]
	filename := formattedInfo["current_database"].(string)
	if dump.Table != "" {
		filename = filename + "_" + dump.Table
	}
	reg := regexp.MustCompile(`[^._\\w]+`)
	cleanFilename := reg.ReplaceAllString(filename, "")

	c.Header(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s.sql.gz"`, cleanFilename),
	)

	err = dump.Export(db.ConnectionString, c.Writer)
	if err != nil {
		badRequest(c, err)
	}
}
