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
	"github.com/sosedoff/pgweb/pkg/data"
)

var DbClient *client.Client

func GetHome(c *gin.Context) {
	data, err := data.Asset("static/index.html")

	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Data(200, "text/html; charset=utf-8", data)
}

func GetConnect(c *gin.Context) {
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

func GetGetDatabases(c *gin.Context) {
	names, err := DbClient.Databases()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, names)
}

func GetRunQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	GetHandleQuery(query, c)
}

func GetExplainQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	GetHandleQuery(fmt.Sprintf("EXPLAIN ANALYZE %s", query), c)
}

func GetGetSchemas(c *gin.Context) {
	names, err := DbClient.Schemas()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, names)
}

func GetGetTables(c *gin.Context) {
	names, err := DbClient.Tables()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, names)
}

func GetGetTable(c *gin.Context) {
	res, err := DbClient.Table(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res)
}

func GetGetTableRows(c *gin.Context) {
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

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res)
}

func GetGetTableInfo(c *gin.Context) {
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
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res)
}

func GetTableIndexes(c *gin.Context) {
	res, err := DbClient.TableIndexes(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res)
}

func GetHandleQuery(query string, c *gin.Context) {
	result, err := DbClient.Query(query)

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	q := c.Request.URL.Query()

	if len(q["format"]) > 0 && q["format"][0] == "csv" {
		filename := fmt.Sprintf("pgweb-%v.csv", time.Now().Unix())
		c.Writer.Header().Set("Content-disposition", "attachment;filename="+filename)
		c.Data(200, "text/csv", result.CSV())
		return
	}

	c.JSON(200, result)
}

func GetBookmarks(c *gin.Context) {
	bookmarks, err := bookmarks.ReadAll(bookmarks.Path())

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, bookmarks)
}

func GetAsset(c *gin.Context) {
	path := "static" + c.Params.ByName("path")
	data, err := data.Asset(path)

	if err != nil {
		c.String(400, err.Error())
		return
	}

	if len(data) == 0 {
		c.String(404, "Asset is empty")
		return
	}

	c.Data(200, assetContentType(path), data)
}
