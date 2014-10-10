package main

import (
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
	SQL_INFO         = "SELECT version(), user, current_database(), inet_client_addr(), inet_client_port(), inet_server_addr(), inet_server_port()"
	SQL_TABLES       = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_schema,table_name;"
	SQL_TABLE_SCHEMA = "SELECT column_name, data_type, is_nullable, character_maximum_length, character_set_catalog, column_default FROM information_schema.columns where table_name = '%s';"
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
var history []string

var options struct {
	Url    string `long:"url" description:"Database connection string"`
	Host   string `short:"h" long:"host" description:"Server hostname or IP" default:"localhost"`
	Port   int    `short:"p" long:"port" description:"Server port" default:"5432"`
	User   string `short:"u" long:"user" description:"Database user" default:"postgres"`
	DbName string `short:"d" long:"db" description:"Database name" default:"postgres"`
	Ssl    string `long:"ssl" description:"SSL option" default:"disable"`
	Static string `short:"s" description:"Path to static assets" default:"./static"`
}

func formatResult(res *Result) []map[string]interface{} {
	var items []map[string]interface{}

	for _, row := range res.Rows {
		item := make(map[string]interface{})

		for i, c := range res.Columns {
			item[c] = row[i]
		}

		items = append(items, item)
	}

	return items
}

func getConnectionString() string {
	if options.Url != "" {
		return options.Url
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		options.Host, options.Port,
		options.User, options.DbName,
	)
}

func NewClient() (*Client, error) {
	db, err := sqlx.Open("postgres", getConnectionString())

	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func NewError(err error) Error {
	return Error{err.Error()}
}

func (client *Client) Tables() ([]string, error) {
	res, err := client.Query(SQL_TABLES)

	if err != nil {
		return nil, err
	}

	var tables []string

	for _, row := range res.Rows {
		tables = append(tables, row[0].(string))
	}

	return tables, nil
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

	c.JSON(200, formatResult(res))
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

	c.JSON(200, formatResult(res)[0])
}

func API_HandleQuery(query string, c *gin.Context) {
	result, err := dbClient.Query(query)

	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, result)
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

	router.GET("/info", API_Info)
	router.GET("/tables", API_GetTables)
	router.GET("/tables/:name", API_GetTable)
	router.GET("/select", API_RunQuery)
	router.POST("/select", API_RunQuery)
	router.GET("/history", API_History)

	router.Static("/app", options.Static)

	fmt.Println("Starting server at 0.0.0.0:8080")
	router.Run(":8080")
}
