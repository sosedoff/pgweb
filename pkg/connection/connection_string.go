package connection

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/sosedoff/pgweb/pkg/command"
)

func currentUser() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.Username, nil
	}

	name := os.Getenv("USER")
	if name != "" {
		return name, nil
	}

	return "", errors.New("Unable to detect OS user")
}

func FormatUrl(opts command.Options) (string, error) {
	url := opts.Url

	// Make sure to only accept urls in a standard format
	if !strings.HasPrefix(url, "postgres://") && !strings.HasPrefix(url, "postgresql://") {
		return "", errors.New("Invalid URL. Valid format: postgres://user:password@host:port/db?sslmode=mode")
	}

	// Special handling for local connections
	if strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1") {
		if !strings.Contains(url, "?sslmode") {
			if opts.Ssl == "" {
				url += fmt.Sprintf("?sslmode=%s", "disable")
			} else {
				url += fmt.Sprintf("?sslmode=%s", opts.Ssl)
			}
		}
	}

	// Append sslmode parameter only if its defined as a flag and not present
	// in the connection string.
	if !strings.Contains(url, "?sslmode") && opts.Ssl != "" {
		url += fmt.Sprintf("?sslmode=%s", opts.Ssl)
	}

	return url, nil
}

func IsBlank(opts command.Options) bool {
	return opts.Host == "" && opts.User == "" && opts.DbName == "" && opts.Url == ""
}

func BuildString(opts command.Options) (string, error) {
	if opts.Url != "" {
		return FormatUrl(opts)
	}

	// Try to detect user from current OS user
	if opts.User == "" {
		u, err := currentUser()
		if err == nil {
			opts.User = u
		}
	}

	// Disable ssl for localhost connections, most users have it disabled
	if opts.Host == "localhost" || opts.Host == "127.0.0.1" {
		if opts.Ssl == "" {
			opts.Ssl = "disable"
		}
	}

	url := "postgres://"

	if opts.User != "" {
		url += opts.User
	}

	if opts.Pass != "" {
		url += fmt.Sprintf(":%s", opts.Pass)
	}

	url += fmt.Sprintf("@%s:%d", opts.Host, opts.Port)

	if opts.DbName != "" {
		url += fmt.Sprintf("/%s", opts.DbName)
	}

	if opts.Ssl != "" {
		url += fmt.Sprintf("?sslmode=%s", opts.Ssl)
	}

	return url, nil
}
