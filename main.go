package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"os"
	"os/exec"
	"os/signal"
)

const VERSION = "0.2.0"

var options struct {
	Version  bool   `short:"v" description:"Print version"`
	Debug    bool   `short:"d" description:"Enable debugging mode" default:"false"`
	Url      string `long:"url" description:"Database connection string"`
	Host     string `long:"host" description:"Server hostname or IP" default:"localhost"`
	Port     int    `long:"port" description:"Server port" default:"5432"`
	User     string `long:"user" description:"Database user" default:"postgres"`
	DbName   string `long:"db" description:"Database name" default:"postgres"`
	Ssl      string `long:"ssl" description:"SSL option" default:"disable"`
	HttpPort uint   `long:"listen" description:"HTTP server listen port" default:"8080"`
}

var dbClient *Client

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
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

func initClient() {
	client, err := NewClient()
	if err != nil {
		exitWithMessage(err.Error())
	}

	fmt.Println("Connecting to server...")
	err = client.Test()
	if err != nil {
		exitWithMessage(err.Error())
	}

	fmt.Println("Checking tables...")
	tables, err := client.Tables()
	if err != nil {
		exitWithMessage(err.Error())
	}

	if len(tables) == 0 {
		exitWithMessage("Database does not have any tables")
	}

	dbClient = client
}

func initOptions() {
	_, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		os.Exit(1)
	}

	if options.Version {
		fmt.Printf("pgweb v%s\n", VERSION)
		os.Exit(1)
	}
}

func startServer() {
	router := gin.Default()

	router.GET("/", API_Home)
	router.GET("/databases", API_GetDatabases)
	router.GET("/info", API_Info)
	router.GET("/tables", API_GetTables)
	router.GET("/tables/:table", API_GetTable)
	router.GET("/tables/:table/info", API_GetTableInfo)
	router.GET("/tables/:table/indexes", API_TableIndexes)
	router.GET("/query", API_RunQuery)
	router.POST("/query", API_RunQuery)
	router.GET("/explain", API_ExplainQuery)
	router.POST("/explain", API_ExplainQuery)
	router.GET("/history", API_History)
	router.GET("/static/:type/:name", API_ServeAsset)

	fmt.Println("Starting server...")
	go router.Run(fmt.Sprintf(":%v", options.HttpPort))
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func openPage() {
	url := fmt.Sprintf("http://localhost:%v", options.HttpPort)
	fmt.Println("To view database open", url, "in browser")

	_, err := exec.Command("which", "open").Output()
	if err != nil {
		return
	}

	exec.Command("open", url).Output()
}

func main() {
	initOptions()
	initClient()

	defer dbClient.db.Close()

	if !options.Debug {
		gin.SetMode("release")
	}

	startServer()
	openPage()
	handleSignals()
}
