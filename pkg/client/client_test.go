package client

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sosedoff/pgweb/pkg/command"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testClient     *Client
	testCommands   map[string]string
	serverHost     string
	serverPort     string
	serverUser     string
	serverPassword string
	serverDatabase string
)

func mapKeys(data map[string]*Objects) []string {
	result := []string{}
	for k := range data {
		result = append(result, k)
	}
	return result
}

func objectNames(data []Object) []string {
	names := make([]string, len(data))
	for i, obj := range data {
		names[i] = obj.Name
	}

	sort.Strings(names)
	return names
}

// assertMatches is a helper method to check if src slice contains any elements of expected slice
func assertMatches(t *testing.T, expected, src []string) {
	assert.NotEqual(t, 0, len(expected))
	assert.NotEqual(t, 0, len(src))

	for _, val := range expected {
		assert.Contains(t, src, val)
	}
}

func pgVersion() (int, int) {
	var major, minor int
	fmt.Sscanf(os.Getenv("PGVERSION"), "%d.%d", &major, &minor)
	return major, minor
}

func getVar(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

func initVars() {
	// We need to load default options to make sure all stuff works
	if err := command.SetDefaultOptions(); err != nil {
		log.Fatal(err)
	}

	serverHost = getVar("PGHOST", "localhost")
	serverPort = getVar("PGPORT", "5432")
	serverUser = getVar("PGUSER", "postgres")
	serverPassword = getVar("PGPASSWORD", "postgres")
	serverDatabase = getVar("PGDATABASE", "booktown")
}

func setupCommands() {
	testCommands = map[string]string{
		"createdb": "createdb",
		"psql":     "psql",
		"dropdb":   "dropdb",
	}

	if onWindows() {
		for k, v := range testCommands {
			testCommands[k] = v + ".exe"
		}
	}
}

func onWindows() bool {
	return runtime.GOOS == "windows"
}

func setup() {
	// No pretty JSON for tests
	command.Opts.DisablePrettyJSON = true

	out, err := exec.Command(
		testCommands["createdb"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Database creation failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	out, err = exec.Command(
		testCommands["psql"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		"-f", "../../data/booktown.sql",
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Database import failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func setupClient() {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	testClient, _ = NewFromUrl(url, nil)
}

func teardownClient() {
	if testClient != nil {
		testClient.Close()
	}
}

func teardown(t *testing.T, allowFail bool) {
	output, err := exec.Command(
		testCommands["dropdb"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		serverDatabase,
	).CombinedOutput()

	if err != nil && strings.Contains(err.Error(), "does not exist") {
		t.Log("Teardown error:", err)
		t.Logf("%s\n", output)

		if !allowFail {
			assert.NoError(t, err)
		}
	}
}

func testNewClientFromURL(t *testing.T) {
	t.Run("postgres prefix", func(t *testing.T) {
		url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
		client, err := NewFromUrl(url, nil)

		assert.Equal(t, nil, err)
		assert.Equal(t, url, client.ConnectionString)
		assert.NoError(t, client.Close())
	})

	t.Run("postgresql prefix", func(t *testing.T) {
		url := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
		client, err := NewFromUrl(url, nil)

		assert.Equal(t, nil, err)
		assert.Equal(t, url, client.ConnectionString)
		assert.NoError(t, client.Close())
	})
}

func testClientIdleTime(t *testing.T) {
	examples := map[time.Time]bool{
		time.Now():                         false, // Current time
		time.Now().Add(time.Minute * -30):  false, // 30 minutes ago
		time.Now().Add(time.Minute * -240): true,  // 240 minutes ago
		time.Now().Add(time.Minute * 30):   false, // 30 minutes in future
		time.Now().Add(time.Minute * 128):  false, // 128 minutes in future
	}

	for ts, expected := range examples {
		testClient.lastQueryTime = ts
		assert.Equal(t, expected, testClient.IsIdle())
	}
}

func testTest(t *testing.T) {
	examples := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "success",
			input: fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase),
			err:   nil,
		},
		{
			name:  "connection refused",
			input: "postgresql://localhost:5433/dbname",
			err:   ErrConnectionRefused,
		},
		{
			name:  "invalid user",
			input: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "foo", serverPassword, serverHost, serverPort, serverDatabase),
			err:   ErrAuthFailed,
		},
		// TODO:
		// This test fails when auth method in local pg_hba.conf is set to "trust".
		// When method is changed to "password" or "md5", client tests start prompting for password.
		// Leaving config set to "trust", and commenting out this test for now.
		// {
		// 	name:  "invalid password",
		// 	input: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", serverUser, "foo", serverHost, serverPort, serverDatabase),
		// 	err:   ErrAuthFailed,
		// },
		{
			name:  "invalid database",
			input: fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, "foo"),
			err:   ErrDatabaseNotExist,
		},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			conn, err := NewFromUrl(ex.input, nil)
			require.NoError(t, err)

			require.Equal(t, ex.err, conn.Test())
		})
	}
}

