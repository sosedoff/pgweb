package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	neturl "net/url"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/history"
	"github.com/sosedoff/pgweb/pkg/shared"
	"github.com/sosedoff/pgweb/pkg/statements"
)

type Client struct {
	db               *sqlx.DB
	tunnel           *Tunnel
	serverVersion    string
	serverType       string
	lastQueryTime    time.Time
	queryTimeout     time.Duration
	closed           bool
	External         bool             `json:"external"`
	History          []history.Record `json:"history"`
	ConnectionString string           `json:"connection_string"`
}

func getSchemaAndTable(str string) (string, string) {
	chunks := strings.Split(str, ".")
	if len(chunks) == 1 {
		return "public", chunks[0]
	}
	return chunks[0], chunks[1]
}

func New() (*Client, error) {
	str, err := connection.BuildStringFromOptions(command.Opts)

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

	client.init()
	return &client, nil
}

func NewFromUrl(url string, sshInfo *shared.SSHInfo) (*Client, error) {
	var tunnel *Tunnel

	if sshInfo != nil {
		if command.Opts.DisableSSH {
			return nil, fmt.Errorf("ssh connections are disabled")
		}
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

	uri, err := neturl.Parse(url)
	if err == nil && uri.Path == "" {
		return nil, fmt.Errorf("Database name is not provided")
	}

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	client := Client{
		db:               db,
		tunnel:           tunnel,
		serverType:       postgresType,
		ConnectionString: url,
		History:          history.New(),
	}

	client.init()
	return &client, nil
}

func NewFromBookmark(bookmark *bookmarks.Bookmark) (*Client, error) {
	var (
		connStr string
		err     error
	)

	options := bookmark.ConvertToOptions()
	if options.URL != "" {
		connStr = options.URL
	} else {
		connStr, err = connection.BuildStringFromOptions(options)
		if err != nil {
			return nil, err
		}
	}

	var sshInfo *shared.SSHInfo
	if !bookmark.SSHInfoIsEmpty() {
		sshInfo = bookmark.SSH
	}

	return NewFromUrl(connStr, sshInfo)
}

func (client *Client) init() {
	if command.Opts.QueryTimeout > 0 {
		client.queryTimeout = time.Second * time.Duration(command.Opts.QueryTimeout)
	}

	client.setServerVersion()
}

func (client *Client) setServerVersion() {
	res, err := client.query("SELECT version()")
	if err != nil || len(res.Rows) < 1 {
		return
	}

	version := res.Rows[0][0].(string)
	match, serverType, serverVersion := detectServerTypeAndVersion(version)
	if match {
		client.serverType = serverType
		client.serverVersion = serverVersion
	}
}

func (client *Client) Test() error {
	return client.db.Ping()
}

func (client *Client) TestWithTimeout(timeout time.Duration) (result error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Check connection status right away without waiting for the ticker to kick in.
	// We're expecting to get "connection refused" here for the most part.
	if err := client.db.PingContext(ctx); err == nil {
		return nil
	}

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result = client.db.PingContext(ctx)
			if result == nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (client *Client) Info() (*Result, error) {
	result, err := client.query(statements.Info)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "inet_") && (strings.Contains(msg, "not supported") || strings.Contains(msg, "permission denied")) {
			// Fetch client information without inet_ function calls
			result, err = client.query(statements.InfoSimple)
		}
	}
	return result, err
}

func (client *Client) Databases() ([]string, error) {
	return client.fetchRows(statements.Databases)
}

func (client *Client) Schemas() ([]string, error) {
	return client.fetchRows(statements.Schemas)
}

func (client *Client) Objects() (*Result, error) {
	return client.query(statements.Objects)
}

func (client *Client) Table(table string) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	return client.query(statements.TableSchema, schema, table)
}

func (client *Client) MaterializedView(name string) (*Result, error) {
	return client.query(statements.MaterializedView, name)
}

func (client *Client) Function(id string) (*Result, error) {
	return client.query(statements.Function, id)
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

		sql += fmt.Sprintf(` ORDER BY "%s" %s`, opts.SortColumn, opts.SortOrder)
	}

	if opts.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	if opts.Offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	return client.query(sql)
}

func (client *Client) EstimatedTableRowsCount(table string, opts RowsOptions) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	result, err := client.query(statements.EstimatedTableRowCount, schema, table)
	if err != nil {
		return nil, err
	}
	// float64 to int64 conversion
	estimatedRowsCount := result.Rows[0][0].(float64)
	result.Rows[0] = Row{int64(estimatedRowsCount)}

	return result, nil
}

func (client *Client) TableRowsCount(table string, opts RowsOptions) (*Result, error) {
	// Return postgres estimated rows count on empty filter
	if opts.Where == "" && client.serverType == postgresType {
		res, err := client.EstimatedTableRowsCount(table, opts)
		if err != nil {
			return nil, err
		}
		n := res.Rows[0][0].(int64)
		if n >= 100000 {
			return res, nil
		}
	}

	schema, tableName := getSchemaAndTable(table)
	sql := fmt.Sprintf(`SELECT COUNT(1) FROM "%s"."%s"`, schema, tableName)

	if opts.Where != "" {
		sql += fmt.Sprintf(" WHERE %s", opts.Where)
	}

	return client.query(sql)
}

func (client *Client) TableInfo(table string) (*Result, error) {
	if client.serverType == cockroachType {
		return client.query(statements.TableInfoCockroach)
	}
	schema, table := getSchemaAndTable(table)
	return client.query(statements.TableInfo, fmt.Sprintf(`"%s"."%s"`, schema, table))
}

