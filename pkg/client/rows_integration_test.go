package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sosedoff/pgweb/pkg/command"
)

// integrationClient connects to the database pointed to by PGWEB_TEST_URL.
//
// These tests create and drop their own scratch tables, so they can run
// against any database the user has write access to. The suite is skipped
// when the variable is not set.
func integrationClient(t *testing.T) *Client {
	t.Helper()

	url := os.Getenv("PGWEB_TEST_URL")
	if url == "" {
		t.Skip("PGWEB_TEST_URL is not set, skipping integration tests")
	}

	if err := command.SetDefaultOptions(); err != nil {
		t.Fatal(err)
	}

	cl, err := NewFromUrl(url, nil)
	require.NoError(t, err)
	require.NoError(t, cl.Test())

	return cl
}

func execStatements(t *testing.T, cl *Client, statements ...string) {
	t.Helper()

	for _, stmt := range statements {
		_, err := cl.query(stmt)
		require.NoError(t, err, stmt)
	}
}

// setupCRUDTable creates a scratch table covering the data types supported
// by the row editor.
func setupCRUDTable(t *testing.T, cl *Client) string {
	t.Helper()

	table := "public.pgweb_crud_test"

	execStatements(t, cl,
		fmt.Sprintf("DROP TABLE IF EXISTS %s", table),
		fmt.Sprintf(`CREATE TABLE %s (
			id          bigserial PRIMARY KEY,
			name        text,
			code        varchar(16),
			quantity    integer,
			big_number  bigint,
			price       numeric(10,2),
			active      boolean,
			created_at  timestamp,
			due_on      date,
			external_id uuid
		)`, table),
	)

	t.Cleanup(func() {
		_, _ = cl.query(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	})

	return table
}

func TestIntegrationTableMetadata(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	meta, err := cl.TableMetadata(table)
	require.NoError(t, err)

	assert.Equal(t, "public", meta.Schema)
	assert.Equal(t, "pgweb_crud_test", meta.Table)
	assert.Equal(t, "r", meta.Kind)
	assert.Equal(t, []string{"id"}, meta.PrimaryKey)
	assert.True(t, meta.Insertable)
	assert.True(t, meta.Editable)
	assert.Len(t, meta.Columns, 10)

	id := meta.Column("id")
	require.NotNil(t, id)
	assert.True(t, id.PrimaryKey)
	assert.False(t, id.Nullable)
	require.NotNil(t, id.Default)

	name := meta.Column("name")
	require.NotNil(t, name)
	assert.Equal(t, "text", name.DataType)
	assert.True(t, name.Nullable)

	// Second lookup must be served from the cache
	again, err := cl.TableMetadata(table)
	require.NoError(t, err)
	assert.Same(t, meta, again)

	cl.ResetTableMetadata(table)
	fresh, err := cl.TableMetadata(table)
	require.NoError(t, err)
	assert.NotSame(t, meta, fresh)
}

func TestIntegrationTableMetadataUnknownTable(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	_, err := cl.TableMetadata("public.pgweb_does_not_exist")
	assert.Error(t, err)
}

func TestIntegrationRowCRUD(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	// --- INSERT ---
	res, err := cl.InsertRow(table, map[string]interface{}{
		"name":        "Widget",
		"code":        "W-1",
		"quantity":    10,
		"big_number":  9007199254740993,
		"price":       "19.99",
		"active":      true,
		"created_at":  "2026-01-02 03:04:05",
		"due_on":      "2026-02-03",
		"external_id": "11111111-2222-3333-4444-555555555555",
	})
	require.NoError(t, err)
	require.Len(t, res.Rows, 1)

	row := res.Format()[0]
	assert.Equal(t, "Widget", row["name"])
	assert.Equal(t, "W-1", row["code"])
	assert.EqualValues(t, 10, row["quantity"])
	assert.Equal(t, true, row["active"])

	id := row["id"]
	require.NotNil(t, id)

	// --- UPDATE: only the modified field is written ---
	res, err = cl.UpdateRow(table,
		map[string]interface{}{"id": id},
		map[string]interface{}{"name": "Updated widget"},
	)
	require.NoError(t, err)
	require.Len(t, res.Rows, 1)

	updated := res.Format()[0]
	assert.Equal(t, "Updated widget", updated["name"])
	assert.Equal(t, "W-1", updated["code"], "untouched columns must be preserved")

	// --- UPDATE: null values are supported ---
	res, err = cl.UpdateRow(table,
		map[string]interface{}{"id": id},
		map[string]interface{}{"name": nil},
	)
	require.NoError(t, err)
	assert.Nil(t, res.Format()[0]["name"])

	// --- DELETE ---
	res, err = cl.DeleteRow(table, map[string]interface{}{"id": id})
	require.NoError(t, err)
	assert.Len(t, res.Rows, 1)

	// Deleting the same row again must not silently succeed
	_, err = cl.DeleteRow(table, map[string]interface{}{"id": id})
	assert.Equal(t, ErrRowNotFound, err)
}

func TestIntegrationInsertDefaultValues(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	res, err := cl.InsertRow(table, map[string]interface{}{})
	require.NoError(t, err)
	require.Len(t, res.Rows, 1)
	assert.NotNil(t, res.Format()[0]["id"])
}

func TestIntegrationRowCRUDValidation(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	_, err := cl.InsertRow(table, map[string]interface{}{"nope": 1})
	assert.Error(t, err)

	_, err = cl.UpdateRow(table, map[string]interface{}{"id": 1}, map[string]interface{}{})
	assert.Equal(t, ErrNoValues, err)

	_, err = cl.UpdateRow(table, map[string]interface{}{}, map[string]interface{}{"name": "x"})
	assert.Error(t, err)

	_, err = cl.UpdateRow(table, map[string]interface{}{"id": 999999}, map[string]interface{}{"name": "x"})
	assert.Equal(t, ErrRowNotFound, err)
}

func TestIntegrationRowCRUDWithoutPrimaryKey(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := "public.pgweb_nopk_test"
	execStatements(t, cl,
		fmt.Sprintf("DROP TABLE IF EXISTS %s", table),
		fmt.Sprintf("CREATE TABLE %s (name text)", table),
	)
	defer cl.query(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)) //nolint

	meta, err := cl.TableMetadata(table)
	require.NoError(t, err)
	assert.True(t, meta.Insertable)
	assert.False(t, meta.Editable)

	// Inserts are fine without a primary key
	_, err = cl.InsertRow(table, map[string]interface{}{"name": "x"})
	require.NoError(t, err)

	// Updates and deletes are not
	_, err = cl.UpdateRow(table, map[string]interface{}{"name": "x"}, map[string]interface{}{"name": "y"})
	assert.Equal(t, ErrNoPrimaryKey, err)

	_, err = cl.DeleteRow(table, map[string]interface{}{"name": "x"})
	assert.Equal(t, ErrNoPrimaryKey, err)
}