func testInfo(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		expected := []string{
			"session_user",
			"current_user",
			"current_database",
			"current_schemas",
			"inet_client_addr",
			"inet_client_port",
			"inet_server_addr",
			"inet_server_port",
			"version",
		}

		res, err := testClient.Info()
		assert.NoError(t, err)
		assert.Equal(t, expected, res.Columns)
	})

	t.Run("with restrictions", func(t *testing.T) {
		expected := []string{
			"session_user",
			"current_user",
			"current_database",
			"current_schemas",
			"version",
		}

		// Prepare a new user and database
		testClient.db.MustExec("DROP DATABASE IF EXISTS testdb")
		testClient.db.Exec("DROP OWNED BY IF EXISTS testuser") //nolint:all
		testClient.db.MustExec("DROP ROLE IF EXISTS testuser")
		testClient.db.MustExec("CREATE ROLE testuser WITH PASSWORD 'secret' LOGIN NOSUPERUSER NOINHERIT")
		testClient.db.MustExec("CREATE DATABASE testdb OWNER testuser")

		// Disable access to inet_ calls for new user
		url := fmt.Sprintf("postgres://%s:@%s:%s/testdb?sslmode=disable", serverUser, serverHost, serverPort)
		client, err := NewFromUrl(url, nil)
		assert.NoError(t, err)
		client.db.MustExec("REVOKE EXECUTE ON FUNCTION inet_client_addr() FROM PUBLIC")
		assert.NoError(t, client.Close())

		// Connect using new user
		url = fmt.Sprintf("postgres://testuser:secret@%s:%s/testdb?sslmode=disable", serverHost, serverPort)
		client, err = NewFromUrl(url, nil)
		assert.NoError(t, err)
		defer client.Close()

		res, err := client.Info()
		assert.NoError(t, err)
		assert.Equal(t, expected, res.Columns)
	})
}

func testActivity(t *testing.T) {
	expected := []string{"datid", "pid", "query", "query_start", "state", "client_addr"}

	res, err := testClient.Activity()
	assert.NoError(t, err)
	assertMatches(t, expected, res.Columns)
}

func testDatabases(t *testing.T) {
	res, err := testClient.Databases()
	assert.NoError(t, err)
	assertMatches(t, []string{"booktown", "postgres"}, res)
}

func testSchemas(t *testing.T) {
	res, err := testClient.Schemas()
	assert.NoError(t, err)
	assert.Equal(t, []string{"public"}, res)
}

func testObjects(t *testing.T) {
	res, err := testClient.Objects()
	objects := ObjectsFromResult(res)

	tables := []string{
		"alternate_stock",
		"authors",
		"book_backup",
		"book_queue",
		"books",
		"customers",
		"daily_inventory",
		"distinguished_authors",
		"dummies",
		"editions",
		"employees",
		"favorite_authors",
		"favorite_books",
		"money_example",
		"my_list",
		"numeric_values",
		"publishers",
		"schedules",
		"shipments",
		"states",
		"stock",
		"stock_backup",
		"subjects",
		"text_sorting",
	}

	functions := []string{
		"add_shipment",
		"add_two_loop",
		"books_by_subject",
		"compound_word",
		"count_by_two",
		"double_price",
		"extract_all_titles",
		"extract_all_titles2",
		"extract_title",
		"first",
		"get_author",
		"get_author",
		"get_customer_id",
		"get_customer_name",
		"html_linebreaks",
		"in_stock",
		"isbn_to_title",
		"mixed",
		"raise_test",
		"ship_item",
		"stock_amount",
		"test",
		"title",
		"triple_price",
	}

	assert.NoError(t, err)
	assert.Equal(t, []string{"oid", "schema", "name", "type", "owner", "comment"}, res.Columns)
	assert.Equal(t, []string{"public"}, mapKeys(objects))
	assert.Equal(t, tables, objectNames(objects["public"].Tables))
	assertMatches(t, functions, objectNames(objects["public"].Functions))
	assert.Equal(t, []string{"recent_shipments", "stock_view"}, objectNames(objects["public"].Views))
	assert.Equal(t, []string{"author_ids", "book_ids", "shipments_ship_id_seq", "subject_ids"}, objectNames(objects["public"].Sequences))

	major, minor := pgVersion()
	if minor == 0 || minor >= 3 {
		assert.Equal(t, []string{"m_stock_view"}, objectNames(objects["public"].MaterializedViews))
	} else {
		t.Logf("Skipping materialized view on %d.%d\n", major, minor)
	}
}

