package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const VERSION = "0.3.1"

var options struct {
	Version  bool   `short:"v" long:"version" description:"Print version"`
	Debug    bool   `short:"d" long:"debug" description:"Enable debugging mode" default:"false"`
	Url      string `long:"url" description:"Database connection string"`
	Host     string `long:"host" description:"Server hostname or IP"`
	Port     int    `long:"port" description:"Server port" default:"5432"`
	User     string `long:"user" description:"Database user"`
	Pass     string `long:"pass" description:"Password for user"`
	DbName   string `long:"db" description:"Database name"`
	Ssl      string `long:"ssl" description:"SSL option" default:"disable"`
	HttpHost string `long:"bind" description:"HTTP server host" default:"localhost"`
	HttpPort uint   `long:"listen" description:"HTTP server listen port" default:"8080"`
	AuthUser string `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass string `long:"auth-pass" description:"HTTP basic auth password"`
	SkipOpen bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
}

var dbClient *Client

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func getConnectionString() string {
	if options.Url != "" {
		url := options.Url

		if options.Ssl != "" && !strings.Contains(url, "sslmode") {
			url += fmt.Sprintf("?sslmode=%s", options.Ssl)
		}

		return url
	}

	str := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s",
		options.Host, options.Port,
		options.User, options.DbName,
		options.Ssl,
	)

	if options.Pass != "" {
		str += fmt.Sprintf(" password=%s", options.Pass)
	}

	return str
}

func connectionSettingsBlank() bool {
	return options.Host == "" &&
		options.User == "" &&
		options.DbName == "" &&
		options.Url == ""
}

func initClient() {
	if connectionSettingsBlank() {
		return
	}

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

	if options.Url == "" {
		options.Url = os.Getenv("DATABASE_URL")
	}

	if options.Version {
		fmt.Printf("pgweb v%s\n", VERSION)
		os.Exit(0)
	}
}

func startServer() {
	router := gin.Default()

	// Enable HTTP basic authentication only if both user and password are set
	if options.AuthUser != "" && options.AuthPass != "" {
		auth := map[string]string{options.AuthUser: options.AuthPass}
		router.Use(gin.BasicAuth(auth))
	}

	router.GET("/", API_Home)
	router.POST("/connect", API_Connect)
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
	go router.Run(fmt.Sprintf("%v:%v", options.HttpHost, options.HttpPort))
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func openPage() {
	url := fmt.Sprintf("http://%v:%v", options.HttpHost, options.HttpPort)
	fmt.Println("To view database open", url, "in browser")

	if options.SkipOpen {
		return
	}

	_, err := exec.Command("which", "open").Output()
	if err != nil {
		return
	}

	exec.Command("open", url).Output()
}

func main() {
	initOptions()
	initClient()

	if dbClient != nil {
		defer dbClient.db.Close()
	}

	if !options.Debug {
		gin.SetMode("release")
	}

	startServer()
	openPage()
	handleSignals()
}
