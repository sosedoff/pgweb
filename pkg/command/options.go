package command

import (
	"errors"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version                      bool   `short:"v" long:"version" description:"Print version"`
	Debug                        bool   `short:"d" long:"debug" description:"Enable debugging mode"`
	Url                          string `long:"url" description:"Database connection string"`
	Host                         string `long:"host" description:"Server hostname or IP" default:"localhost"`
	Port                         int    `long:"port" description:"Server port" default:"5432"`
	User                         string `long:"user" description:"Database user"`
	Pass                         string `long:"pass" description:"Password for user"`
	DbName                       string `long:"db" description:"Database name"`
	Ssl                          string `long:"ssl" description:"SSL option"`
	HttpHost                     string `long:"bind" description:"HTTP server host" default:"localhost"`
	HttpPort                     uint   `long:"listen" description:"HTTP server listen port" default:"8081"`
	AuthUser                     string `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass                     string `long:"auth-pass" description:"HTTP basic auth password"`
	SkipOpen                     bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
	Sessions                     bool   `long:"sessions" description:"Enable multiple database sessions"`
	Prefix                       string `long:"prefix" description:"Add a url prefix"`
	ReadOnly                     bool   `long:"readonly" description:"Run database connection in readonly mode"`
	LockSession                  bool   `long:"lock-session" description:"Lock session to a single database connection"`
	Bookmark                     string `short:"b" long:"bookmark" description:"Bookmark to use for connection. Bookmark files are stored under $HOME/.pgweb/bookmarks/*.toml" default:""`
	BookmarksDir                 string `long:"bookmarks-dir" description:"Overrides default directory for bookmark files to search" default:""`
	DisablePrettyJson            bool   `long:"no-pretty-json" description:"Disable JSON formatting feature for result export"`
	DisableSSH                   bool   `long:"no-ssh" description:"Disable database connections via SSH"`
	ConnectBackend               string `long:"connect-backend" description:"Enable database authentication through a third party backend"`
	ConnectToken                 string `long:"connect-token" description:"Authentication token for the third-party connect backend"`
	ConnectHeaders               string `long:"connect-headers" description:"List of headers to pass to the connect backend"`
	DisableConnectionIdleTimeout bool   `long:"no-idle-timeout" description:"Disable connection idle timeout"`
	ConnectionIdleTimeout        int    `long:"idle-timeout" description:"Set connection idle timeout in minutes" default:"180"`
	Cors                         bool   `long:"cors" description:"Enable Cross-Origin Resource Sharing (CORS)"`
	CorsOrigin                   string `long:"cors-origin" description:"Allowed CORS origins" default:"*"`

	// Global flat to enable discovery feature
	Discovery bool `long:"discovery"`

	// Heroku Provider
	Heroku      bool   `long:"heroku"`
	HerokuToken string `long:"heroku-token"`

	// Amazon Web Services Provider
	AWS          bool   `long:"aws"`
	AWSAccessKey string `long:"aws-access-key"`
	AWSSecretKey string `long:"aws-secret-key"`
	AWSRegion    string `long:"aws-region"`
}

var Opts Options

func ParseOptions(args []string) (Options, error) {
	var opts = Options{}

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return opts, err
	}

	if opts.Url == "" {
		opts.Url = os.Getenv("DATABASE_URL")
	}

	if opts.Prefix == "" {
		opts.Prefix = os.Getenv("URL_PREFIX")
	}

	if os.Getenv("SESSIONS") != "" {
		opts.Sessions = true
	}

	if os.Getenv("LOCK_SESSION") != "" {
		opts.LockSession = true
		opts.Sessions = false
	}

	if opts.Sessions || opts.ConnectBackend != "" {
		opts.Bookmark = ""
		opts.Url = ""
		opts.Host = ""
		opts.User = ""
		opts.Pass = ""
		opts.DbName = ""
		opts.Ssl = ""
	}

	if opts.Prefix != "" && !strings.Contains(opts.Prefix, "/") {
		opts.Prefix = opts.Prefix + "/"
	}

	if opts.AuthUser == "" && os.Getenv("AUTH_USER") != "" {
		opts.AuthUser = os.Getenv("AUTH_USER")
	}

	if opts.AuthPass == "" && os.Getenv("AUTH_PASS") != "" {
		opts.AuthPass = os.Getenv("AUTH_PASS")
	}

	if opts.ConnectBackend != "" {
		if !opts.Sessions {
			return opts, errors.New("--sessions flag must be set")
		}
		if opts.ConnectToken == "" {
			return opts, errors.New("--connect-token flag must be set")
		}
	} else {
		if opts.ConnectToken != "" || opts.ConnectHeaders != "" {
			return opts, errors.New("--connect-backend flag must be set")
		}
	}

	return opts, nil
}

func SetDefaultOptions() error {
	opts, err := ParseOptions([]string{})
	if err != nil {
		return err
	}
	Opts = opts
	return nil
}
