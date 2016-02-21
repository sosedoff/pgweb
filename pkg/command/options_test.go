package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Options(t *testing.T) {
	err := ParseOptions()

	assert.NoError(t, err)
	assert.Equal(t, false, Opts.Sessions)
	assert.Equal(t, "", Opts.Prefix)
}

func Test_SessionsOption(t *testing.T) {
	oldargs := os.Args
	defer func() { os.Args = oldargs }()

	os.Args = []string{"--sessions", "1"}
	assert.NoError(t, ParseOptions())
	assert.Equal(t, true, Opts.Sessions)
}

func Test_PrefixOption(t *testing.T) {
	oldargs := os.Args
	defer func() { os.Args = oldargs }()

	os.Args = []string{"--prefix", "pgweb"}
	assert.NoError(t, ParseOptions())
	assert.Equal(t, "pgweb/", Opts.Prefix)

	os.Args = []string{"--prefix", "pgweb/"}
	assert.NoError(t, ParseOptions())
	assert.Equal(t, "pgweb/", Opts.Prefix)
}
