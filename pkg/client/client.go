package client

import (
	"fmt"
	"log"
	neturl "net/url"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/history"
	"github.com/sosedoff/pgweb/pkg/shared"
	"github.com/sosedoff/pgweb/pkg/statements"
)

type Client struct {
	db               *sqlx.DB
	tunnel           *Tunnel
	History          []history.Record `json:"history"`
	ConnectionString string           `json:"connection_string"`
}

// Struct to hold table rows browsing options
type RowsOptions struct {
	Where      string // Custom filter
	Offset     int    // Number of rows to skip
	Limit      int    // Number of rows to fetch
	SortColumn string // Column to sort by
	SortOrder  string // Sort direction (ASC, DESC)
}

func getSchemaAndTable(str string) (string, string) {
	chunks := strings.Split(str, ".")
	if len(chunks) == 1 {
		return "public", chunks[0]
	}
	return chunks[0], chunks[1]
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

func NewFromUrl(url string, sshInfo *shared.SSHInfo) (*Client, error) {
	var tunnel *Tunnel

	if sshInfo != nil {
		if command.Opts.Debug {
			fmt.Println("Opening SSH tunnel for:", sshInfo)
		}

		tunnel, err := NewTunnel(sshInfo, url)
		if err != nil {
			tunnel.Close()
			return nil, err
		}

		err = tunnel.Configure()
		if err != nil {
			tunnel.Close()
			return nil, err
		}

		go tunnel.Start()

		uri, err := neturl.Parse(url)
		if err != nil {
			tunnel.Close()
			return nil, err
		}

		// Override remote postgres port with local proxy port
		url = strings.Replace(url, uri.Host, fmt.Sprintf("127.0.0.1:%v", tunnel.Port), 1)
	}

	if command.Opts.Debug {
		fmt.Println("Creating a new client for:", url)
	}

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	client := Client{
		db:               db,
		tunnel:           tunnel,
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

func (client *Client) Objects() (*Result, error) {
	return client.query(statements.PG_OBJECTS)
}

func (client *Client) Table(table string) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	return client.query(statements.PG_TABLE_SCHEMA, schema, table)
}

func (client *Client) MaterializedView(name string) (*Result, error) {
	return client.query(statements.PG_MATERIALIZED_VIEW_SCHEMA, name)
}

func (client *Client) TableRows(table string, opts RowsOptions) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	sql := fmt.Sprintf(`SELECT * FROM "%s"."%s"`, schema, table)

	if opts.Where != "" {
		sql += fmt.Sprintf(" WHERE %s", opts.Where)
	}

	if opts.SortColumn != "" {
		if opts.SortOrder == "" {
			opts.SortOrder = "ASC"
		}

		sql += fmt.Sprintf(" ORDER BY %s %s", opts.SortColumn, opts.SortOrder)
	}

	if opts.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	if opts.Offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	return client.query(sql)
}

func (client *Client) TableRowsCount(table string, opts RowsOptions) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	sql := fmt.Sprintf(`SELECT COUNT(1) FROM "%s"."%s"`, schema, table)

	if opts.Where != "" {
		sql += fmt.Sprintf(" WHERE %s", opts.Where)
	}

	return client.query(sql)
}

func (client *Client) TableInfo(table string) (*Result, error) {
	return client.query(statements.PG_TABLE_INFO, table)
}

func (client *Client) TableIndexes(table string) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	res, err := client.query(statements.PG_TABLE_INDEXES, schema, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (client *Client) TableConstraints(table string) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	res, err := client.query(statements.PG_TABLE_CONSTRAINTS, schema, table)

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

	timeout := command.Opts.Timeout
	res, err := client.query("set statement_timeout to " + strconv.Itoa(timeout*1000) + "; --SELECT setting FROM pg_settings where name = 'statement_timeout';")

	if command.Opts.Debug {
		log.Println("Query Timeout Seconds: ", timeout)
	}

	res, err = client.query(query)

	// Save history records only if query did not fail
	if err == nil && !client.hasHistoryRecord(query) {
		client.History = append(client.History, history.NewRecord(query))
	}

	return res, err
}

func (client *Client) query(query string, args ...interface{}) (*Result, error) {
	action := strings.ToLower(strings.Split(query, " ")[0])
	if action == "update" || action == "delete" {
		res, err := client.db.Exec(query, args...)
		if err != nil {
			return nil, err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		result := Result{
			Columns: []string{"Rows Affected"},
			Rows: []Row{
				Row{affected},
			},
		}

		return &result, nil
	}

	rows, err := client.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Make sure to never return null colums
	if cols == nil {
		cols = []string{}
	}

	result := Result{
		Columns: cols,
		Rows:    []Row{},
	}

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

	result.PrepareBigints()

	return &result, nil
}

// Close database connection
func (client *Client) Close() error {
	if client.tunnel != nil {
		client.tunnel.Close()
	}

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

func (client *Client) hasHistoryRecord(query string) bool {
	result := false

	for _, record := range client.History {
		if record.Query == query {
			result = true
			break
		}
	}

	return result
}
