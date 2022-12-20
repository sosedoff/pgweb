package connection

import (
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/jackc/pgpassfile"

	"github.com/sosedoff/pgweb/pkg/command"
)

// Common errors
var (
	errCantDetectUser   = errors.New("Could not detect default username")
	errInvalidURLFormat = errors.New("Invalid URL. Valid format: postgres://user:password@host:port/db?sslmode=mode")
)

// currentUser returns a current user name
func currentUser() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.Username, nil
	}

	name := os.Getenv("USER")
	if name != "" {
		return name, nil
	}

	return "", errCantDetectUser
}

// Check if connection url has a correct postgres prefix
func hasValidPrefix(str string) bool {
	return strings.HasPrefix(str, "postgres://") || strings.HasPrefix(str, "postgresql://")
}

// Extract all query vals and return as a map
func valsFromQuery(vals neturl.Values) map[string]string {
	result := map[string]string{}
	for k, v := range vals {
		result[strings.ToLower(k)] = v[0]
	}
	return result
}

// FormatURL reformats the existing connection string
func FormatURL(opts command.Options) (string, error) {
	url := opts.URL

	// Validate connection string prefix
	if !hasValidPrefix(url) {
		return "", errInvalidURLFormat
	}

	// Validate the URL
	uri, err := neturl.Parse(url)
	if err != nil {
		return "", errInvalidURLFormat
	}

	// Get query params
	params := valsFromQuery(uri.Query())

	// Determine if we need to specify sslmode if it's missing
	if params["sslmode"] == "" {
		if opts.SSLMode == "" {
			// Only modify sslmode for local connections
			if strings.Contains(uri.Host, "localhost") || strings.Contains(uri.Host, "127.0.0.1") {
				params["sslmode"] = "disable"
			}
		} else {
			params["sslmode"] = opts.SSLMode
		}
	}

	// When password is not provided, look it up from a .pgpass file
	if uri.User != nil {
		pass, _ := uri.User.Password()
		if pass == "" && opts.Passfile != "" {
			pass = lookupPassword(opts, uri)
			if pass != "" {
				uri.User = neturl.UserPassword(uri.User.Username(), pass)
			}
		}
	}

	// Configure default connect timeout
	if opts.OpenTimeout > 0 {
		params["connect_timeout"] = strconv.Itoa(opts.OpenTimeout)
	}

	// Rebuild query params
	query := neturl.Values{}
	for k, v := range params {
		query.Add(k, v)
	}
	uri.RawQuery = query.Encode()

	return uri.String(), nil
}

// IsBlank returns true if command options do not contain connection details
func IsBlank(opts command.Options) bool {
	return opts.Host == "" && opts.User == "" && opts.DbName == "" && opts.URL == ""
}

// BuildStringFromOptions returns a new connection string built from options
func BuildStringFromOptions(opts command.Options) (string, error) {
	query := neturl.Values{}

	// If connection string is provided we just use that
	if opts.URL != "" {
		return FormatURL(opts)
	}

	// Try to detect user from current OS user
	if opts.User == "" {
		u, err := currentUser()
		if err == nil {
			opts.User = u
		}
	}

	if opts.SSLMode != "" {
		query.Add("sslmode", opts.SSLMode)
	} else {
		if opts.Host == "localhost" || opts.Host == "127.0.0.1" {
			query.Add("sslmode", "disable")
		}
	}
	if opts.SSLCert != "" {
		query.Add("sslcert", opts.SSLCert)
	}
	if opts.SSLKey != "" {
		query.Add("sslkey", opts.SSLKey)
	}
	if opts.SSLRootCert != "" {
		query.Add("sslrootcert", opts.SSLRootCert)
	}

	// Grab password from .pgpass file if it's available
	if opts.Pass == "" && opts.Passfile != "" {
		opts.Pass = lookupPassword(opts, nil)
	}

	// Configure default connect timeout
	if opts.OpenTimeout > 0 {
		query.Add("connect_timeout", strconv.Itoa(opts.OpenTimeout))
	}

	url := neturl.URL{
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%v:%v", opts.Host, opts.Port),
		User:     neturl.UserPassword(opts.User, opts.Pass),
		Path:     fmt.Sprintf("/%s", opts.DbName),
		RawQuery: query.Encode(),
	}

	return url.String(), nil
}

func lookupPassword(opts command.Options, url *neturl.URL) string {
	if opts.Passfile == "" {
		return ""
	}

	passfile, err := pgpassfile.ReadPassfile(opts.Passfile)
	if err != nil {
		fmt.Println("[WARN] .pgpassfile", opts.Passfile, "is not readable")
		return ""
	}

	if url != nil {
		var dbName string
		fmt.Sscanf(url.Path, "/%s", &dbName)

		return passfile.FindPassword(
			url.Hostname(),
			url.Port(),
			dbName,
			url.User.Username(),
		)
	}

	return passfile.FindPassword(
		opts.Host,
		fmt.Sprintf("%d", opts.Port),
		opts.DbName,
		opts.User,
	)
}
