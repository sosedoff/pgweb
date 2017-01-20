package connection

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type PassFileEntry struct {
	Hostname string
	Port     string
	Database string
	Username string
	Password string
}

type PassFile struct {
	Path    string
	Entries []PassFileEntry
}

// ReadPassFile reads a postgresl password file from user's home directory
// On OSX and Linux operating systems file is located at ~/.pgpass.
// On Windows its located at %APPDATA%\postgresql\pgpass.conf
func ReadPassFile(path string) (*PassFile, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// According to documentation, permissions on pgpass file should be 0600,
	// otherwise file should be ignored.
	if stat.Mode() != 0600 {
		return nil, fmt.Errorf("pgpass file has invalid permissions")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	passfile := PassFile{
		Path:    path,
		Entries: []PassFileEntry{},
	}

	for _, line := range strings.Split(string(data), "\n") {
		// Skip empty and comment lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")

		// We only expect 5 entries separated by colon
		if len(parts) != 5 {
			continue
		}

		passfile.Entries = append(passfile.Entries, PassFileEntry{
			Hostname: parts[0],
			Port:     parts[1],
			Database: parts[2],
			Username: parts[3],
			Password: parts[4],
		})
	}

	return &passfile, nil
}