func (client *Client) TableIndexes(table string) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	res, err := client.query(statements.TableIndexes, schema, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (client *Client) TableConstraints(table string) (*Result, error) {
	schema, table := getSchemaAndTable(table)
	res, err := client.query(statements.TableConstraints, schema, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (client *Client) TablesStats() (*Result, error) {
	return client.query(statements.TablesStats)
}

// Returns all active queriers on the server
func (client *Client) Activity() (*Result, error) {
	if client.serverType == cockroachType {
		return client.query("SHOW QUERIES")
	}

	version := getMajorMinorVersionString(client.serverVersion)
	query := statements.Activity[version]
	if query == "" {
		query = statements.Activity["default"]
	}

	return client.query(query)
}

func (client *Client) Query(query string) (*Result, error) {
	res, err := client.query(query)

	// Save history records only if query did not fail
	if err == nil && !client.hasHistoryRecord(query) {
		client.History = append(client.History, history.NewRecord(query))
	}

	return res, err
}

func (client *Client) SetReadOnlyMode() error {
	var value string
	if err := client.db.Get(&value, "SHOW default_transaction_read_only;"); err != nil {
		return err
	}

	if value == "off" {
		_, err := client.db.Exec("SET default_transaction_read_only=on;")
		return err
	}

	return nil
}

func (client *Client) ServerVersionInfo() string {
	return fmt.Sprintf("%s %s", client.serverType, client.serverVersion)
}

func (client *Client) ServerVersion() string {
	return client.serverVersion
}

func (client *Client) context() (context.Context, context.CancelFunc) {
	if client.queryTimeout > 0 {
		return context.WithTimeout(context.Background(), client.queryTimeout)
	}
	return context.Background(), func() {}
}

func (client *Client) exec(query string, args ...interface{}) (*Result, error) {
	ctx, cancel := client.context()
	defer cancel()

	queryStart := time.Now()
	res, err := client.db.ExecContext(ctx, query, args...)
	queryFinish := time.Now()
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
			{affected},
		},
		Stats: &ResultStats{
			ColumnsCount:    1,
			RowsCount:       1,
			QueryStartTime:  queryStart.UTC(),
			QueryFinishTime: queryFinish.UTC(),
			QueryDuration:   queryFinish.Sub(queryStart).Milliseconds(),
		},
	}

	return &result, nil
}

func (client *Client) query(query string, args ...interface{}) (*Result, error) {
	if client.db == nil {
		return nil, nil
	}

	// Update the last usage time
	defer func() {
		client.lastQueryTime = time.Now().UTC()
	}()

	// We're going to force-set transaction mode on every query.
	// This is needed so that default mode could not be changed by user.
	if command.Opts.ReadOnly {
		if err := client.SetReadOnlyMode(); err != nil {
			return nil, err
		}
		if containsRestrictedKeywords(query) {
			return nil, errors.New("query contains keywords not allowed in read-only mode")
		}
	}

	action := strings.ToLower(strings.Split(query, " ")[0])
	hasReturnValues := strings.Contains(strings.ToLower(query), " returning ")

	if (action == "update" || action == "delete") && !hasReturnValues {
		return client.exec(query, args...)
	}

	ctx, cancel := client.context()
	defer cancel()

	queryStart := time.Now()
	rows, err := client.db.QueryxContext(ctx, query, args...)
	queryFinish := time.Now()
	if err != nil {
		if command.Opts.Debug {
			log.Println("Failed query:", query, "\nArgs:", args)
		}
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Make sure to never return null columns
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

	result.Stats = &ResultStats{
		ColumnsCount:    len(cols),
		RowsCount:       len(result.Rows),
		QueryStartTime:  queryStart.UTC(),
		QueryFinishTime: queryFinish.UTC(),
		QueryDuration:   queryFinish.Sub(queryStart).Milliseconds(),
	}

	result.PostProcess()

	return &result, nil
}

// Close database connection
func (client *Client) Close() error {
	if client.closed {
		return nil
	}
	defer func() {
		client.closed = true
	}()

	if client.tunnel != nil {
		client.tunnel.Close()
	}

	if client.db != nil {
		return client.db.Close()
	}

	return nil
}

func (c *Client) IsClosed() bool {
	return c.closed
}

func (c *Client) LastQueryTime() time.Time {
	return c.lastQueryTime
}

func (client *Client) IsIdle() bool {
	mins := int(time.Since(client.lastQueryTime).Minutes())

	if command.Opts.ConnectionIdleTimeout > 0 {
		return mins >= command.Opts.ConnectionIdleTimeout
	}

	return false
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

type ConnContext struct {
	Host     string
	User     string
	Database string
	Mode     string
}

func (c ConnContext) String() string {
	return fmt.Sprintf(
		"host=%q user=%q database=%q mode=%q",
		c.Host, c.User, c.Database, c.Mode,
	)
}

// ConnContext returns information about current database connection
func (client *Client) GetConnContext() (*ConnContext, error) {
	url, err := neturl.Parse(client.ConnectionString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	connCtx := ConnContext{
		Host: url.Hostname(),
		Mode: "default",
	}

	if command.Opts.ReadOnly {
		connCtx.Mode = "readonly"
	}

	row := client.db.QueryRowContext(ctx, "SELECT current_user, current_database()")
	if err := row.Scan(&connCtx.User, &connCtx.Database); err != nil {
		return nil, err
	}

	return &connCtx, nil
}
