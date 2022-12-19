package bookmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagerList(t *testing.T) {
	examples := []struct {
		dir string
		num int
		err string
	}{
		{"../../data", 3, ""},
		{"../../data/bookmark.toml", 0, "is not a directory"},
		{"../../data2", 0, ""},
		{"", 0, ""},
	}

	for _, ex := range examples {
		t.Run(ex.dir, func(t *testing.T) {
			bookmarks, err := NewManager(ex.dir).List()
			if ex.err != "" {
				assert.Contains(t, err.Error(), ex.err)
			}
			assert.Len(t, bookmarks, ex.num)
		})
	}
}

func TestManagerListIDs(t *testing.T) {
	ids, err := NewManager("../../data").ListIDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"bookmark", "bookmark_invalid_ssl", "bookmark_url"}, ids)
}

func TestManagerGet(t *testing.T) {
	manager := NewManager("../../data")

	b, err := manager.Get("bookmark")
	assert.NoError(t, err)
	assert.Equal(t, "bookmark", b.ID)

	b, err = manager.Get("foo")
	assert.Equal(t, "bookmark foo not found", err.Error())
	assert.Nil(t, b)
}

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
		assert.Equal(t, "bookmark file foobar does not exist", err.Error())
	})

	t.Run("invalid syntax", func(t *testing.T) {
		_, err := readBookmark("../../data/invalid.toml")
		assert.Equal(t, "toml: line 1: expected '.' or '=', but got 'e' instead", err.Error())
	})
}
