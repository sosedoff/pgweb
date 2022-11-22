package command

import (
	"errors"
	"os"
	"os/user"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Version                      bool   `short:"v" long:"version" description:"Print version"`
	Debug                        bool   `short:"d" long:"debug" description:"Enable debugging mode"`
	URL                          string `long:"url" description:"Database connection string"`
	Host                         string `long:"host" description:"Server hostname or IP" default:"localhost"`
	Port                         int    `long:"port" description:"Server port" default:"5432"`
	User                         string `long:"user" description:"Database user"`
	Pass                         string `long:"pass" description:"Password for user"`
	DbName                       string `long:"db" description:"Database name"`
	Ssl                          string `long:"ssl" description:"SSL mode"`
	SslRootCert                  string `long:"ssl-rootcert" description:"SSL certificate authority file"`
	SslCert                      string `long:"ssl-cert" description:"SSL client certificate file"`
	SslKey                       string `long:"ssl-key" description:"SSL client certificate key file"`
	HTTPHost                     string `long:"bind" description:"HTTP server host" default:"localhost"`
	HTTPPort                     uint   `long:"listen" description:"HTTP server listen port" default:"8081"`
	AuthUser                     string `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass                     string `long:"auth-pass" description:"HTTP basic auth password"`
	SkipOpen                     bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
	Sessions                     bool   `long:"sessions" description:"Enable multiple database sessions"`
	Prefix                       string `long:"prefix" description:"Add a url prefix"`
	ReadOnly                     bool   `long:"readonly" description:"Run database connection in readonly mode"`
	LockSession                  bool   `long:"lock-session" description:"Lock session to a single database connection"`
	Bookmark                     string `short:"b" long:"bookmark" description:"Bookmark to use for connection. Bookmark files are stored under $HOME/.pgweb/bookmarks/*.toml" default:""`
	BookmarksDir                 string `long:"bookmarks-dir" description:"Overrides default directory for bookmark files to search" default:""`
	DisablePrettyJSON            bool   `long:"no-pretty-json" description:"Disable JSON formatting feature for result export"`
	DisableSSH                   bool   `long:"no-ssh" description:"Disable database connections via SSH"`
	ConnectBackend               string `long:"connect-backend" description:"Enable database authentication through a third party backend"`
	ConnectToken                 string `long:"connect-token" description:"Authentication token for the third-party connect backend"`
	ConnectHeaders               string `long:"connect-headers" description:"List of headers to pass to the connect backend"`
	DisableConnectionIdleTimeout bool   `long:"no-idle-timeout" description:"Disable connection idle timeout"`
	ConnectionIdleTimeout        int    `long:"idle-timeout" description:"Set connection idle timeout in minutes" default:"180"`
	Cors                         bool   `long:"cors" description:"Enable Cross-Origin Resource Sharing (CORS)"`
	CorsOrigin                   string `long:"cors-origin" description:"Allowed CORS origins" default:"*"`
	BinaryCodec                  string `long:"binary-codec" description:"Codec for binary data serialization, one of 'none', 'hex', 'base58', 'base64'" default:"none"`
}

var Opts Options

// ParseOptions returns a new options struct from the input arguments
func ParseOptions(args []string) (Options, error) {
	var opts = Options{}

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return opts, err
	}

	if opts.URL == "" {
		opts.URL = os.Getenv("DATABASE_URL")
	}

	if opts.Prefix == "" {
		opts.Prefix = os.Getenv("URL_PREFIX")
	}

	// Handle edge case where pgweb is started with a default host `localhost` and no user.
	// When user is not set the `lib/pq` connection will fail and cause pgweb's termination.
	if (opts.Host == "localhost" || opts.Host == "127.0.0.1") && opts.User == "" {
		if username := getCurrentUser(); username != "" {
			opts.User = username
		} else {
			opts.Host = ""
		}
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
		opts.URL = ""
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

// SetDefaultOptions parses and assigns the options
func SetDefaultOptions() error {
	opts, err := ParseOptions([]string{})
	if err != nil {
		return err
	}
	Opts = opts
	return nil
}

// getCurrentUser returns a current user name
func getCurrentUser() string {
	u, _ := user.Current()
	if u != nil {
		return u.Username
	}
	return os.Getenv("USER")
}
