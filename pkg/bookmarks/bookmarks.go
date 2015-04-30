package bookmarks

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type Bookmark struct {
	Url      string `json:"url"`      // Postgres connection URL
	Host     string `json:"host"`     // Server hostname
	Port     string `json:"port"`     // Server port
	User     string `json:"user"`     // Database user
	Password string `json:"password"` // User password
	Database string `json:"database"` // Database name
	Ssl      string `json:"ssl"`      // Connection SSL mode
}

func readServerConfig(path string) (Bookmark, error) {
	bookmark := Bookmark{}

	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return bookmark, err
	}

	_, err = toml.Decode(string(buff), &bookmark)
	return bookmark, err
}

func fileBasename(path string) string {
	filename := filepath.Base(path)
	return strings.Replace(filename, filepath.Ext(path), "", 1)
}

func Path() string {
	path, _ := homedir.Dir()
	return fmt.Sprintf("%s/.pgweb/bookmarks", path)
}

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

		fullPath := path + "/" + file.Name()
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
