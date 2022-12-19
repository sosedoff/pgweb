package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var hdir string
	if d, err := homedir.Dir(); err == nil {
		hdir = d
	}

	t.Run("defaults", func(t *testing.T) {
		opts, err := ParseOptions([]string{})
		assert.NoError(t, err)
		assert.Equal(t, false, opts.Sessions)
		assert.Equal(t, "", opts.Prefix)
		assert.Equal(t, "", opts.ConnectToken)
		assert.Equal(t, "", opts.ConnectHeaders)
		assert.Equal(t, false, opts.DisableSSH)
		assert.Equal(t, false, opts.DisablePrettyJSON)
		assert.Equal(t, false, opts.DisableConnectionIdleTimeout)
		assert.Equal(t, 180, opts.ConnectionIdleTimeout)
		assert.Equal(t, false, opts.Cors)
		assert.Equal(t, "*", opts.CorsOrigin)
		assert.Equal(t, "", opts.Passfile)
		assert.Equal(t, filepath.Join(hdir, ".pgweb/bookmarks"), opts.BookmarksDir)
	})

	t.Run("sessions", func(t *testing.T) {
		opts, err := ParseOptions([]string{"--sessions", "1"})
		assert.NoError(t, err)
		assert.Equal(t, true, opts.Sessions)
	})

	t.Run("url prefix", func(t *testing.T) {
		opts, err := ParseOptions([]string{"--prefix", "pgweb"})
		assert.NoError(t, err)
		assert.Equal(t, "pgweb/", opts.Prefix)

		opts, err = ParseOptions([]string{"--prefix", "pgweb/"})
		assert.NoError(t, err)
		assert.Equal(t, "pgweb/", opts.Prefix)
	})

	t.Run("connect backend", func(t *testing.T) {
		_, err := ParseOptions([]string{"--connect-backend", "test"})
		assert.EqualError(t, err, "--sessions flag must be set")

		_, err = ParseOptions([]string{"--connect-backend", "test", "--sessions"})
		assert.EqualError(t, err, "--connect-token flag must be set")

		_, err = ParseOptions([]string{"--connect-backend", "test", "--sessions", "--connect-token", "token"})
		assert.NoError(t, err)
	})

	t.Run("passfile", func(t *testing.T) {
		defer os.Unsetenv("PGPASSFILE")

		// File does not exist
		os.Setenv("PGPASSFILE", "/tmp/foo")
		opts, err := ParseOptions([]string{})
		assert.NoError(t, err)
		assert.Equal(t, "", opts.Passfile)

		// File exists and valid
		os.Setenv("PGPASSFILE", "../../data/passfile")
		opts, err = ParseOptions([]string{})
		assert.NoError(t, err)
		assert.Equal(t, "../../data/passfile", opts.Passfile)

		// Set via flag
		os.Unsetenv("PGPASSFILE")
		opts, err = ParseOptions([]string{"--passfile", "../../data/passfile"})
		assert.NoError(t, err)
		assert.Equal(t, "../../data/passfile", opts.Passfile)
	})
}
