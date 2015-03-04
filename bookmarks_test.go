package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Invalid_Bookmark_Files(t *testing.T) {
	examples := map[string]string{
		"foobar":                       "open foobar: no such file or directory",
		"./fixtures/invalid.toml":      "Near line 1, key 'invalid encoding': Near line 2: Expected key separator '=', but got '\\n' instead.",
		"./fixtures/invalid_port.toml": "Type mismatch for 'main.Bookmark.Port': Expected string but found 'int64'.",
	}

	for path, message := range examples {
		_, err := readServerConfig(path)
		assert.Error(t, err)
		assert.Equal(t, message, err.Error())
	}
}

func Test_Bookmark(t *testing.T) {
	bookmark, err := readServerConfig("./fixtures/bookmark.toml")

	assert.Equal(t, nil, err)
	assert.Equal(t, "localhost", bookmark.Host)
	assert.Equal(t, "5432", bookmark.Port)
	assert.Equal(t, "postgres", bookmark.User)
	assert.Equal(t, "mydatabase", bookmark.Database)
	assert.Equal(t, "disable", bookmark.Ssl)
	assert.Equal(t, "", bookmark.Password)
	assert.Equal(t, "", bookmark.Url)
}

func Test_Bookmark_URL(t *testing.T) {
	bookmark, err := readServerConfig("./fixtures/bookmark_url.toml")

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://username:password@host:port/database?sslmode=disable", bookmark.Url)
	assert.Equal(t, "", bookmark.Host)
	assert.Equal(t, "", bookmark.Port)
	assert.Equal(t, "", bookmark.User)
	assert.Equal(t, "", bookmark.Database)
	assert.Equal(t, "", bookmark.Ssl)
	assert.Equal(t, "", bookmark.Password)
}

func Test_Bookmarks_Path(t *testing.T) {
	path := fmt.Sprintf("%s/.pgweb/bookmarks", os.Getenv("HOME"))
	assert.Equal(t, path, bookmarksPath())
}

func Test_Basename(t *testing.T) {
	assert.Equal(t, "filename", fileBasename("filename.toml"))
	assert.Equal(t, "filename", fileBasename("path/filename.toml"))
	assert.Equal(t, "filename", fileBasename("~/long/path/filename.toml"))
	assert.Equal(t, "filename", fileBasename("filename"))
}

func Test_ReadBookmarks_Invalid(t *testing.T) {
	bookmarks, err := readAllBookmarks("foobar")

	assert.Error(t, err)
	assert.Equal(t, "open foobar: no such file or directory", err.Error())
	assert.Equal(t, 0, len(bookmarks))
}

func Test_ReadBookmarks(t *testing.T) {
	bookmarks, err := readAllBookmarks("./fixtures")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(bookmarks))
}
