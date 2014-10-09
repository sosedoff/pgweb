package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	SQL_GET_TABLES = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_schema,table_name;"
)

type Client struct {
	db *sqlx.DB
}

type Result struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type Error struct {
	Message string `json:"error"`
}

var dbClient *Client

var options struct {
	Host   string `short:"h" long:"host" description:"Server hostname or IP" default:"localhost"`
	Port   int    `short:"p" long:"port" description:"Server port" default:"5432"`
	User   string `short:"u" long:"user" description:"Database user" default:"postgres"`
	DbName string `short:"d" long:"db" description:"Database name" default:"postgres"`
}

func getConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		options.Host, options.Port,
		options.User, options.DbName,
	)
}

func NewClient() (*Client, error) {
	fmt.Println(getConnectionString())
	db, err := sqlx.Open("postgres", getConnectionString())

	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func NewError(err error) Error {
	return Error{err.Error()}
}

func (client *Client) Query(query string) (*Result, error) {
	rows, err := client.db.Queryx(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	result := Result{
		Columns: cols,
	}

	for rows.Next() {
		obj, err := rows.SliceScan()

		for i, item := range obj {
			if item == nil {
				obj[i] = nil
			} else {
				t := reflect.TypeOf(item).Kind().String()

				if t == "slice" {
					obj[i] = string(item.([]byte))
				}
			}
		}

		if err == nil {
			result.Rows = append(result.Rows, obj)
		}
	}

	return &result, nil
}

func jsonError(err error) string {
	e := NewError(err)
	buff, _ := json.Marshal(e)
	return string(buff)
}

func API_RunQuery(c *gin.Context) {
	query := strings.TrimSpace(c.Request.FormValue("query"))

	if query == "" {
		c.String(400, jsonError(errors.New("Query parameter is missing")))
		return
	}

	API_HandleQuery(query, c)
}

func API_GetTables(c *gin.Context) {
	API_HandleQuery(SQL_GET_TABLES, c)
}

func API_HandleQuery(query string, c *gin.Context) {
	result, err := dbClient.Query(query)

	if err != nil {
		c.String(400, jsonError(err))
		return
	}

	buff, err := json.Marshal(result)

	if err != nil {
		c.String(400, jsonError(err))
		return
	}

	c.String(200, string(buff))
}

func initClient() {
	client, err := NewClient()

	if err != nil {
		log.Fatal(err)
	}

	dbClient = client
}

func initOptions() {
	_, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initOptions()
	initClient()
	defer dbClient.db.Close()

	router := gin.Default()
	router.GET("/tables", API_GetTables)
	router.GET("/select", API_RunQuery)
	router.POST("/select", API_RunQuery)

	fmt.Println("Starting server at 0.0.0.0:8080")
	router.Run(":8080")
}
