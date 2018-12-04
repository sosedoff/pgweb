package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	neturl "net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuvistavie/securerandom"

	"github.com/sosedoff/pgweb/pkg/bookmarks"
	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
	"github.com/sosedoff/pgweb/pkg/discovery"
	"github.com/sosedoff/pgweb/pkg/shared"
)

var (
	// DbClient represents the active database connection in a single-session mode
	DbClient *client.Client

	// DbSessions represents the mapping for client connections
	DbSessions = map[string]*client.Client{}

	// Providers represents a list of all enabled providers
	Providers = map[string]discovery.Provider{}
)

// DB returns a database connection from the client context
func DB(c *gin.Context) *client.Client {
	if command.Opts.Sessions {
		return DbSessions[getSessionId(c.Request)]
	} else {
		return DbClient
	}
}

func setClient(c *gin.Context, newClient *client.Client) error {
	currentClient := DB(c)
	if currentClient != nil {
		currentClient.Close()
	}

	if !command.Opts.Sessions {
		DbClient = newClient
		return nil
	}

	sessionId := getSessionId(c.Request)
	if sessionId == "" {
		return errors.New("Session ID is required")
	}

	DbSessions[sessionId] = newClient
	return nil
}

func GetHome(c *gin.Context) {
	serveStaticAsset("/index.html", c)
}

func GetAsset(c *gin.Context) {
	serveStaticAsset(c.Params.ByName("path"), c)
}

func GetSessions(c *gin.Context) {
	// In debug mode endpoint will return a lot of sensitive information
	// like full database connection string and all query history.
	if command.Opts.Debug {
		successResponse(c, DbSessions)
		return
	}
	successResponse(c, gin.H{"sessions": len(DbSessions)})
}

func ConnectWithBackend(c *gin.Context) {
	// Setup a new backend client
	backend := Backend{
		Endpoint:    command.Opts.ConnectBackend,
		Token:       command.Opts.ConnectToken,
		PassHeaders: command.Opts.ConnectHeaders,
	}

	// Fetch connection credentials
	cred, err := backend.FetchCredential(c.Param("resource"), c)
	if err != nil {
		badRequest(c, err)
		return
	}

	// Make the new session
	sessionId, err := securerandom.Uuid()
	if err != nil {
		badRequest(c, err)
		return
	}
	c.Request.Header.Add("x-session-id", sessionId)

	// Connect to the database
	cl, err := client.NewFromUrl(cred.DatabaseURL, nil)
	if err != nil {
		badRequest(c, err)
		return
	}
	cl.External = true

	// Finalize session seetup
	_, err = cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		badRequest(c, err)
		return
	}

	c.Redirect(301, fmt.Sprintf("/%s?session=%s", command.Opts.Prefix, sessionId))
}

func Connect(c *gin.Context) {
	if command.Opts.LockSession {
		badRequest(c, "Session is locked")
		return
	}

	var sshInfo *shared.SSHInfo
	url := c.Request.FormValue("url")

	if url == "" {
		badRequest(c, "Url parameter is required")
		return
	}

	opts := command.Options{Url: url}
	url, err := connection.FormatUrl(opts)

	if err != nil {
		badRequest(c, err)
		return
	}

	if c.Request.FormValue("ssh") != "" {
		sshInfo = parseSshInfo(c)
	}

	cl, err := client.NewFromUrl(url, sshInfo)
	if err != nil {
		badRequest(c, err)
		return
	}

	err = cl.Test()
	if err != nil {
		badRequest(c, err)
		return
	}

	info, err := cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		badRequest(c, err)
		return
	}

	successResponse(c, info.Format()[0])
}

