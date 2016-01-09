package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
)

var DbClient *client.Client
var DbSessions = map[string]*client.Client{}

func DB(c *gin.Context) *client.Client {
	return DbSessions[getSessionId(c)]
}

func GetHome(c *gin.Context) {
	serveStaticAsset("/index.html", c)
}

func GetAsset(c *gin.Context) {
	serveStaticAsset(c.Params.ByName("path"), c)
}

func GetSessions(c *gin.Context) {
	c.JSON(200, DbSessions)
}

func Connect(c *gin.Context) {
	url := c.Request.FormValue("url")

	if url == "" {
		c.JSON(400, Error{"Url parameter is required"})
		return
	}

	sessionId := getSessionId(c)
	if sessionId == "" {
		c.JSON(400, Error{"Session ID is required"})
		return
	}

	opts := command.Options{Url: url}
	url, err := connection.FormatUrl(opts)

	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	cl, err := client.NewFromUrl(url)
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
		db := DbSessions[sessionId]
		if db != nil {
			db.Close()
		}

		DbSessions[sessionId] = cl
	}

	c.JSON(200, info.Format()[0])
}

func GetDatabases(c *gin.Context) {
	names, err := DB(c).Databases()
	serveResult(names, err, c)
}

func RunQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	HandleQuery(query, c)
}

func ExplainQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	HandleQuery(fmt.Sprintf("EXPLAIN ANALYZE %s", query), c)
}

func GetSchemas(c *gin.Context) {
	names, err := DB(c).Schemas()
	serveResult(names, err, c)
}

func GetTables(c *gin.Context) {
	names, err := DB(c).Tables()
	serveResult(names, err, c)
}

func GetTable(c *gin.Context) {
	res, err := DB(c).Table(c.Params.ByName("table"))
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

	c.JSON(200, res.Format()[0])
}

func GetSequences(c *gin.Context) {
	res, err := DB(c).Sequences()
	serveResult(res, err, c)
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
	rawQuery, err := base64.StdEncoding.DecodeString(query)
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
