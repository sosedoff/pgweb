package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"os"
)

var options struct {
	Url    string `long:"url" description:"Database connection string"`
	Host   string `short:"h" long:"host" description:"Server hostname or IP" default:"localhost"`
	Port   int    `short:"p" long:"port" description:"Server port" default:"5432"`
	User   string `short:"u" long:"user" description:"Database user" default:"postgres"`
	DbName string `short:"d" long:"db" description:"Database name" default:"postgres"`
	Ssl    string `long:"ssl" description:"SSL option" default:"disable"`
	Static string `short:"s" description:"Path to static assets" default:"./static"`
}

var dbClient *Client
var history []string

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

func initClient() {
	client, err := NewClient()

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	_, err = client.Query(SQL_INFO)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dbClient = client
}

func initOptions() {
	_, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		fmt.Println("___")
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
	router.GET("/query", API_RunQuery)
	router.POST("/query", API_RunQuery)
	router.GET("/history", API_History)
	router.Static("/app", options.Static)

	fmt.Println("Starting server at 0.0.0.0:8080")
	router.Run(":8080")
}
