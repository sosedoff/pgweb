package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/shared"
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
		if err != nil {
			cl.Close()
			c.JSON(400, Error{err.Error()})
			return
		}
	}

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
		c.Data(200, "applicaiton/json", result.JSON())
	case "xml":
		c.XML(200, result)
	default:
		c.JSON(200, result)
	}
}

func GetBookmarks(c *gin.Context) {
	bookmarks, err := bookmarks.ReadAll(bookmarks.Path())
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
