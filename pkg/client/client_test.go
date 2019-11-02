package client

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/sosedoff/pgweb/pkg/command"

	"github.com/stretchr/testify/assert"
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
	for k, _ := range data {
		result = append(result, k)
	}
	return result
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
		testClient.db.Close()
	}
}

func teardown() {
	_, err := exec.Command(
		testCommands["dropdb"],
		"-U", serverUser,
		"-h", serverHost,
		"-p", serverPort,
		serverDatabase,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Teardown error:", err)
	}
}

func testNewClientFromUrl(t *testing.T) {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, err := NewFromUrl(url, nil)

	if err != nil {
		defer client.Close()
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, url, client.ConnectionString)
}

func testNewClientFromUrl2(t *testing.T) {
	url := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, err := NewFromUrl(url, nil)

	if err != nil {
		defer client.Close()
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, url, client.ConnectionString)
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
	assert.Equal(t, nil, testClient.Test())
}

func testInfo(t *testing.T) {
	res, err := testClient.Info()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func testActivity(t *testing.T) {
	res, err := testClient.Activity()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func testDatabases(t *testing.T) {
	res, err := testClient.Databases()

	assert.Equal(t, nil, err)
	assert.Contains(t, res, "booktown")
	assert.Contains(t, res, "postgres")
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

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"schema", "name", "type", "owner", "comment"}, res.Columns)
	assert.Equal(t, []string{"public"}, mapKeys(objects))
	assert.Equal(t, tables, objects["public"].Tables)
	assert.Equal(t, []string{"recent_shipments", "stock_view"}, objects["public"].Views)
	assert.Equal(t, []string{"author_ids", "book_ids", "shipments_ship_id_seq", "subject_ids"}, objects["public"].Sequences)

	major, minor := pgVersion()
	if minor == 0 || minor >= 3 {
		assert.Equal(t, []string{"m_stock_view"}, objects["public"].MaterializedViews)
	} else {
		t.Logf("Skipping materialized view on %d.%d\n", major, minor)
	}
}

func testTable(t *testing.T) {
	res, err := testClient.Table("books")

	columns := []string{
		"column_name",
		"data_type",
		"is_nullable",
		"character_maximum_length",
		"character_set_catalog",
		"column_default",
		"comment",
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, columns, res.Columns)
	assert.Equal(t, 4, len(res.Rows))
}

