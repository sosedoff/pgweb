package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	// Test default behavior
	opts, err := ParseOptions([]string{})
	assert.NoError(t, err)
	assert.Equal(t, false, opts.Sessions)
	assert.Equal(t, "", opts.Prefix)
	assert.Equal(t, "", opts.ConnectToken)
	assert.Equal(t, "", opts.ConnectHeaders)
	assert.Equal(t, false, opts.DisableSSH)
	assert.Equal(t, false, opts.DisablePrettyJson)
	assert.Equal(t, false, opts.DisableConnectionIdleTimeout)
	assert.Equal(t, 180, opts.ConnectionIdleTimeout)
	assert.Equal(t, false, opts.Cors)
	assert.Equal(t, "*", opts.CorsOrigin)

	// Test sessions
	opts, err = ParseOptions([]string{"--sessions", "1"})
	assert.NoError(t, err)
	assert.Equal(t, true, opts.Sessions)

	// Test url prefix
	opts, err = ParseOptions([]string{"--prefix", "pgweb"})
	assert.NoError(t, err)
	assert.Equal(t, "pgweb/", opts.Prefix)

	opts, err = ParseOptions([]string{"--prefix", "pgweb/"})
	assert.NoError(t, err)
	assert.Equal(t, "pgweb/", opts.Prefix)

	// Test connect backend options
	opts, err = ParseOptions([]string{"--connect-backend", "test"})
	assert.EqualError(t, err, "--sessions flag must be set")

	opts, err = ParseOptions([]string{"--connect-backend", "test", "--sessions"})
	assert.EqualError(t, err, "--connect-token flag must be set")

	opts, err = ParseOptions([]string{"--connect-backend", "test", "--sessions", "--connect-token", "token"})
	assert.NoError(t, err)
}
