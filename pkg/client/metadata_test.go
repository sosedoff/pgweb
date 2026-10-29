package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuoteIdent(t *testing.T) {
	assert.Equal(t, `"users"`, quoteIdent("users"))
	assert.Equal(t, `"user name"`, quoteIdent("user name"))
	assert.Equal(t, `"we""ird"`, quoteIdent(`we"ird`))
	assert.Equal(t, `"a"";DROP TABLE x;--"`, quoteIdent(`a";DROP TABLE x;--`))
}

func TestBuildTableMetadata(t *testing.T) {
	res := &Result{
		Columns: []string{
			"column_name", "data_type", "udt_name", "is_nullable", "column_default",
			"ordinal_position", "is_identity", "is_generated", "is_updatable",
			"key_position", "relkind",
		},
		Rows: []Row{
			{"id", "integer", "int4", false, "nextval('t_id_seq')", int64(1), false, false, true, int64(1), "r"},
			{"name", "character varying", "varchar", true, nil, int64(2), false, false, true, int64(0), "r"},
			{"total", "numeric", "numeric", true, nil, int64(3), false, true, false, int64(0), "r"},
		},
	}

	meta, err := buildTableMetadata("public", "items", res)
	require.NoError(t, err)

	assert.Equal(t, "public", meta.Schema)
	assert.Equal(t, "items", meta.Table)
	assert.Equal(t, "r", meta.Kind)
	assert.Equal(t, []string{"id"}, meta.PrimaryKey)
	assert.True(t, meta.Insertable)
	assert.True(t, meta.Editable)
	assert.Len(t, meta.Columns, 3)

	id := meta.Column("id")
	require.NotNil(t, id)
	assert.True(t, id.PrimaryKey)
	assert.False(t, id.Nullable)
	assert.True(t, id.Writable())
	require.NotNil(t, id.Default)
	assert.Equal(t, "nextval('t_id_seq')", *id.Default)

	name := meta.Column("name")
	require.NotNil(t, name)
	assert.True(t, name.Nullable)
	assert.Nil(t, name.Default)

	// Generated columns are never writable
	total := meta.Column("total")
	require.NotNil(t, total)
	assert.True(t, total.Generated)
	assert.False(t, total.Writable())

	assert.Nil(t, meta.Column("missing"))
}

func TestBuildTableMetadataCompositeKeyOrder(t *testing.T) {
	res := &Result{
		Columns: []string{"column_name", "ordinal_position", "is_updatable", "key_position"},
		Rows: []Row{
			{"tenant_id", int64(1), true, int64(2)},
			{"code", int64(2), true, int64(1)},
			{"label", int64(3), true, int64(0)},
		},
	}

	meta, err := buildTableMetadata("public", "codes", res)
	require.NoError(t, err)

	// Primary key order must follow the index, not the column order
	assert.Equal(t, []string{"code", "tenant_id"}, meta.PrimaryKey)
	assert.True(t, meta.Editable)
}

func TestBuildTableMetadataWithoutPrimaryKey(t *testing.T) {
	res := &Result{
		Columns: []string{"column_name", "ordinal_position", "is_updatable", "key_position"},
		Rows: []Row{
			{"name", int64(1), true, int64(0)},
		},
	}

	meta, err := buildTableMetadata("public", "logs", res)
	require.NoError(t, err)

	assert.Empty(t, meta.PrimaryKey)
	assert.True(t, meta.Insertable)
	assert.False(t, meta.Editable)
}

func TestBuildTableMetadataNoRows(t *testing.T) {
	_, err := buildTableMetadata("public", "nope", &Result{Columns: []string{"column_name"}})
	assert.Error(t, err)

	_, err = buildTableMetadata("public", "nope", nil)
	assert.Error(t, err)
}

func TestBuildTableMetadataNonUpdatableRelation(t *testing.T) {
	res := &Result{
		Columns: []string{"column_name", "ordinal_position", "is_updatable", "key_position"},
		Rows: []Row{
			{"total", int64(1), false, int64(0)},
		},
	}

	meta, err := buildTableMetadata("public", "summary", res)
	require.NoError(t, err)

	assert.False(t, meta.Insertable)
	assert.False(t, meta.Editable)
}
