package bookmarks

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/shared"
)

// Bookmark contains information about bookmarked database connection
type Bookmark struct {
	URL      string          `json:"url"`      // Postgres connection URL
	Host     string          `json:"host"`     // Server hostname
	Port     int             `json:"port"`     // Server port
	User     string          `json:"user"`     // Database user
	Password string          `json:"password"` // User password
	Database string          `json:"database"` // Database name
	Ssl      string          `json:"ssl"`      // Connection SSL mode
	SSH      *shared.SSHInfo `json:"ssh"`      // SSH tunnel config
}

// SSHInfoIsEmpty returns true if ssh configration is not provided
func (b Bookmark) SSHInfoIsEmpty() bool {
	return b.SSH == nil || b.SSH.User == "" && b.SSH.Host == "" && b.SSH.Port == ""
}

// ConvertToOptions returns an options struct from connection details
func (b Bookmark) ConvertToOptions() command.Options {
	return command.Options{
		URL:    b.URL,
		Host:   b.Host,
		Port:   b.Port,
		User:   b.User,
		Pass:   b.Password,
		DbName: b.Database,
		Ssl:    b.Ssl,
	}
}

func readServerConfig(path string) (Bookmark, error) {
	bookmark := Bookmark{}

	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return bookmark, err
	}

	_, err = toml.Decode(string(buff), &bookmark)

	if bookmark.Port == 0 {
		bookmark.Port = 5432
	}

	// List of all supported postgres modes
	modes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	valid := false

	for _, mode := range modes {
		if bookmark.Ssl == mode {
			valid = true
			break
		}
	}

	// Fall back to a default mode if mode is not set or invalid
	// Typical typo: ssl mode set to "disabled"
	if bookmark.Ssl == "" || !valid {
		bookmark.Ssl = "disable"
	}

	// Set default SSH port if it's not provided by user
	if bookmark.SSH != nil && bookmark.SSH.Port == "" {
		bookmark.SSH.Port = "22"
	}

	return bookmark, err
}

func fileBasename(path string) string {
	filename := filepath.Base(path)
	return strings.Replace(filename, filepath.Ext(path), "", 1)
}

// Path returns bookmarks storage path
func Path(overrideDir string) string {
	if overrideDir == "" {
		path, _ := homedir.Dir()
		return fmt.Sprintf("%s/.pgweb/bookmarks", path)
	}
	return overrideDir
}

// ReadAll returns all available bookmarks
func ReadAll(path string) (map[string]Bookmark, error) {
	results := map[string]Bookmark{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return results, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".toml" {
			continue
		}

		fullPath := filepath.Join(path, file.Name())
		key := fileBasename(file.Name())
		config, err := readServerConfig(fullPath)

		if err != nil {
			fmt.Printf("%s parse error: %s\n", fullPath, err)
			continue
		}

		results[key] = config
	}

	return results, nil
}

// GetBookmark reads an existing bookmark
func GetBookmark(bookmarkPath string, bookmarkName string) (Bookmark, error) {
	bookmarks, err := ReadAll(bookmarkPath)
	if err != nil {
		return Bookmark{}, err
	}

	bookmark, ok := bookmarks[bookmarkName]
	if !ok {
		return Bookmark{}, fmt.Errorf("couldn't find a bookmark with name %s", bookmarkName)
	}

	return bookmark, nil
}
