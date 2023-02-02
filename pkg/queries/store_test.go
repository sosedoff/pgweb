//go:build !windows

package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreReadAll(t *testing.T) {
	t.Run("valid dir", func(t *testing.T) {
		queries, err := NewStore("../../data").ReadAll()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(queries))
	})

	t.Run("invalid dir", func(t *testing.T) {
		queries, err := NewStore("../../data2").ReadAll()
		assert.Equal(t, err.Error(), "queries directory does not exist")
		assert.Equal(t, 0, len(queries))
	})
}

func TestStoreRead(t *testing.T) {
	examples := []struct {
		id    string
		err   string
		check func(*testing.T, *Query)
	}{
		{id: "foo", err: "query file does not exist"},
		{id: "lc_no_meta"},
		{id: "lc_invalid_meta", err: `invalid "mode" field value: "foo"`},
		{
			id: "lc_example1",
			check: func(t *testing.T, q *Query) {
				assert.Equal(t, "lc_example1", q.ID)
				assert.Equal(t, "../../data/lc_example1.sql", q.Path)
				assert.Equal(t, "select 'foo'", q.Data)
				assert.Equal(t, "localhost", q.Meta.Host.String())
				assert.Equal(t, "*", q.Meta.User.String())
				assert.Equal(t, "*", q.Meta.Database.String())
			},
		},
		{
			id: "lc_example2",
			check: func(t *testing.T, q *Query) {
				assert.Equal(t, "lc_example2", q.ID)
				assert.Equal(t, "../../data/lc_example2.sql", q.Path)
				assert.Equal(t, "-- some comment\nselect 'foo'", q.Data)
				assert.Equal(t, "localhost", q.Meta.Host.String())
				assert.Equal(t, "foo", q.Meta.User.String())
				assert.Equal(t, "*", q.Meta.Database.String())
			},
		},
	}

	store := NewStore("../../data")

	for _, ex := range examples {
		t.Run(ex.id, func(t *testing.T) {
			query, err := store.Read(ex.id)
			if ex.err != "" || err != nil {
				assert.Equal(t, ex.err, err.Error())
			}
			if ex.check != nil {
				ex.check(t, query)
			}
		})
	}
}
