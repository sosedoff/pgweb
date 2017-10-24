package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"

	"github.com/sosedoff/pgweb/pkg/api"
	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/shared"
	"github.com/sosedoff/pgweb/pkg/util"
)

var options command.Options

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func initClientUsingBookmark(bookmarkPath, bookmarkName string) (*client.Client, error) {
	bookmark, err := bookmarks.GetBookmark(bookmarkPath, bookmarkName)
	if err != nil {
		return nil, err
	}

	opt := bookmark.ConvertToOptions()
	var connStr string

	if opt.Url != "" { // if the bookmark has url set, use it
		connStr = opt.Url
	} else {
		connStr, err = connection.BuildString(opt)
		if err != nil {
			return nil, fmt.Errorf("error building connection string: %v", err)
		}
	}

	var ssh *shared.SSHInfo
	if !bookmark.SSHInfoIsEmpty() {
		ssh = bookmark.Ssh
	}

	return client.NewFromUrl(connStr, ssh)
}

func initClient() {
	if connection.IsBlank(command.Opts) && options.Bookmark == "" {
		return
	}

	var cl *client.Client
	var err error

	if options.Bookmark != "" {
		cl, err = initClientUsingBookmark(bookmarks.Path(options.BookmarksDir), options.Bookmark)
	} else {
		cl, err = client.New()
	}

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

	if !command.Opts.Sessions {
		fmt.Printf("Server runs PostgreSQL v%s\n", cl.ServerVersion())
	}

	fmt.Println("Checking database objects...")
	_, err = cl.Objects()
	if err != nil {
		exitWithMessage(err.Error())
	}

	api.DbClient = cl
}

func initOptions() {
	err := command.ParseOptions()
	if err != nil {
		switch err.(type) {
		case *flags.Error:
			// no need to print error, flags package already does that
		default:
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	options = command.Opts

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

func Run() {
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

	// Start session cleanup worker
	if options.Sessions && !command.Opts.DisableConnectionIdleTimeout {
		go api.StartSessionCleanup()
	}

	startServer()
	openPage()
	handleSignals()
}
