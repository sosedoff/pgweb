package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"error"`
}

func assetContentType(name string) string {
	switch filepath.Ext(name) {
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".icon":
		return "image/x-icon"
	default:
		return "text/plain"
	}
}

func API_Home(c *gin.Context) {
	data, err := Asset("static/index.html")

	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Data(200, "text/html; charset=utf-8", data)
}

func API_Connect(c *gin.Context) {
	url := c.Request.FormValue("url")

	if url == "" {
		c.JSON(400, Error{"Url parameter is required"})
		return
	}

	client, err := NewClientFromUrl(url)
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	err = client.Test()
	if err != nil {
		c.JSON(400, Error{err.Error()})
		return
	}

	info, err := client.Info()

	if err == nil {
		if dbClient != nil {
			dbClient.db.Close()
		}

		dbClient = client
	}

	c.JSON(200, info.Format()[0])
}

func API_GetDatabases(c *gin.Context) {
	names, err := dbClient.Databases()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, names)
}

func API_RunQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	API_HandleQuery(query, c)
}

func API_ExplainQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	API_HandleQuery(fmt.Sprintf("EXPLAIN ANALYZE %s", query), c)
}

func API_GetTables(c *gin.Context) {
	names, err := dbClient.Tables()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, names)
}

func API_GetTable(c *gin.Context) {
	res, err := dbClient.Table(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res)
}

func API_GetTableInfo(c *gin.Context) {
	res, err := dbClient.TableInfo(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format()[0])
}

func API_History(c *gin.Context) {
	c.JSON(200, dbClient.history)
}

func API_Info(c *gin.Context) {
	if dbClient == nil {
		c.JSON(400, Error{"Not connected"})
		return
	}

	res, err := dbClient.Info()

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format()[0])
}

func API_TableIndexes(c *gin.Context) {
	res, err := dbClient.TableIndexes(c.Params.ByName("table"))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res)
}

func API_HandleQuery(query string, c *gin.Context) {
	result, err := dbClient.Query(query)

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	q := c.Request.URL.Query()

	if len(q["format"]) > 0 {
		if q["format"][0] == "csv" {
			c.Data(200, "text/csv", result.CSV())
			return
		}
	}

	c.JSON(200, result)
}

func API_ServeAsset(c *gin.Context) {
	file := fmt.Sprintf(
		"static/%s/%s",
		c.Params.ByName("type"),
		c.Params.ByName("name"),
	)

	data, err := Asset(file)

	if err != nil {
		c.String(400, err.Error())
		return
	}

	if len(data) == 0 {
		c.String(404, "Asset is empty")
		return
	}

	c.Data(200, assetContentType(file), data)
}
