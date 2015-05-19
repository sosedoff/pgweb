package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
)

var DbClient *client.Client

func GetHome(c *gin.Context) {
	serveStaticAsset("/index.html", c)
}

func GetAsset(c *gin.Context) {
	serveStaticAsset(c.Params.ByName("path"), c)
}

func Connect(c *gin.Context) {
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
		if DbClient != nil {
			DbClient.Close()
		}

		DbClient = cl
	}

	c.JSON(200, info.Format()[0])
}

func GetDatabases(c *gin.Context) {
	names, err := DbClient.Databases()
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
	names, err := DbClient.Schemas()
	serveResult(names, err, c)
}

func GetTables(c *gin.Context) {
	names, err := DbClient.Tables()
	serveResult(names, err, c)
}

func GetTable(c *gin.Context) {
	res, err := DbClient.Table(c.Params.ByName("table"))
	serveResult(res, err, c)
}

func GetTableRows(c *gin.Context) {
	limit := 1000 // Number of rows to fetch
	limitVal := c.Request.FormValue("limit")

	if limitVal != "" {
		num, err := strconv.Atoi(limitVal)

		if err != nil {
			c.JSON(400, Error{"Invalid limit value"})
			return
		}

		if num <= 0 {
			c.JSON(400, Error{"Limit should be greater than 0"})
			return
		}

		limit = num
	}

	opts := client.RowsOptions{
		Limit:      limit,
		SortColumn: c.Request.FormValue("sort_column"),
		SortOrder:  c.Request.FormValue("sort_order"),
	}

	res, err := DbClient.TableRows(c.Params.ByName("table"), opts)
	serveResult(res, err, c)
}

func GetTableInfo(c *gin.Context) {
	res, err := DbClient.TableInfo(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format()[0])
}

func GetHistory(c *gin.Context) {
	c.JSON(200, DbClient.History)
}

func GetConnectionInfo(c *gin.Context) {
	res, err := DbClient.Info()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format()[0])
}

func GetActivity(c *gin.Context) {
	res, err := DbClient.Activity()
	serveResult(res, err, c)
}

func GetTableIndexes(c *gin.Context) {
	res, err := DbClient.TableIndexes(c.Params.ByName("table"))
	serveResult(res, err, c)
}

func HandleQuery(query string, c *gin.Context) {
	result, err := DbClient.Query(query)

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	q := c.Request.URL.Query()

	if len(q["format"]) > 0 && q["format"][0] == "csv" {
		filename := fmt.Sprintf("pgweb-%v.csv", time.Now().Unix())
		if len(q["filename"]) > 0 && q["filename"][0] != "" {
			filename = q["filename"][0]
		}

		c.Writer.Header().Set("Content-disposition", "attachment;filename="+filename)
		c.Data(200, "text/csv", result.CSV())
		return
	}

	c.JSON(200, result)
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
