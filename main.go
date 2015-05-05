package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/sosedoff/pgweb/pkg/api"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/util"
)

var options command.Options

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func initClient() {
	if connection.IsBlank(command.Opts) {
		return
	}

	cl, err := client.New()
	if err != nil {
		exitWithMessage(err.Error())
	}

	if command.Opts.Debug {
		fmt.Println("Server connection string:", cl.ConnectionString)
	}

	fmt.Println("Connecting to server...")
	err = cl.Test()
	if err != nil {
		exitWithMessage(err.Error())
	}

	fmt.Println("Checking tables...")
	_, err = cl.Tables()
	if err != nil {
		exitWithMessage(err.Error())
	}

	api.DbClient = cl
}

func initOptions() {
	err := command.ParseOptions()
	if err != nil {
		os.Exit(1)
	}

	options = command.Opts

	if options.Version {
		printVersion()
		os.Exit(0)
	}

	printVersion()
}

func printVersion() {
	str := fmt.Sprintf("Pgweb v%s", command.VERSION)
	if command.GitCommit != "" {
		str += fmt.Sprintf(" (git: %s)", command.GitCommit)
	}

	fmt.Println(str)
}

func startServer() {
	router := gin.Default()

	// Enable HTTP basic authentication only if both user and password are set
	if options.AuthUser != "" && options.AuthPass != "" {
		auth := map[string]string{options.AuthUser: options.AuthPass}
		router.Use(gin.BasicAuth(auth))
	}

	api.SetupRoutes(router)

	fmt.Println("Starting server...")
	go func() {
		err := router.Run(fmt.Sprintf("%v:%v", options.HttpHost, options.HttpPort))
		if err != nil {
			fmt.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
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

	if api.DbClient != nil {
		defer api.DbClient.Close()
	}

	if !options.Debug {
		gin.SetMode("release")
	}

	// Print memory usage every 30 seconds with debug flag
	if options.Debug {
		util.StartProfiler()
	}

	startServer()
	openPage()
	handleSignals()
}
