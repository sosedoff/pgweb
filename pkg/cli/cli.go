package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/api"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/discovery/aws"
	"github.com/sosedoff/pgweb/pkg/discovery/bookmark"
	"github.com/sosedoff/pgweb/pkg/discovery/heroku"
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

	var cl *client.Client
	var err error

	cl, err = client.New()
	if err != nil {
		exitWithMessage(err.Error())
	}

	if command.Opts.Debug {
		fmt.Println("Server connection string:", cl.ConnectionString)
	}

	fmt.Println("Connecting to server...")
	if err := cl.Test(); err != nil {
		msg := err.Error()

		// Check if we're trying to connect to the default database.
		if command.Opts.DbName == "" && command.Opts.Url == "" {
			// If database does not exist, allow user to connect from the UI.
			if strings.Contains(msg, "database") && strings.Contains(msg, "does not exist") {
				fmt.Println("Error:", msg)
				return
			}
			// Do not bail if local server is not running.
			if strings.Contains(msg, "connection refused") {
				fmt.Println("Error:", msg)
				return
			}
		}

		exitWithMessage(msg)
	}

	if !command.Opts.Sessions {
		fmt.Printf("Connected to %s\n", cl.ServerVersion())
	}

	fmt.Println("Checking database objects...")
	_, err = cl.Objects()
	if err != nil {
		exitWithMessage(err.Error())
	}

	api.DbClient = cl
}

func initOptions() {
	opts, err := command.ParseOptions(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	command.Opts = opts
	options = opts

	if options.Version {
		printVersion()
		os.Exit(0)
	}

	if options.ReadOnly {
		msg := `------------------------------------------------------
SECURITY WARNING: You are running pgweb in read-only mode.
This mode is designed for environments where users could potentially delete / change data.
For proper read-only access please follow postgresql role management documentation.
------------------------------------------------------`
		fmt.Println(msg)
	}

	printVersion()
}

func printVersion() {
	str := fmt.Sprintf("Pgweb v%s", command.Version)
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
			if strings.Contains(err.Error(), "address already in use") {
				openPage()
			}
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
	url := fmt.Sprintf("http://%v:%v/%s", options.HttpHost, options.HttpPort, options.Prefix)
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

func initProviders() {
	if command.Opts.DisableDiscovery {
		return
	}

	provider, err := bookmark.New(command.Opts)
	if err != nil {
		exitWithMessage(err.Error())
		return
	}
	api.RegisterProvider(provider)

	if command.Opts.Heroku {
		provider, err := heroku.New(command.Opts)
		if err != nil {
			exitWithMessage(err.Error())
			return
		}
		api.RegisterProvider(provider)
	}

	if command.Opts.AWS {
		provider, err := aws.New(command.Opts)
		if err != nil {
			exitWithMessage(err.Error())
			return
		}
		api.RegisterProvider(provider)
	}
}

func Run() {
	initOptions()
	initProviders()
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

	// Start session cleanup worker
	if options.Sessions && !command.Opts.DisableConnectionIdleTimeout {
		go api.StartSessionCleanup()
	}

	startServer()
	openPage()
	handleSignals()
}