func TestIntegrationRowCRUDCompositePrimaryKey(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := "public.pgweb_composite_test"
	execStatements(t, cl,
		fmt.Sprintf("DROP TABLE IF EXISTS %s", table),
		fmt.Sprintf(`CREATE TABLE %s (
			tenant_id integer NOT NULL,
			code      varchar(8) NOT NULL,
			label     text,
			PRIMARY KEY (tenant_id, code)
		)`, table),
	)
	defer cl.query(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)) //nolint

	meta, err := cl.TableMetadata(table)
	require.NoError(t, err)
	assert.Equal(t, []string{"tenant_id", "code"}, meta.PrimaryKey)
	assert.True(t, meta.Editable)

	_, err = cl.InsertRow(table, map[string]interface{}{
		"tenant_id": 1,
		"code":      "AA",
		"label":     "first",
	})
	require.NoError(t, err)

	key := map[string]interface{}{"tenant_id": 1, "code": "AA"}

	res, err := cl.UpdateRow(table, key, map[string]interface{}{"label": "second"})
	require.NoError(t, err)
	assert.Equal(t, "second", res.Format()[0]["label"])

	// A partial key must be rejected instead of matching multiple rows
	_, err = cl.UpdateRow(table, map[string]interface{}{"tenant_id": 1}, map[string]interface{}{"label": "x"})
	assert.Error(t, err)

	_, err = cl.DeleteRow(table, key)
	require.NoError(t, err)
}

func TestIntegrationRowCRUDReadOnly(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	cl.readonly = true
	defer func() { cl.readonly = false }()

	_, err := cl.InsertRow(table, map[string]interface{}{"name": "x"})
	assert.Equal(t, ErrReadOnlyMode, err)

	_, err = cl.UpdateRow(table, map[string]interface{}{"id": 1}, map[string]interface{}{"name": "x"})
	assert.Equal(t, ErrReadOnlyMode, err)

	_, err = cl.DeleteRow(table, map[string]interface{}{"id": 1})
	assert.Equal(t, ErrReadOnlyMode, err)
}