func SwitchDb(c *gin.Context) {
	if command.Opts.LockSession {
		badRequest(c, "Session is locked")
		return
	}

	name := c.Request.URL.Query().Get("db")
	if name == "" {
		name = c.Request.FormValue("db")
	}
	if name == "" {
		badRequest(c, "Database name is not provided")
		return
	}

	conn := DB(c)
	if conn == nil {
		badRequest(c, "Not connected")
		return
	}

	// Do not allow switching databases for connections from third-party backends
	if conn.External {
		badRequest(c, "Session is locked")
		return
	}

	currentUrl, err := neturl.Parse(conn.ConnectionString)
	if err != nil {
		badRequest(c, "Unable to parse current connection string")
		return
	}

	currentUrl.Path = name

	cl, err := client.NewFromUrl(currentUrl.String(), nil)
	if err != nil {
		badRequest(c, err)
		return
	}

	err = cl.Test()
	if err != nil {
		badRequest(c, err)
		return
	}

	info, err := cl.Info()
	if err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		badRequest(c, err)
		return
	}

	conn.Close()

	successResponse(c, info.Format()[0])
}

func Disconnect(c *gin.Context) {
	if command.Opts.LockSession {
		badRequest(c, "Session is locked")
		return
	}

	conn := DB(c)

	if conn == nil {
		badRequest(c, "Not connected")
		return
	}

	err := conn.Close()
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, gin.H{"success": true})
}

func GetDatabases(c *gin.Context) {
	conn := DB(c)
	if conn.External {
		errorResponse(c, 403, "Not permitted")
		return
	}

	names, err := DB(c).Databases()
	serveResult(c, names, err)
}

func GetObjects(c *gin.Context) {
	result, err := DB(c).Objects()
	if err != nil {
		badRequest(c, err)
		return
	}
	successResponse(c, client.ObjectsFromResult(result))
}

func RunQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		badRequest(c, "Query parameter is missing")
		return
	}

	HandleQuery(query, c)
}

func ExplainQuery(c *gin.Context) {
	query := cleanQuery(c.Request.FormValue("query"))

	if query == "" {
		badRequest(c, "Query parameter is missing")
		return
	}

	HandleQuery(fmt.Sprintf("EXPLAIN ANALYZE %s", query), c)
}

func GetSchemas(c *gin.Context) {
	res, err := DB(c).Schemas()
	serveResult(c, res, err)
}

func GetTable(c *gin.Context) {
	var res *client.Result
	var err error

	if c.Request.FormValue("type") == "materialized_view" {
		res, err = DB(c).MaterializedView(c.Params.ByName("table"))
	} else {
		res, err = DB(c).Table(c.Params.ByName("table"))
	}

	serveResult(c, res, err)
}

func GetTableRows(c *gin.Context) {
	offset, err := parseIntFormValue(c, "offset", 0)
	if err != nil {
		badRequest(c, err)
		return
	}

	limit, err := parseIntFormValue(c, "limit", 100)
	if err != nil {
		badRequest(c, err)
		return
	}

	opts := client.RowsOptions{
		Limit:      limit,
		Offset:     offset,
		SortColumn: c.Request.FormValue("sort_column"),
		SortOrder:  c.Request.FormValue("sort_order"),
		Where:      c.Request.FormValue("where"),
	}

	res, err := DB(c).TableRows(c.Params.ByName("table"), opts)
	if err != nil {
		badRequest(c, err)
		return
	}

	countRes, err := DB(c).TableRowsCount(c.Params.ByName("table"), opts)
	if err != nil {
		badRequest(c, err)
		return
	}

	numFetch := int64(opts.Limit)
	numOffset := int64(opts.Offset)
	numRows := countRes.Rows[0][0].(int64)
	numPages := numRows / numFetch

	if numPages*numFetch < numRows {
		numPages++
	}

	res.Pagination = &client.Pagination{
		Rows:    numRows,
		Page:    (numOffset / numFetch) + 1,
		Pages:   numPages,
		PerPage: numFetch,
	}

	serveResult(c, res, err)
}

func GetTableInfo(c *gin.Context) {
	res, err := DB(c).TableInfo(c.Params.ByName("table"))
	if err == nil {
		successResponse(c, res.Format()[0])
	} else {
		badRequest(c, err)
	}
}

func GetHistory(c *gin.Context) {
	successResponse(c, DB(c).History)
}

func GetConnectionInfo(c *gin.Context) {
	res, err := DB(c).Info()

	if err != nil {
		badRequest(c, err)
		return
	}

	info := res.Format()[0]
	info["session_lock"] = command.Opts.LockSession

	successResponse(c, info)
}

