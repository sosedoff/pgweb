package bookmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Invalid_Bookmark_Files(t *testing.T) {
	_, err := readServerConfig("foobar")
	assert.Error(t, err)

	_, err = readServerConfig("../../data/invalid.toml")
	assert.Error(t, err)
	assert.Equal(t, "Near line 1, key 'invalid encoding': Near line 2: Expected key separator '=', but got '\\n' instead.", err.Error())

	_, err = readServerConfig("../../data/invalid_port.toml")
	assert.Error(t, err)
	assert.Equal(t, "Type mismatch for 'bookmarks.Bookmark.Port': Expected string but found 'int64'.", err.Error())
}

func Test_Bookmark(t *testing.T) {
	bookmark, err := readServerConfig("../../data/bookmark.toml")

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
	bookmark, err := readServerConfig("../../data/bookmark_url.toml")

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
	assert.NotEqual(t, "/.pgweb/bookmarks", Path())
}

func Test_Basename(t *testing.T) {
	assert.Equal(t, "filename", fileBasename("filename.toml"))
	assert.Equal(t, "filename", fileBasename("path/filename.toml"))
	assert.Equal(t, "filename", fileBasename("~/long/path/filename.toml"))
	assert.Equal(t, "filename", fileBasename("filename"))
}

func Test_ReadBookmarks_Invalid(t *testing.T) {
	bookmarks, err := ReadAll("foobar")

	assert.Error(t, err)
	assert.Equal(t, 0, len(bookmarks))
}

func Test_ReadBookmarks(t *testing.T) {
	bookmarks, err := ReadAll("../../data")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(bookmarks))
}