func testTable(t *testing.T) {
	columns := []string{
		"column_name",
		"data_type",
		"is_nullable",
		"character_maximum_length",
		"character_set_catalog",
		"column_default",
		"comment",
	}

	res, err := testClient.Table("books")
	assert.NoError(t, err)
	assert.Equal(t, columns, res.Columns)
	assert.Equal(t, 4, len(res.Rows))
}

func testTableRows(t *testing.T) {
	res, err := testClient.TableRows("books", RowsOptions{})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func testTableInfo(t *testing.T) {
	res, err := testClient.TableInfo("books")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 1, len(res.Rows))
}

func testEstimatedTableRowsCount(t *testing.T) {
	res, err := testClient.EstimatedTableRowsCount("books", RowsOptions{})
	assert.NoError(t, err)
	assert.Equal(t, []string{"reltuples"}, res.Columns)
	assert.Equal(t, []Row{{int64(15)}}, res.Rows)
}

func testTableRowsCount(t *testing.T) {
	res, err := testClient.TableRowsCount("books", RowsOptions{})
	assert.NoError(t, err)
	assert.Equal(t, []string{"count"}, res.Columns)
	assert.Equal(t, []Row{{int64(15)}}, res.Rows)
}

func testTableRowsCountWithLargeTable(t *testing.T) {
	testClient.db.MustExec(`CREATE TABLE large_table AS SELECT s FROM generate_series(1,1000000) s;`)
	testClient.db.MustExec(`VACUUM large_table;`)

	res, err := testClient.TableRowsCount("large_table", RowsOptions{})
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"reltuples"}, res.Columns)
	assert.Equal(t, []Row{{int64(1000000)}}, res.Rows)
}

func testTableIndexes(t *testing.T) {
	res, err := testClient.TableIndexes("books")
	assert.NoError(t, err)
	assert.Equal(t, []string{"index_name", "index_size", "index_definition"}, res.Columns)
	assert.Equal(t, 2, len(res.Rows))
}

func testTableConstraints(t *testing.T) {
	res, err := testClient.TableConstraints("editions")
	assert.NoError(t, err)
	assert.Equal(t, []string{"name", "definition"}, res.Columns)
	assert.Equal(t, Row{"pkey", "PRIMARY KEY (isbn)"}, res.Rows[0])
	assert.Equal(t, Row{"integrity", "CHECK (book_id IS NOT NULL AND edition IS NOT NULL)"}, res.Rows[1])
}

func testTableNameWithCamelCase(t *testing.T) {
	testClient.db.MustExec(`CREATE TABLE "exampleTable" (id int, name varchar);`)
	testClient.db.MustExec(`INSERT INTO "exampleTable" (id, name) VALUES (1, 'foo'), (2, 'bar');`)

	_, err := testClient.Table("exampleTable")
	assert.NoError(t, err)

	_, err = testClient.TableInfo("exampleTable")
	assert.NoError(t, err)

	_, err = testClient.TableConstraints("exampleTable")
	assert.NoError(t, err)

	_, err = testClient.TableIndexes("exampleTable")
	assert.NoError(t, err)

	_, err = testClient.TableRowsCount("exampleTable", RowsOptions{})
	assert.NoError(t, err)

	_, err = testClient.EstimatedTableRowsCount("exampleTable", RowsOptions{})
	assert.NoError(t, err)
}

