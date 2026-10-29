package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMetadata() *TableMetadata {
	return &TableMetadata{
		Schema: "public",
		Table:  "users",
		Columns: []ColumnMeta{
			{Name: "id", Updatable: true, PrimaryKey: true},
			{Name: "status", Updatable: true},
			{Name: "country", Updatable: true},
			{Name: "created_at", Updatable: true},
		},
		PrimaryKey: []string{"id"},
		Insertable: true,
		Editable:   true,
	}
}

func TestBuildFiltersEmpty(t *testing.T) {
	sql, args, err := buildFilters(nil, testMetadata(), 0)

	require.NoError(t, err)
	assert.Equal(t, "", sql)
	assert.Empty(t, args)
}

func TestBuildFiltersMultipleConditions(t *testing.T) {
	filters := []Filter{
		{Column: "status", Op: "equal", Value: "ACTIVE"},
		{Column: "country", Op: "equal", Value: "FR"},
		{Column: "created_at", Op: "greater", Value: "2026-01-01"},
	}

	sql, args, err := buildFilters(filters, testMetadata(), 0)

	require.NoError(t, err)
	assert.Equal(t, `"status" = $1 AND "country" = $2 AND "created_at" > $3`, sql)
	assert.Equal(t, []interface{}{"ACTIVE", "FR", "2026-01-01"}, args)
}

func TestBuildFiltersAllOperators(t *testing.T) {
	examples := []struct {
		op       string
		value    string
		expected string
		args     []interface{}
	}{
		{"equal", "a", `"status" = $1`, []interface{}{"a"}},
		{"not_equal", "a", `"status" != $1`, []interface{}{"a"}},
		{"greater", "1", `"status" > $1`, []interface{}{"1"}},
		{"greater_eq", "1", `"status" >= $1`, []interface{}{"1"}},
		{"less", "1", `"status" < $1`, []interface{}{"1"}},
		{"less_eq", "1", `"status" <= $1`, []interface{}{"1"}},
		{"like", "a%", `"status" LIKE $1`, []interface{}{"a%"}},
		{"ilike", "a%", `"status" ILIKE $1`, []interface{}{"a%"}},
		{"in", "a, b ,c", `"status" IN ($1, $2, $3)`, []interface{}{"a", "b", "c"}},
		{"null", "", `"status" IS NULL`, []interface{}{}},
		{"not_null", "", `"status" IS NOT NULL`, []interface{}{}},
	}

	for _, example := range examples {
		sql, args, err := buildFilters(
			[]Filter{{Column: "status", Op: example.op, Value: example.value}},
			testMetadata(),
			0,
		)

		require.NoError(t, err, example.op)
		assert.Equal(t, example.expected, sql, example.op)
		assert.Equal(t, example.args, args, example.op)
	}
}

func TestBuildFiltersPlaceholderOffset(t *testing.T) {
	filters := []Filter{
		{Column: "status", Op: "equal", Value: "ACTIVE"},
		{Column: "country", Op: "in", Value: "FR,DE"},
	}

	sql, args, err := buildFilters(filters, testMetadata(), 2)

	require.NoError(t, err)
	assert.Equal(t, `"status" = $3 AND "country" IN ($4, $5)`, sql)
	assert.Len(t, args, 3)
}

func TestBuildFiltersRejectsUnknownColumn(t *testing.T) {
	_, _, err := buildFilters(
		[]Filter{{Column: "password", Op: "equal", Value: "x"}},
		testMetadata(),
		0,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown filter column")
}

func TestBuildFiltersRejectsUnknownOperator(t *testing.T) {
	_, _, err := buildFilters(
		[]Filter{{Column: "status", Op: "; DROP TABLE users; --", Value: "x"}},
		testMetadata(),
		0,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported filter operator")
}

func TestBuildFiltersRejectsEmptyColumn(t *testing.T) {
	_, _, err := buildFilters([]Filter{{Column: "  ", Op: "equal"}}, testMetadata(), 0)
	require.Error(t, err)
}

func TestBuildFiltersRejectsEmptyInList(t *testing.T) {
	_, _, err := buildFilters(
		[]Filter{{Column: "status", Op: "in", Value: " , "}},
		testMetadata(),
		0,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one value")
}

// Injection attempts must end up as bind parameters, never as SQL.
func TestBuildFiltersValuesAreNeverInlined(t *testing.T) {
	payload := "' OR 1=1; DROP TABLE users; --"

	sql, args, err := buildFilters(
		[]Filter{{Column: "status", Op: "equal", Value: payload}},
		testMetadata(),
		0,
	)

	require.NoError(t, err)
	assert.Equal(t, `"status" = $1`, sql)
	assert.Equal(t, []interface{}{payload}, args)
	assert.NotContains(t, sql, "DROP")
}

func TestSplitListValue(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, splitListValue(" a , b ,c "))
	assert.Equal(t, []string{}, splitListValue("  "))
	assert.Equal(t, []string{"a"}, splitListValue("a"))
}

func TestFilterOperators(t *testing.T) {
	assert.Len(t, FilterOperators(), 11)
}