func GetActivity(c *gin.Context) {
	res, err := DB(c).Activity()
	serveResult(c, res, err)
}

func GetTableIndexes(c *gin.Context) {
	res, err := DB(c).TableIndexes(c.Params.ByName("table"))
	serveResult(c, res, err)
}

func GetTableConstraints(c *gin.Context) {
	res, err := DB(c).TableConstraints(c.Params.ByName("table"))
	serveResult(c, res, err)
}

func HandleQuery(query string, c *gin.Context) {
	rawQuery, err := base64.StdEncoding.DecodeString(desanitize64(query))
	if err == nil {
		query = string(rawQuery)
	}

	result, err := DB(c).Query(query)
	if err != nil {
		badRequest(c, err)
		return
	}

	format := getQueryParam(c, "format")
	filename := getQueryParam(c, "filename")

	if filename == "" {
		filename = fmt.Sprintf("pgweb-%v.%v", time.Now().Unix(), format)
	}

	if format != "" {
		c.Writer.Header().Set("Content-disposition", "attachment;filename="+filename)
	}

	switch format {
	case "csv":
		c.Data(200, "text/csv", result.CSV())
	case "json":
		c.Data(200, "application/json", result.JSON())
	case "xml":
		c.XML(200, result)
	default:
		c.JSON(200, result)
	}
}

func GetBookmarks(c *gin.Context) {
	bookmarks, err := bookmarks.ReadAll(bookmarks.Path(command.Opts.BookmarksDir))
	serveResult(c, bookmarks, err)
}

func GetInfo(c *gin.Context) {
	successResponse(c, gin.H{
		"version":    command.Version,
		"git_sha":    command.GitCommit,
		"build_time": command.BuildTime,
	})
}

// Export database or table data
func DataExport(c *gin.Context) {
	db := DB(c)

	info, err := db.Info()
	if err != nil {
		badRequest(c, err)
		return
	}

	dump := client.Dump{
		Table: strings.TrimSpace(c.Request.FormValue("table")),
	}

	// If pg_dump is not available the following code will not show an error in browser.
	// This is due to the headers being written first.
	if !dump.CanExport() {
		badRequest(c, "pg_dump is not found")
		return
	}

	formattedInfo := info.Format()[0]
	filename := formattedInfo["current_database"].(string)
	if dump.Table != "" {
		filename = filename + "_" + dump.Table
	}

	c.Header(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s.sql.gz"`, filename),
	)

	err = dump.Export(db.ConnectionString, c.Writer)
	if err != nil {
		badRequest(c, err)
	}
}

// DiscoveryIndex returns a list of all configured providers
func DiscoveryIndex(c *gin.Context) {
	if !command.Opts.Discovery {
		badRequest(c, "Discovery is not enabled")
		return
	}

	result := []map[string]string{}

	for _, provider := range Providers {
		result = append(result, map[string]string{
			"id":   provider.ID(),
			"name": provider.Name(),
		})
	}
	successResponse(c, result)
}

// DiscoveryList returns a list of all provider resources
func DiscoveryList(c *gin.Context) {
	if !command.Opts.Discovery {
		badRequest(c, "Discovery is not enabled")
		return
	}

	provider, ok := Providers[c.Param("provider")]
	if !ok {
		badRequest(c, "Invalid provider")
		return
	}

	resources, err := provider.List()
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, resources)
}

// DiscoveryConnect performs a provider resource lookup and connects to the database
// if the resource was found.
func DiscoveryConnect(c *gin.Context) {
	if !command.Opts.Discovery {
		badRequest(c, "Discovery is not enabled")
		return
	}

	provider, ok := Providers[c.Param("provider")]
	if !ok {
		badRequest(c, "Invalid provider")
		return
	}

	credential, err := provider.Get(c.Param("id"))
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, credential)
}

func RegisterProvider(p discovery.Provider) error {
	id := p.ID()
	if _, ok := Providers[id]; ok {
		return fmt.Errorf("Provider %s already registered", id)
	}
	Providers[id] = p
	return nil
}
