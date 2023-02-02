package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFields(t *testing.T) {
	examples := []struct {
		input string
		err   error
		vals  map[string]string
	}{
		{input: "", err: nil, vals: nil},
		{input: "foobar", err: nil, vals: nil},
		{input: "-- no pgweb meta", err: nil, vals: nil},
		{
			input: `--pgweb: foo=bar`,
			err:   nil,
			vals:  map[string]string{},
		},
		{
			input: `--pgweb: host="localhost"`,
			err:   nil,
			vals:  map[string]string{"host": "localhost"},
		},
		{
			input: `--pgweb: host="*" user="admin" database  ="mydb"; mode = "readonly"`,
			err:   nil,
			vals: map[string]string{
				"host":     "*",
				"database": "mydb",
				"user":     "admin",
				"mode":     "readonly",
			},
		},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			fields, err := parseFields(ex.input)
			assert.Equal(t, ex.err, err)
			assert.Equal(t, ex.vals, fields)
		})
	}
}

func Test_parseMetadata(t *testing.T) {
	examples := []struct {
		input string
		err   string
		check func(meta *Metadata) bool
	}{
		{
			input: `--pgweb: `,
			err:   `host field must be set`,
		},
		{
			input: `--pgweb: hello="world"`,
			err:   `unknown key: "hello"`,
		},
		{
			input: `--pgweb: host="localhost" user="anyuser" database="anydb" mode="foo"`,
			err:   `invalid "mode" field value: "foo"`,
		},
		{
			input: "--pgweb2:",
			check: func(m *Metadata) bool {
				return m == nil
			},
		},
		{
			input: `--pgweb: host="localhost"`,
			check: func(m *Metadata) bool {
				return m.Host.value == "localhost" &&
					m.User.value == "*" &&
					m.Database.value == "*" &&
					m.Mode.value == "*" &&
					m.Timeout == nil
			},
		},
		{
			input: `--pgweb: host="localhost" user="anyuser" database="anydb" mode="*"`,
			check: func(m *Metadata) bool {
				return m.Host.value == "localhost" &&
					m.Host.re == nil &&
					m.User.value == "anyuser" &&
					m.Database.value == "anydb" &&
					m.Mode.value == "*" &&
					m.Timeout == nil
			},
		},
		{
			input: `--pgweb: host="localhost" timeout="foo"`,
			err:   `error initializing "timeout" field: strconv.Atoi: parsing "foo": invalid syntax`,
		},
		{
			input: `-- pgweb: host="local(host|dev)"`,
			check: func(m *Metadata) bool {
				return m.Host.value == "local(host|dev)" && m.Host.re != nil &&
					m.Host.matches("localhost") && m.Host.matches("localdev") &&
					!m.Host.matches("localfoo") && !m.Host.matches("superlocaldev")
			},
		},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			meta, err := parseMetadata(ex.input)
			if ex.err != "" {
				assert.Contains(t, err.Error(), ex.err)
			}
			if ex.check != nil {
				assert.Equal(t, true, ex.check(meta))
			}
		})
	}
}

func Test_sanitizeMetadata(t *testing.T) {
	examples := []struct {
		input  string
		output string
	}{
		{input: "", output: ""},
		{input: "foo", output: "foo"},
		{
			input: `
-- pgweb: metadata
query1
-- pgweb: more metadata

query2

`,
			output: "query1\nquery2",
		},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			assert.Equal(t, ex.output, sanitizeMetadata(ex.input))
		})
	}
}
