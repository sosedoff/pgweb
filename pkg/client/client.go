package client

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/history"
	"github.com/sosedoff/pgweb/pkg/statements"
)

type Client struct {
	db               *sqlx.DB
	History          []history.Record
	ConnectionString string
}

type Row []interface{}

type Result struct {
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

// Struct to hold table rows browsing options
type RowsOptions struct {
	Limit      int    // Number of rows to fetch
	SortColumn string // Column to sort by
	SortOrder  string // Sort direction (ASC, DESC)
}

func New() (*Client, error) {
	str, err := connection.BuildString(command.Opts)

	if command.Opts.Debug && str != "" {
		fmt.Println("Creating a new client for:", str)
	}

	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("postgres", str)

	if err != nil {
		return nil, err
	}

	client := Client{
		db:               db,
		ConnectionString: str,
		History:          history.New(),
	}

	return &client, nil
}

func NewFromUrl(url string) (*Client, error) {
	if command.Opts.Debug {
		fmt.Println("Creating a new client for:", url)
	}

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	client := Client{
		db:               db,
		ConnectionString: url,
		History:          history.New(),
	}

	return &client, nil
}

func (client *Client) Test() error {
	return client.db.Ping()
}

func (client *Client) Info() (*Result, error) {
	return client.query(statements.PG_INFO)
}

func (client *Client) Databases() ([]string, error) {
	return client.fetchRows(statements.PG_DATABASES)
}

func (client *Client) Schemas() ([]string, error) {
	return client.fetchRows(statements.PG_SCHEMAS)
}

func (client *Client) Tables() ([]string, error) {
	return client.fetchRows(statements.PG_TABLES)
}

func (client *Client) Table(table string) (*Result, error) {
	return client.query(statements.PG_TABLE_SCHEMA, table)
}

func (client *Client) TableRows(table string, opts RowsOptions) (*Result, error) {
	sql := fmt.Sprintf(`SELECT * FROM "%s"`, table)

	if opts.SortColumn != "" {
		if opts.SortOrder == "" {
			opts.SortOrder = "ASC"
		}

		sql += fmt.Sprintf(" ORDER BY %s %s", opts.SortColumn, opts.SortOrder)
	}

	if opts.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	return client.query(sql)
}

func (client *Client) TableInfo(table string) (*Result, error) {
	return client.query(statements.PG_TABLE_INFO, table)
}

func (client *Client) TableIndexes(table string) (*Result, error) {
	res, err := client.query(statements.PG_TABLE_INDEXES, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Returns all active queriers on the server
func (client *Client) Activity() (*Result, error) {
	return client.query(statements.PG_ACTIVITY)
}

func (client *Client) Query(query string) (*Result, error) {
	res, err := client.query(query)

	// Save history records only if query did not fail
	if err == nil {
		client.History = append(client.History, history.NewRecord(query))
	}

	return res, err
}

func (client *Client) query(query string, args ...interface{}) (*Result, error) {
	rows, err := client.db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := Result{Columns: cols}

	for rows.Next() {
		obj, err := rows.SliceScan()

		for i, item := range obj {
			if item == nil {
				obj[i] = nil
			} else {
				t := reflect.TypeOf(item).Kind().String()

				if t == "slice" {
					obj[i] = string(item.([]byte))
				}
			}
		}

		if err == nil {
			result.Rows = append(result.Rows, obj)
		}
	}

	return &result, nil
}

func (res *Result) Format() []map[string]interface{} {
	var items []map[string]interface{}

	for _, row := range res.Rows {
		item := make(map[string]interface{})

		for i, c := range res.Columns {
			item[c] = row[i]
		}

		items = append(items, item)
	}

	return items
}

func (res *Result) CSV() []byte {
	buff := &bytes.Buffer{}
	writer := csv.NewWriter(buff)

	writer.Write(res.Columns)

	for _, row := range res.Rows {
		record := make([]string, len(res.Columns))

		for i, item := range row {
			if item != nil {
				record[i] = fmt.Sprintf("%v", item)
			} else {
				record[i] = ""
			}
		}

		err := writer.Write(record)

		if err != nil {
			fmt.Println(err)
			break
		}
	}

	writer.Flush()
	return buff.Bytes()
}

// Close database connection
func (client *Client) Close() error {
	if client.db != nil {
		return client.db.Close()
	}
	return nil
}

// Fetch all rows as strings for a single column
func (client *Client) fetchRows(q string) ([]string, error) {
	res, err := client.query(q)

	if err != nil {
		return nil, err
	}

	// Init empty slice so json.Marshal will encode it to "[]" instead of "null"
	results := make([]string, 0)

	for _, row := range res.Rows {
		results = append(results, row[0].(string))
	}

	return results, nil
}