func testQuery(t *testing.T) {
	t.Run("basic query", func(t *testing.T) {
		res, err := testClient.Query("SELECT * FROM books")
		assert.NoError(t, err)
		assert.Equal(t, 4, len(res.Columns))
		assert.Equal(t, 15, len(res.Rows))
	})

	t.Run("error", func(t *testing.T) {
		res, err := testClient.Query("SELCT * FROM books")
		assert.NotNil(t, err)
		assert.Equal(t, "pq: syntax error at or near \"SELCT\"", err.Error())
		assert.Nil(t, res)
	})

	t.Run("invalid table", func(t *testing.T) {
		res, err := testClient.Query("SELECT * FROM books2")
		assert.NotNil(t, err)
		assert.Equal(t, "pq: relation \"books2\" does not exist", err.Error())
		assert.Nil(t, res)
	})

	t.Run("timeout", func(t *testing.T) {
		testClient.queryTimeout = time.Millisecond * 100
		defer func() {
			testClient.queryTimeout = 0
		}()

		res, err := testClient.query("SELECT pg_sleep(1);")
		assert.Equal(t, "pq: canceling statement due to user request", err.Error())
		assert.Nil(t, res)
	})
}

func testUpdateQuery(t *testing.T) {
	t.Run("updating data", func(t *testing.T) {
		// Add new row
		testClient.db.MustExec("INSERT INTO books (id, title) VALUES (8888, 'Test Book'), (8889, 'Test Book 2')")

		// Update without return values
		res, err := testClient.Query("UPDATE books SET title = 'Foo' WHERE id >= 8888 AND id <= 8889")
		assert.NoError(t, err)
		assert.Equal(t, "Rows Affected", res.Columns[0])
		assert.Equal(t, int64(2), res.Rows[0][0])

		// Update with return values
		res, err = testClient.Query("UPDATE books SET title = 'Foo2' WHERE id >= 8888 AND id <= 8889 RETURNING id, title")
		assert.NoError(t, err)
		assert.Equal(t, []string{"id", "title"}, res.Columns)
		assert.Equal(t, Row{int64(8888), "Foo2"}, res.Rows[0])
		assert.Equal(t, Row{int64(8889), "Foo2"}, res.Rows[1])
	})

	t.Run("deleting data", func(t *testing.T) {
		// Add new row
		testClient.db.MustExec("INSERT INTO books (id, title) VALUES (9999, 'Test Book')")

		// Delete the existing row
		res, err := testClient.Query("DELETE FROM books WHERE id = 9999")
		assert.NoError(t, err)
		assert.Equal(t, "Rows Affected", res.Columns[0])
		assert.Equal(t, int64(1), res.Rows[0][0])

		// Deleting already deleted row
		res, err = testClient.Query("DELETE FROM books WHERE id = 9999")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), res.Rows[0][0])

		// Delete with returning value
		testClient.db.MustExec("INSERT INTO books (id, title) VALUES (9999, 'Test Book')")

		res, err = testClient.Query("DELETE FROM books WHERE id = 9999 RETURNING id")
		assert.NoError(t, err)
		assert.Equal(t, int64(9999), res.Rows[0][0])
	})
}

func testTableRowsOrderEscape(t *testing.T) {
	rows, err := testClient.TableRows("dummies", RowsOptions{SortColumn: "isDummy"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rows.Rows))

	rows, err = testClient.TableRows("dummies", RowsOptions{SortColumn: "isdummy"})
	assert.NotNil(t, err)
	assert.Equal(t, `pq: column "isdummy" does not exist`, err.Error())
	assert.Nil(t, rows)
}

func testFunctions(t *testing.T) {
	funcName := "get_customer_name"
	funcID := ""

	res, err := testClient.Objects()
	assert.NoError(t, err)

	for _, row := range res.Rows {
		if row[2] == funcName {
			funcID = row[0].(string)
			break
		}
	}

	res, err = testClient.Function("12345")
	assert.NoError(t, err)
	assertMatches(t, []string{"oid", "proname", "functiondef"}, res.Columns)
	assert.Equal(t, 0, len(res.Rows))

	res, err = testClient.Function(funcID)
	assert.NoError(t, err)
	assertMatches(t, []string{"oid", "proname", "functiondef"}, res.Columns)
	assert.Equal(t, 1, len(res.Rows))
	assert.Equal(t, funcName, res.Rows[0][1])
	assert.Contains(t, res.Rows[0][len(res.Columns)-1], "SELECT INTO customer_fname, customer_lname")
}

