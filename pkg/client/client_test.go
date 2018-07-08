package client

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sosedoff/pgweb/pkg/command"
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

func test_NewClientFromUrl(t *testing.T) {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, err := NewFromUrl(url, nil)

	if err != nil {
		defer client.Close()
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, url, client.ConnectionString)
}

func test_NewClientFromUrl2(t *testing.T) {
	url := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, err := NewFromUrl(url, nil)

	if err != nil {
		defer client.Close()
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, url, client.ConnectionString)
}

func test_ClientIdleTime(t *testing.T) {
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

func test_Test(t *testing.T) {
	assert.Equal(t, nil, testClient.Test())
}

func test_Info(t *testing.T) {
	res, err := testClient.Info()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func test_Activity(t *testing.T) {
	res, err := testClient.Activity()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func test_Databases(t *testing.T) {
	res, err := testClient.Databases()

	assert.Equal(t, nil, err)
	assert.Contains(t, res, "booktown")
	assert.Contains(t, res, "postgres")
}

func test_Objects(t *testing.T) {
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

func test_Table(t *testing.T) {
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

func test_TableRows(t *testing.T) {
	res, err := testClient.TableRows("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func test_TableInfo(t *testing.T) {
	res, err := testClient.TableInfo("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 1, len(res.Rows))
}

func test_EstimatedTableRowsCount(t *testing.T) {
	var count int64 = 15
	res, err := testClient.EstimatedTableRowsCount("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"reltuples"}, res.Columns)
	assert.Equal(t, []Row{Row{count}}, res.Rows)
}

func test_TableRowsCount(t *testing.T) {
	var count int64 = 15
	res, err := testClient.TableRowsCount("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"count"}, res.Columns)
	assert.Equal(t, []Row{Row{count}}, res.Rows)
}

func test_TableRowsCountWithLargeTable(t *testing.T) {
	var count int64 = 100010
	testClient.db.MustExec(`create table large_table as select s from generate_Series(1,100010) s;`)
	testClient.db.MustExec(`VACUUM large_table;`)
	res, err := testClient.TableRowsCount("large_table", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"reltuples"}, res.Columns)
	assert.Equal(t, []Row{Row{count}}, res.Rows)
}

func test_TableIndexes(t *testing.T) {
	res, err := testClient.TableIndexes("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(res.Columns))
	assert.Equal(t, 2, len(res.Rows))
}

func test_TableConstraints(t *testing.T) {
	res, err := testClient.TableConstraints("editions")

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"name", "definition"}, res.Columns)
	assert.Equal(t, Row{"pkey", "PRIMARY KEY (isbn)"}, res.Rows[0])
	assert.Equal(t, Row{"integrity", "CHECK (book_id IS NOT NULL AND edition IS NOT NULL)"}, res.Rows[1])
}

func test_Query(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func test_QueryError(t *testing.T) {
	res, err := testClient.Query("SELCT * FROM books")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "pq: syntax error at or near \"SELCT\"", err.Error())
	assert.Equal(t, true, res == nil)
}

func test_QueryInvalidTable(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM books2")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "pq: relation \"books2\" does not exist", err.Error())
	assert.Equal(t, true, res == nil)
}

func test_TableRowsOrderEscape(t *testing.T) {
	rows, err := testClient.TableRows("dummies", RowsOptions{SortColumn: "isDummy"})
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(rows.Rows))

	rows, err = testClient.TableRows("dummies", RowsOptions{SortColumn: "isdummy"})
	assert.NotEqual(t, nil, err)
	assert.Equal(t, `pq: column "isdummy" does not exist`, err.Error())
	assert.Equal(t, true, rows == nil)
}

func test_ResultCsv(t *testing.T) {
	res, _ := testClient.Query("SELECT * FROM books ORDER BY id ASC LIMIT 1")
	csv := res.CSV()

	expected := "id,title,author_id,subject_id\n156,The Tell-Tale Heart,115,9\n"

	assert.Equal(t, expected, string(csv))
}

func test_History(t *testing.T) {
	_, err := testClient.Query("SELECT * FROM books WHERE id = 12345")
	query := testClient.History[len(testClient.History)-1].Query

	assert.Equal(t, nil, err)
	assert.Equal(t, "SELECT * FROM books WHERE id = 12345", query)
}

func test_HistoryError(t *testing.T) {
	_, err := testClient.Query("SELECT * FROM books123")
	query := testClient.History[len(testClient.History)-1].Query

	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, "SELECT * FROM books123", query)
}

func test_HistoryUniqueness(t *testing.T) {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, _ := NewFromUrl(url, nil)

	client.Query("SELECT * FROM books WHERE id = 1")
	client.Query("SELECT * FROM books WHERE id = 1")

	assert.Equal(t, 1, len(client.History))
	assert.Equal(t, "SELECT * FROM books WHERE id = 1", client.History[0].Query)
}

func test_ReadOnlyMode(t *testing.T) {
	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", serverUser, serverHost, serverPort, serverDatabase)
	client, _ := NewFromUrl(url, nil)

	err := client.SetReadOnlyMode()
	assert.Equal(t, nil, err)

	_, err = client.Query("CREATE TABLE foobar(id integer);")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "in a read-only transaction")
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

	test_NewClientFromUrl(t)
	test_ClientIdleTime(t)
	test_Test(t)
	test_Info(t)
	test_Activity(t)
	test_Databases(t)
	test_Objects(t)
	test_Table(t)
	test_TableRows(t)
	test_TableInfo(t)
	test_EstimatedTableRowsCount(t)
	test_TableRowsCount(t)
	test_TableRowsCountWithLargeTable(t)
	test_TableIndexes(t)
	test_TableConstraints(t)
	test_Query(t)
	test_QueryError(t)
	test_QueryInvalidTable(t)
	test_TableRowsOrderEscape(t)
	test_ResultCsv(t)
	test_History(t)
	test_HistoryUniqueness(t)
	test_HistoryError(t)
	test_ReadOnlyMode(t)
	test_DumpExport(t)

	teardownClient()
	teardown()
}