func TestIntegrationRowFilters(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	rows := []map[string]interface{}{
		{"name": "alpha", "code": "FR", "quantity": 1, "active": true},
		{"name": "beta", "code": "FR", "quantity": 5, "active": false},
		{"name": "gamma", "code": "DE", "quantity": 9, "active": true},
		{"name": nil, "code": "ES", "quantity": 3, "active": true},
	}
	for _, values := range rows {
		_, err := cl.InsertRow(table, values)
		require.NoError(t, err)
	}

	examples := []struct {
		name    string
		filters []Filter
		count   int
	}{
		{"single equality", []Filter{{Column: "code", Op: "equal", Value: "FR"}}, 2},
		{"multiple AND conditions", []Filter{
			{Column: "code", Op: "equal", Value: "FR"},
			{Column: "quantity", Op: "greater", Value: "2"},
		}, 1},
		{"not equal", []Filter{{Column: "code", Op: "not_equal", Value: "FR"}}, 2},
		{"range", []Filter{
			{Column: "quantity", Op: "greater_eq", Value: "3"},
			{Column: "quantity", Op: "less_eq", Value: "5"},
		}, 2},
		{"like", []Filter{{Column: "name", Op: "like", Value: "a%"}}, 1},
		{"ilike", []Filter{{Column: "name", Op: "ilike", Value: "A%"}}, 1},
		{"in", []Filter{{Column: "code", Op: "in", Value: "FR,DE"}}, 3},
		{"is null", []Filter{{Column: "name", Op: "null"}}, 1},
		{"is not null", []Filter{{Column: "name", Op: "not_null"}}, 3},
		{"boolean", []Filter{{Column: "active", Op: "equal", Value: "true"}}, 3},
		{"no match", []Filter{{Column: "code", Op: "equal", Value: "ZZ"}}, 0},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			opts := RowsOptions{Limit: 100, Filters: example.filters}

			res, err := cl.TableRows(table, opts)
			require.NoError(t, err)
			assert.Len(t, res.Rows, example.count)

			// Pagination must be computed against the filtered set
			countRes, err := cl.TableRowsCount(table, opts)
			require.NoError(t, err)
			assert.EqualValues(t, example.count, countRes.Rows[0][0])
		})
	}
}

func TestIntegrationRowFiltersRejectUnknownColumn(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	_, err := cl.TableRows(table, RowsOptions{
		Limit:   100,
		Filters: []Filter{{Column: "no_such_column", Op: "equal", Value: "x"}},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown filter column")
}

// Filter values must never be able to change the statement structure.
func TestIntegrationRowFiltersAreInjectionSafe(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	_, err := cl.InsertRow(table, map[string]interface{}{"name": "alpha"})
	require.NoError(t, err)

	res, err := cl.TableRows(table, RowsOptions{
		Limit:   100,
		Filters: []Filter{{Column: "name", Op: "equal", Value: "' OR 1=1 --"}},
	})

	require.NoError(t, err)
	assert.Len(t, res.Rows, 0, "injected payload must be treated as a literal value")

	// The table must still exist after a DROP attempt through a filter value
	_, err = cl.TableRows(table, RowsOptions{
		Limit:   100,
		Filters: []Filter{{Column: "name", Op: "equal", Value: "x'; DROP TABLE pgweb_crud_test; --"}},
	})
	require.NoError(t, err)

	meta, err := cl.TableMetadata(table)
	require.NoError(t, err)
	assert.NotEmpty(t, meta.Columns)
}

// Legacy free-form filters must keep working alongside structured filters.
func TestIntegrationRowFiltersLegacyWhere(t *testing.T) {
	cl := integrationClient(t)
	defer cl.Close()

	table := setupCRUDTable(t, cl)

	for _, code := range []string{"FR", "FR", "DE"} {
		_, err := cl.InsertRow(table, map[string]interface{}{"code": code, "quantity": 7})
		require.NoError(t, err)
	}

	res, err := cl.TableRows(table, RowsOptions{Limit: 100, Where: `"code" = 'FR'`})
	require.NoError(t, err)
	assert.Len(t, res.Rows, 2)

	// Combined with a structured filter
	res, err = cl.TableRows(table, RowsOptions{
		Limit:   100,
		Where:   `"code" = 'FR'`,
		Filters: []Filter{{Column: "quantity", Op: "equal", Value: "7"}},
	})
	require.NoError(t, err)
	assert.Len(t, res.Rows, 2)
}
