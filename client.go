package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type Client struct {
	db               *sqlx.DB
	history          []HistoryRecord
	connectionString string
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

func NewClient() (*Client, error) {
	str, err := buildConnectionString(options)

	if options.Debug && str != "" {
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
		connectionString: str,
		history:          NewHistory(),
	}

	return &client, nil
}

func NewClientFromUrl(url string) (*Client, error) {
	if options.Debug {
		fmt.Println("Creating a new client for:", url)
	}

	db, err := sqlx.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	client := Client{
		db:               db,
		connectionString: url,
		history:          NewHistory(),
	}

	return &client, nil
}

func (client *Client) Test() error {
	return client.db.Ping()
}

func (client *Client) Info() (*Result, error) {
	return client.query(PG_INFO)
}

func (client *Client) Databases() ([]string, error) {
	return client.fetchRows(PG_DATABASES)
}

func (client *Client) Tables() ([]string, error) {
	return client.fetchRows(PG_TABLES)
}

func (client *Client) Table(table string) (*Result, error) {
	return client.query(PG_TABLE_SCHEMA, table)
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
	return client.query(PG_TABLE_INFO, table)
}

func (client *Client) TableIndexes(table string) (*Result, error) {
	res, err := client.query(PG_TABLE_INDEXES, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (client *Client) Query(query string) (*Result, error) {
	res, err := client.query(query)

	// Save history records only if query did not fail
	if err == nil {
		client.history = append(client.history, NewHistoryRecord(query))
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
