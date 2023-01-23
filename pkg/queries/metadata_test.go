package queries

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseMetadata(t *testing.T) {
	type vals struct {
		host     string
		database string
		user     string
		mode     string
	}

	examples := []struct {
		input  string
		err    error
		values *vals
	}{
		{input: "", err: nil},
		{input: "foobar", err: nil},
		{input: "-- no pgweb meta", err: nil},
		{input: "--pgweb: foo=bar", err: errors.New(`invalid meta attribute: "foo"`)},
		{input: "--pgweb: mode=writeonly", err: errors.New(`invalid value for "mode" attribute: "writeonly"`)},
		{
			input: "--pgweb: host=localhost",
			err:   nil,
			values: &vals{
				host:     "localhost",
				database: "*",
				user:     "*",
				mode:     "*",
			},
		},
		{
			input: "--pgweb: host=*; user=admin; database  =mydb; mode = readonly",
			err:   nil,
			values: &vals{
				host:     "*",
				database: "mydb",
				user:     "admin",
				mode:     "readonly",
			},
		},
	}

	for _, ex := range examples {
		meta, err := parseMetadata(ex.input)
		assert.Equal(t, ex.err, err)
		if ex.values != nil && meta != nil {
			assert.Equal(t, ex.values.host, meta.host.input())
			assert.Equal(t, ex.values.database, meta.database.input())
			assert.Equal(t, ex.values.user, meta.user.input())
			assert.Equal(t, ex.values.mode, meta.mode.input())
		}
	}
}

func Test_matcher(t *testing.T) {
	examples := []struct {
		result   bool
		input    string
		host     string
		user     string
		database string
		mode     string
	}{
		{true, "-- pgweb: host=localhost", "localhost", "_", "_", "readonly"},
		{false, "-- pgweb: host=localhost", "anyhost", "_", "_", "readonly"},
		{false, "-- pgweb: host=localhost_(dev|test)", "localhost_foo", "_", "_", "readonly"},
		{false, "-- pgweb: host=localhost_(dev|test)", "localhost_development", "_", "_", "readonly"},
		{true, "-- pgweb: host=localhost_(dev|test)", "localhost_dev", "_", "_", "readonly"},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			meta, err := parseMetadata(ex.input)
			assert.NoError(t, err)
			assert.Equal(t, ex.result, meta.match(ex.host, ex.database, ex.user, ex.mode))
		})
	}
}
