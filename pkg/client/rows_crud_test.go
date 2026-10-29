package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWritableColumns(t *testing.T) {
	meta := testMetadata()

	columns, err := writableColumns(meta, map[string]interface{}{
		"status":  "ACTIVE",
		"country": "FR",
	})

	require.NoError(t, err)
	// Sorted to keep generated statements deterministic
	assert.Equal(t, []string{"country", "status"}, columns)
}

func TestWritableColumnsRejectsUnknownColumn(t *testing.T) {
	_, err := writableColumns(testMetadata(), map[string]interface{}{"nope": 1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown column")
}

func TestWritableColumnsRejectsGeneratedColumn(t *testing.T) {
	meta := testMetadata()
	meta.Columns = append(meta.Columns, ColumnMeta{Name: "total", Updatable: true, Generated: true})

	_, err := writableColumns(meta, map[string]interface{}{"total": 1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not editable")
}

func TestPrimaryKeyCondition(t *testing.T) {
	sql, args, err := primaryKeyCondition(testMetadata(), map[string]interface{}{"id": 42}, 0)

	require.NoError(t, err)
	assert.Equal(t, `"id" = $1`, sql)
	assert.Equal(t, []interface{}{42}, args)
}

func TestPrimaryKeyConditionWithOffset(t *testing.T) {
	sql, args, err := primaryKeyCondition(testMetadata(), map[string]interface{}{"id": 42}, 3)

	require.NoError(t, err)
	assert.Equal(t, `"id" = $4`, sql)
	assert.Len(t, args, 1)
}

func TestPrimaryKeyConditionComposite(t *testing.T) {
	meta := &TableMetadata{
		Schema:     "public",
		Table:      "codes",
		Columns:    []ColumnMeta{{Name: "code", Updatable: true}, {Name: "tenant_id", Updatable: true}},
		PrimaryKey: []string{"code", "tenant_id"},
	}

	sql, args, err := primaryKeyCondition(meta, map[string]interface{}{"tenant_id": 1, "code": "A"}, 0)

	require.NoError(t, err)
	assert.Equal(t, `"code" = $1 AND "tenant_id" = $2`, sql)
	assert.Equal(t, []interface{}{"A", 1}, args)
}

func TestPrimaryKeyConditionRejectsPartialKey(t *testing.T) {
	meta := &TableMetadata{
		Schema:     "public",
		Table:      "codes",
		PrimaryKey: []string{"code", "tenant_id"},
	}

	_, _, err := primaryKeyCondition(meta, map[string]interface{}{"code": "A"}, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tenant_id is required")
}

func TestPrimaryKeyConditionRejectsExtraColumns(t *testing.T) {
	_, _, err := primaryKeyCondition(testMetadata(), map[string]interface{}{
		"id":     1,
		"status": "ACTIVE",
	}, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly 1 column")
}

func TestPrimaryKeyConditionRejectsNullValue(t *testing.T) {
	_, _, err := primaryKeyCondition(testMetadata(), map[string]interface{}{"id": nil}, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "can't be null")
}

func TestPrimaryKeyConditionWithoutPrimaryKey(t *testing.T) {
	meta := &TableMetadata{Schema: "public", Table: "logs"}

	_, _, err := primaryKeyCondition(meta, map[string]interface{}{"id": 1}, 0)

	assert.Equal(t, ErrNoPrimaryKey, err)
}

func TestQuotedName(t *testing.T) {
	meta := &TableMetadata{Schema: "my schema", Table: `we"ird`}

	assert.Equal(t, `"my schema"."we""ird"`, meta.quotedName())
}
