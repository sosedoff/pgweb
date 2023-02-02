package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryIsPermitted(t *testing.T) {
	examples := []struct {
		name     string
		query    Query
		args     []string
		expected bool
	}{
		{
			name:     "no input provided",
			query:    makeQuery("localhost", "someuser", "somedb", "default"),
			args:     makeArgs("", "", "", ""),
			expected: false,
		},
		{
			name:     "match on host",
			query:    makeQuery("localhost", "*", "*", "*"),
			args:     makeArgs("localhost", "user", "db", "default"),
			expected: true,
		},
		{
			name:     "match on full set",
			query:    makeQuery("localhost", "user", "database", "mode"),
			args:     makeArgs("localhost", "someuser", "somedb", "default"),
			expected: false,
		},
		{
			name:     "match on partial database",
			query:    makeQuery("localhost", "*", "myapp_*", "*"),
			args:     makeArgs("localhost", "user", "myapp_development", "default"),
			expected: true,
		},
		{
			name:     "match on full set but not mode",
			query:    makeQuery("localhost", "*", "*", "readonly"),
			args:     makeArgs("localhost", "user", "db", "default"),
			expected: false,
		},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			result := ex.query.IsPermitted(ex.args[0], ex.args[1], ex.args[2], ex.args[3])
			assert.Equal(t, ex.expected, result)
		})
	}
}

func makeArgs(vals ...string) []string {
	return vals
}

func makeQuery(host, user, database, mode string) Query {
	mustfield := func(input string) field {
		f, err := newField(input)
		if err != nil {
			panic(err)
		}
		return f
	}

	return Query{
		Meta: &Metadata{
			Host:     mustfield(host),
			User:     mustfield(user),
			Database: mustfield(database),
			Mode:     mustfield(mode),
		},
	}
}
