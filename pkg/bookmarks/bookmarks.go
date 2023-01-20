package bookmarks

import (
	"os"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/shared"
)

// Bookmark contains information about bookmarked database connection
type Bookmark struct {
	ID          string          // ID generated from the filename
	URL         string          // Postgres connection URL
	Host        string          // Server hostname
	Port        int             // Server port
	User        string          // Database user
	UserVar     string          // Database user environment variable
	Password    string          // User password
	PasswordVar string          // User password environment variable
	Database    string          // Database name
	SSLMode     string          // Connection SSL mode
	SSH         *shared.SSHInfo // SSH tunnel config
}

// SSHInfoIsEmpty returns true if ssh configuration is not provided
func (b Bookmark) SSHInfoIsEmpty() bool {
	return b.SSH == nil || b.SSH.User == "" && b.SSH.Host == "" && b.SSH.Port == ""
}

// ConvertToOptions returns an options struct from connection details
func (b Bookmark) ConvertToOptions() command.Options {
	user := b.User
	if b.User == "" {
		user = os.Getenv(b.UserVar)
	}

	pass := b.Password
	if b.Password == "" {
		pass = os.Getenv(b.PasswordVar)
	}

	return command.Options{
		URL:     b.URL,
		Host:    b.Host,
		Port:    b.Port,
		User:    user,
		Pass:    pass,
		DbName:  b.Database,
		SSLMode: b.SSLMode,
	}
}
