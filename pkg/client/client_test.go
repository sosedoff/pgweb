package client

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

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

func test_Test(t *testing.T) {
	assert.Equal(t, nil, testClient.Test())
}

func test_Info(t *testing.T) {
	res, err := testClient.Info()

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
	assert.Equal(t, []string{"schema", "name", "type", "owner"}, res.Columns)
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

func test_TableIndexes(t *testing.T) {
	res, err := testClient.TableIndexes("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(res.Columns))
	assert.Equal(t, 2, len(res.Rows))
}

func test_TableConstraints(t *testing.T) {
	res, err := testClient.TableConstraints("editions")

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(res.Columns))
	assert.Equal(t, 2, len(res.Rows))
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

func TestAll(t *testing.T) {
	if onWindows() {
		// Dont have access to windows machines at the moment...
		return
	}

	initVars()
	setupCommands()
	teardown()
	setup()
	setupClient()

	test_NewClientFromUrl(t)
	test_Test(t)
	test_Info(t)
	test_Databases(t)
	test_Objects(t)
	test_Table(t)
	test_TableRows(t)
	test_TableInfo(t)
	test_TableIndexes(t)
	test_TableConstraints(t)
	test_Query(t)
	test_QueryError(t)
	test_QueryInvalidTable(t)
	test_ResultCsv(t)
	test_History(t)
	test_HistoryError(t)

	teardownClient()
	teardown()
}