func testResult(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		result, err := testClient.Query("SELECT * FROM books LIMIT 1")
		assert.NoError(t, err)
		assert.Equal(t, `[{"author_id":4156,"id":7808,"subject_id":9,"title":"The Shining"}]`, string(result.JSON()))

		result, err = testClient.Query("SELECT 'NaN'::float AS value;")
		assert.NoError(t, err)
		assert.Equal(t, `[{"value":null}]`, string(result.JSON()))
	})

	t.Run("csv", func(t *testing.T) {
		expected := "id,title,author_id,subject_id\n156,The Tell-Tale Heart,115,9\n"

		res, err := testClient.Query("SELECT * FROM books ORDER BY id ASC LIMIT 1")
		assert.NoError(t, err)
		assert.Equal(t, expected, string(res.CSV()))
	})
}

func testHistory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, err := testClient.Query("SELECT * FROM books WHERE id = 12345")
		query := testClient.History[len(testClient.History)-1].Query
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM books WHERE id = 12345", query)
	})

	t.Run("failed query", func(t *testing.T) {
		_, err := testClient.Query("SELECT * FROM books123")
		query := testClient.History[len(testClient.History)-1].Query
		assert.NotNil(t, err)
		assert.NotEqual(t, "SELECT * FROM books123", query)
	})

	t.Run("unique queries", func(t *testing.T) {
		url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)

		client, _ := NewFromUrl(url, nil)
		defer client.Close()

		for i := 0; i < 3; i++ {
			_, err := client.Query("SELECT * FROM books WHERE id = 1")
			assert.NoError(t, err)
		}

		assert.Equal(t, 1, len(client.History))
		assert.Equal(t, "SELECT * FROM books WHERE id = 1", client.History[0].Query)
	})
}

func testReadOnlyMode(t *testing.T) {
	command.Opts.ReadOnly = true
	defer func() {
		command.Opts.ReadOnly = false
	}()

	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, _ := NewFromUrl(url, nil)
	defer client.Close()

	err := client.SetReadOnlyMode()
	assert.NoError(t, err)

	_, err = client.Query("\nCREATE TABLE foobar(id integer);\n")
	assert.NotNil(t, err)
	assert.Error(t, err, "query contains keywords not allowed in read-only mode")

	// Turn off guard
	_, err = client.db.Exec("SET default_transaction_read_only=off;")
	assert.NoError(t, err)

	_, err = client.Query("\nCREATE TABLE foobar(id integer);\n")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "query contains keywords not allowed in read-only mode")

	_, err = client.Query("-- CREATE TABLE foobar(id integer);\nSELECT 'foo';")
	assert.NoError(t, err)

	_, err = client.Query("/* CREATE TABLE foobar(id integer); */ SELECT 'foo';")
	assert.NoError(t, err)

	t.Run("with local readonly flag", func(t *testing.T) {
		command.Opts.ReadOnly = false
		client.readonly = true

		_, err := client.Query("INSERT INTO foobar(id) VALUES(1)")
		assert.Error(t, err, "query contains keywords not allowed in read-only mode")
	})
}

func testTablesStats(t *testing.T) {
	columns := []string{
		"schema_name",
		"table_name",
		"total_size",
		"data_size",
		"index_size",
		"estimated_rows_count",
		"estimated_rows",
		"index_to_data_ratio",
		"indexes_count",
		"columns_count",
	}

	result, err := testClient.TablesStats()
	assert.NoError(t, err)
	assert.Equal(t, columns, result.Columns)
}

func testConnContext(t *testing.T) {
	result, err := testClient.GetConnContext()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", result.Host)
	assert.Equal(t, "postgres", result.User)
	assert.Equal(t, "booktown", result.Database)
	assert.Equal(t, "default", result.Mode)
}

func TestAll(t *testing.T) {
	if onWindows() {
		t.Log("Unit testing on Windows platform is not supported.")
		return
	}

	initVars()
	setupCommands()
	teardown(t, false)
	setup()
	setupClient()

	testNewClientFromURL(t)
	testClientIdleTime(t)
	testTest(t)
	testInfo(t)
	testActivity(t)
	testDatabases(t)
	testSchemas(t)
	testObjects(t)
	testTable(t)
	testTableRows(t)
	testTableInfo(t)
	testEstimatedTableRowsCount(t)
	testTableRowsCount(t)
	testTableRowsCountWithLargeTable(t)
	testTableIndexes(t)
	testTableConstraints(t)
	testTableNameWithCamelCase(t)
	testQuery(t)
	testUpdateQuery(t)
	testTableRowsOrderEscape(t)
	testFunctions(t)
	testResult(t)
	testHistory(t)
	testReadOnlyMode(t)
	testDumpExport(t)
	testTablesStats(t)
	testConnContext(t)

	teardownClient()
	teardown(t, true)
}
