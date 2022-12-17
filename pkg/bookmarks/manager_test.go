package bookmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*

expBookmark := Bookmark{

		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Database: "mydatabase",
		SSLMode:  "disable",
	}
	b, err := GetBookmark("../../data", "bookmark")
	if assert.NoError(t, err) {
		assert.Equal(t, expBookmark, b)
	}

	_, err = GetBookmark("../../data", "bar")
	expErrStr := "couldn't find a bookmark with name bar"
	assert.Equal(t, expErrStr, err.Error())

	_, err = GetBookmark("foo", "bookmark")
	assert.Error(t, err)
*/

func Test_fileBasename(t *testing.T) {
	assert.Equal(t, "filename", fileBasename("filename.toml"))
	assert.Equal(t, "filename", fileBasename("path/filename.toml"))
	assert.Equal(t, "filename", fileBasename("~/long/path/filename.toml"))
	assert.Equal(t, "filename", fileBasename("filename"))
}

func Test_readBookmark(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		b, err := readBookmark("../../data/bookmark.toml")
		assert.NoError(t, err)
		assert.Equal(t, "bookmark", b.ID)
		assert.Equal(t, "localhost", b.Host)
		assert.Equal(t, 5432, b.Port)
		assert.Equal(t, "postgres", b.User)
		assert.Equal(t, "mydatabase", b.Database)
		assert.Equal(t, "disable", b.SSLMode)
		assert.Equal(t, "", b.Password)
		assert.Equal(t, "", b.URL)
	})

	t.Run("with url", func(t *testing.T) {
		b, err := readBookmark("../../data/bookmark_url.toml")
		assert.NoError(t, err)
		assert.Equal(t, "postgres://username:password@host:port/database?sslmode=disable", b.URL)
		assert.Equal(t, "", b.Host)
		assert.Equal(t, 5432, b.Port)
		assert.Equal(t, "", b.User)
		assert.Equal(t, "", b.Database)
		assert.Equal(t, "disable", b.SSLMode)
		assert.Equal(t, "", b.Password)
	})

	t.Run("invalid ssl", func(t *testing.T) {
		b, err := readBookmark("../../data/bookmark_invalid_ssl.toml")
		assert.NoError(t, err)
		assert.Equal(t, "disable", b.SSLMode)
	})

	t.Run("invalid file", func(t *testing.T) {
		_, err := readBookmark("foobar")
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("invalid syntax", func(t *testing.T) {
		_, err := readBookmark("../../data/invalid.toml")
		assert.Equal(t, "toml: line 1: expected '.' or '=', but got 'e' instead", err.Error())
	})
}
