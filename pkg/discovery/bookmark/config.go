package bookmark

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sosedoff/pgweb/pkg/bookmarks"
)

var (
	// sslmodes is a list of allosed "sslmode" values
	sslmodes = []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
)

// bookmarksDir returns full path to bookmarks directory
func bookmarksDir() string {
	return filepath.Join(os.Getenv("HOME"), ".pgweb/bookmarks")
}

// isValidSSLMode returns true if given mode is valid
func isValidSSLMode(val string) bool {
	for _, mode := range sslmodes {
		if val == mode {
			return true
		}
	}
	return false
}

// readBookmark reads and parses the bookmark file
func readBookmark(path string) (*bookmarks.Bookmark, error) {
	b := &bookmarks.Bookmark{}
	if _, err := toml.DecodeFile(path, b); err != nil {
		return nil, err
	}

	// Fill in port value
	if b.Port == 0 {
		b.Port = 5432
	}

	// Check the sslmode value
	if b.Ssl == "" {
		b.Ssl = "disable"
	}
	if !isValidSSLMode(b.Ssl) {
		log.Printf("[bookmark] invalid ssl mode %q in %q", b.Ssl, path)
		return nil, nil
	}

	// Check SSH port
	if b.Ssh != nil && b.Ssh.Port == "" {
		b.Ssh.Port = "22"
	}

	return b, nil
}
