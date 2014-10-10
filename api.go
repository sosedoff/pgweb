package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Error struct {
	Message string `json:"error"`
}

func API_RunQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.JSON(400, errors.New("Query parameter is missing"))
		return
	}

	history = append(history, query)
	API_HandleQuery(query, c)
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
	res, err := dbClient.Query(fmt.Sprintf(SQL_TABLE_SCHEMA, c.Params.ByName("name")))

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format())
}

func API_History(c *gin.Context) {
	c.JSON(200, history)
}

func API_Info(c *gin.Context) {
	res, err := dbClient.Query(SQL_INFO)

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, res.Format()[0])
}

func API_HandleQuery(query string, c *gin.Context) {
	result, err := dbClient.Query(query)

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, result)
}
