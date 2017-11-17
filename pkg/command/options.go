package command

import (
	"errors"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version                      bool    `short:"v" long:"version" description:"Print version"`
	Debug                        bool    `short:"d" long:"debug" description:"Enable debugging mode" default:"false"`
	Url                          string  `long:"url" description:"Database connection string"`
	Host                         string  `long:"host" description:"Server hostname or IP"`
	Port                         int     `long:"port" description:"Server port" default:"5432"`
	User                         string  `long:"user" description:"Database user"`
	Pass                         string  `long:"pass" description:"Password for user"`
	DbName                       string  `long:"db" description:"Database name"`
	Ssl                          string  `long:"ssl" description:"SSL option"`
	HttpHost                     string  `long:"bind" description:"HTTP server host" default:"localhost"`
	HttpPort                     uint    `long:"listen" description:"HTTP server listen port" default:"8081"`
	AuthUser                     string  `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass                     string  `long:"auth-pass" description:"HTTP basic auth password"`
	SkipOpen                     bool    `short:"s" long:"skip-open" description:"Skip browser open on start"`
	Sessions                     bool    `long:"sessions" description:"Enable multiple database sessions" default:"false"`
	Prefix                       string  `long:"prefix" description:"Add a url prefix"`
	ReadOnly                     bool    `long:"readonly" description:"Run database connection in readonly mode"`
	LockSession                  bool    `long:"lock-session" description:"Lock session to a single database connection" default:"false"`
	Bookmark                     string  `short:"b" long:"bookmark" description:"Bookmark to use for connection. Bookmark files are stored under $HOME/.pgweb/bookmarks/*.toml" default:""`
	BookmarksDir                 string  `long:"bookmarks-dir" description:"Overrides default directory for bookmark files to search" default:""`
	DisablePrettyJson            bool    `long:"no-pretty-json" description:"Disable JSON formatting feature for result export" default:"false"`
	DisableSSH                   bool    `long:"no-ssh" description:"Disable database connections via SSH" default:"false"`
	ConnectBackend               string  `long:"connect-backend" description:"Enable database authentication through a third party backend"`
	ConnectToken                 string  `long:"connect-token" description:"Authentication token for the third-party connect backend"`
	ConnectHeaders               string  `long:"connect-headers" description:"List of headers to pass to the connect backend"`
	DisableConnectionIdleTimeout bool    `long:"no-idle-timeout" description:"Disable connection idle timeout" default:"false"`
	ConnectionIdleTimeout        float64 `long:"idle-timeout" description:"Set connection idle timeout in minutes" default:"180"`
	Cors                         bool    `long:"cors" description:"Enable Cross-Origin Resource Sharing (CORS)" default:"false"`
	CorsOrigin                   string  `long:"cors-origin" description:"Allowed CORS origins" default:"*"`
}

var Opts Options

func ParseOptions() error {
	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		return err
	}

	if Opts.Url == "" {
		Opts.Url = os.Getenv("DATABASE_URL")
	}

	if os.Getenv("SESSIONS") != "" {
		Opts.Sessions = true
	}

	if os.Getenv("LOCK_SESSION") != "" {
		Opts.LockSession = true
		Opts.Sessions = false
	}

	if Opts.Prefix != "" && !strings.Contains(Opts.Prefix, "/") {
		Opts.Prefix = Opts.Prefix + "/"
	}

	if Opts.AuthUser == "" && os.Getenv("AUTH_USER") != "" {
		Opts.AuthUser = os.Getenv("AUTH_USER")
	}

	if Opts.AuthPass == "" && os.Getenv("AUTH_PASS") != "" {
		Opts.AuthPass = os.Getenv("AUTH_PASS")
	}

	if Opts.Bookmark != "" && Opts.Sessions {
		return errors.New("--bookmark is not allowed in multi-session mode")
	}

	if Opts.ConnectBackend != "" {
		if !Opts.Sessions {
			return errors.New("--sessions flag must be set")
		}
		if Opts.ConnectToken == "" {
			return errors.New("--connect-token flag must be set")
		}
	} else {
		if Opts.ConnectToken != "" || Opts.ConnectHeaders != "" {
			return errors.New("--connect-backend flag must be set")
		}
	}

	return nil
}
