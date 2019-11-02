package connection

import (
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"os/user"
	"strings"

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
		if opts.Ssl == "" {
			// Only modify sslmode for local connections
			if strings.Contains(uri.Host, "localhost") || strings.Contains(uri.Host, "127.0.0.1") {
				params["sslmode"] = "disable"
			}
		} else {
			params["sslmode"] = opts.Ssl
		}
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

	// Disable ssl for localhost connections, most users have it disabled
	if opts.Ssl == "" && (opts.Host == "localhost" || opts.Host == "127.0.0.1") {
		opts.Ssl = "disable"
	}

	query := neturl.Values{}
	if opts.Ssl != "" {
		query.Add("sslmode", opts.Ssl)
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