func testTableRows(t *testing.T) {
	res, err := testClient.TableRows("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func testTableInfo(t *testing.T) {
	res, err := testClient.TableInfo("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 1, len(res.Rows))
}

func testEstimatedTableRowsCount(t *testing.T) {
	var count int64 = 15
	res, err := testClient.EstimatedTableRowsCount("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"reltuples"}, res.Columns)
	assert.Equal(t, []Row{Row{count}}, res.Rows)
}

func testTableRowsCount(t *testing.T) {
	var count int64 = 15
	res, err := testClient.TableRowsCount("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"count"}, res.Columns)
	assert.Equal(t, []Row{Row{count}}, res.Rows)
}

func testTableRowsCountWithLargeTable(t *testing.T) {
	var count int64 = 100010
	testClient.db.MustExec(`CREATE TABLE large_table AS SELECT s FROM generate_Series(1,100010) s;`)
	testClient.db.MustExec(`VACUUM large_table;`)
	res, err := testClient.TableRowsCount("large_table", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"reltuples"}, res.Columns)
	assert.Equal(t, []Row{Row{count}}, res.Rows)
}

func testTableIndexes(t *testing.T) {
	res, err := testClient.TableIndexes("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(res.Columns))
	assert.Equal(t, 2, len(res.Rows))
}

func testTableConstraints(t *testing.T) {
	res, err := testClient.TableConstraints("editions")

	assert.Equal(t, nil, err)
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
	res, err := testClient.Query("SELECT * FROM books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func testQueryError(t *testing.T) {
	res, err := testClient.Query("SELCT * FROM books")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "pq: syntax error at or near \"SELCT\"", err.Error())
	assert.Equal(t, true, res == nil)
}

func testQueryInvalidTable(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM books2")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "pq: relation \"books2\" does not exist", err.Error())
	assert.Equal(t, true, res == nil)
}

func testTableRowsOrderEscape(t *testing.T) {
	rows, err := testClient.TableRows("dummies", RowsOptions{SortColumn: "isDummy"})
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(rows.Rows))

	rows, err = testClient.TableRows("dummies", RowsOptions{SortColumn: "isdummy"})
	assert.NotEqual(t, nil, err)
	assert.Equal(t, `pq: column "isdummy" does not exist`, err.Error())
	assert.Equal(t, true, rows == nil)
}

func testResultJSON(t *testing.T) {
	result, err := testClient.Query("SELECT 'NaN'::float AS value;")

	assert.NoError(t, err)
	assert.Equal(t, `[{"value":null}]`, string(result.JSON()))
}

func testResultCsv(t *testing.T) {
	res, _ := testClient.Query("SELECT * FROM books ORDER BY id ASC LIMIT 1")
	csv := res.CSV()

	expected := "id,title,author_id,subject_id\n156,The Tell-Tale Heart,115,9\n"

	assert.Equal(t, expected, string(csv))
}

func testHistory(t *testing.T) {
	_, err := testClient.Query("SELECT * FROM books WHERE id = 12345")
	query := testClient.History[len(testClient.History)-1].Query

	assert.Equal(t, nil, err)
	assert.Equal(t, "SELECT * FROM books WHERE id = 12345", query)
}

func testHistoryError(t *testing.T) {
	_, err := testClient.Query("SELECT * FROM books123")
	query := testClient.History[len(testClient.History)-1].Query

	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, "SELECT * FROM books123", query)
}

func testHistoryUniqueness(t *testing.T) {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, _ := NewFromUrl(url, nil)

	client.Query("SELECT * FROM books WHERE id = 1")
	client.Query("SELECT * FROM books WHERE id = 1")

	assert.Equal(t, 1, len(client.History))
	assert.Equal(t, "SELECT * FROM books WHERE id = 1", client.History[0].Query)
}

func testReadOnlyMode(t *testing.T) {
	command.Opts.ReadOnly = true
	defer func() {
		command.Opts.ReadOnly = false
	}()

	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, _ := NewFromUrl(url, nil)

	err := client.SetReadOnlyMode()
	assert.NoError(t, err)

	_, err = client.Query("\nCREATE TABLE foobar(id integer);\n")
	assert.NotNil(t, err)
	assert.Error(t, err, "query contains keywords not allowed in read-only mode")

	// Turn off guard
	client.db.Exec("SET default_transaction_read_only=off;")

	_, err = client.Query("\nCREATE TABLE foobar(id integer);\n")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "query contains keywords not allowed in read-only mode")

	_, err = client.Query("-- CREATE TABLE foobar(id integer);\nSELECT 'foo';")
	assert.NoError(t, err)

	_, err = client.Query("/* CREATE TABLE foobar(id integer); */ SELECT 'foo';")
	assert.NoError(t, err)
}

func TestAll(t *testing.T) {
	if onWindows() {
		t.Log("Unit testing on Windows platform is not supported.")
		return
	}

	initVars()
	setupCommands()
	teardown()
	setup()
	setupClient()

	testNewClientFromUrl(t)
	testClientIdleTime(t)
	testTest(t)
	testInfo(t)
	testActivity(t)
	testDatabases(t)
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
	testQueryError(t)
	testQueryInvalidTable(t)
	testTableRowsOrderEscape(t)
	testResultJSON(t)
	testResultCsv(t)
	testHistory(t)
	testHistoryUniqueness(t)
	testHistoryError(t)
	testReadOnlyMode(t)
	testDumpExport(t)

	teardownClient()
	teardown()
}
