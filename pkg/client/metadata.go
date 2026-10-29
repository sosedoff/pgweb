package client

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/sosedoff/pgweb/pkg/statements"
)

// metadataCacheLimit caps the number of cached table metadata entries so that
// browsing a database with thousands of tables can't grow the cache unbounded.
const metadataCacheLimit = 256

var (
	// ErrNoPrimaryKey is returned when a row operation requires a primary key
	// but the relation does not have one.
	ErrNoPrimaryKey = errors.New("table does not have a primary key")

	// ErrNotEditable is returned when a relation can't be modified.
	ErrNotEditable = errors.New("table is not editable")
)

// ColumnMeta describes a single relation column.
type ColumnMeta struct {
	Name       string  `json:"name"`
	DataType   string  `json:"data_type"`
	UDTName    string  `json:"udt_name"`
	Nullable   bool    `json:"nullable"`
	Default    *string `json:"default"`
	Position   int     `json:"position"`
	Identity   bool    `json:"identity"`
	Generated  bool    `json:"generated"`
	Updatable  bool    `json:"updatable"`
	PrimaryKey bool    `json:"primary_key"`
}

// Writable returns true if the column value can be supplied by the user.
func (c ColumnMeta) Writable() bool {
	return c.Updatable && !c.Generated
}

// TableMetadata describes a relation and everything needed to build
// INSERT/UPDATE/DELETE statements and filter conditions for it.
type TableMetadata struct {
	Schema     string       `json:"schema"`
	Table      string       `json:"table"`
	Kind       string       `json:"kind"`
	Columns    []ColumnMeta `json:"columns"`
	PrimaryKey []string     `json:"primary_key"`
	Insertable bool         `json:"insertable"`
	Editable   bool         `json:"editable"`
}

// Column returns the column metadata by name, or nil when it does not exist.
func (m *TableMetadata) Column(name string) *ColumnMeta {
	for i := range m.Columns {
		if m.Columns[i].Name == name {
			return &m.Columns[i]
		}
	}
	return nil
}

// quotedName returns a fully qualified and quoted relation name.
func (m *TableMetadata) quotedName() string {
	return fmt.Sprintf("%s.%s", quoteIdent(m.Schema), quoteIdent(m.Table))
}

// quoteIdent escapes a postgres identifier. Callers must always validate the
// identifier against table metadata first, quoting alone is not a whitelist.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// TableMetadata returns the metadata for a given relation. Results are cached
// per connection since schema changes during a browsing session are rare and
// the data is requested on every row operation.
func (client *Client) TableMetadata(table string) (*TableMetadata, error) {
	schema, name := getSchemaAndTable(table)
	key := schema + "." + name

	client.metadataMu.RLock()
	cached := client.metadataCache[key]
	client.metadataMu.RUnlock()

	if cached != nil {
		return cached, nil
	}

	res, err := client.query(statements.TableMetadata, schema, name)
	if err != nil {
		return nil, err
	}

	meta, err := buildTableMetadata(schema, name, res)
	if err != nil {
		return nil, err
	}

	client.metadataMu.Lock()
	if len(client.metadataCache) >= metadataCacheLimit {
		client.metadataCache = map[string]*TableMetadata{}
	}
	if client.metadataCache == nil {
		client.metadataCache = map[string]*TableMetadata{}
	}
	client.metadataCache[key] = meta
	client.metadataMu.Unlock()

	return meta, nil
}

// ResetTableMetadata drops the cached metadata for a single relation, or for
// all relations when the name is empty.
func (client *Client) ResetTableMetadata(table string) {
	client.metadataMu.Lock()
	defer client.metadataMu.Unlock()

	if table == "" {
		client.metadataCache = map[string]*TableMetadata{}
		return
	}

	schema, name := getSchemaAndTable(table)
	delete(client.metadataCache, schema+"."+name)
}

// buildTableMetadata converts a raw metadata query result into TableMetadata.
func buildTableMetadata(schema, table string, res *Result) (*TableMetadata, error) {
	if res == nil || len(res.Rows) == 0 {
		return nil, fmt.Errorf("relation %q does not exist", schema+"."+table)
	}

	meta := TableMetadata{
		Schema:     schema,
		Table:      table,
		Columns:    make([]ColumnMeta, 0, len(res.Rows)),
		PrimaryKey: []string{},
	}

	// Primary key columns must be ordered by their position in the index.
	type pkEntry struct {
		name     string
		position int
	}
	pkEntries := []pkEntry{}

	idx := columnIndexes(res.Columns)

	for _, row := range res.Rows {
		col := ColumnMeta{
			Name:      stringValue(row, idx.get("column_name")),
			DataType:  stringValue(row, idx.get("data_type")),
			UDTName:   stringValue(row, idx.get("udt_name")),
			Nullable:  boolValue(row, idx.get("is_nullable")),
			Position:  intValue(row, idx.get("ordinal_position")),
			Identity:  boolValue(row, idx.get("is_identity")),
			Generated: boolValue(row, idx.get("is_generated")),
			Updatable: boolValue(row, idx.get("is_updatable")),
		}

		if def := row.value(idx.get("column_default")); def != nil {
			str := fmt.Sprintf("%v", def)
			col.Default = &str
		}

		if pos := intValue(row, idx.get("key_position")); pos > 0 {
			col.PrimaryKey = true
			pkEntries = append(pkEntries, pkEntry{name: col.Name, position: pos})
		}

		if meta.Kind == "" {
			meta.Kind = stringValue(row, idx.get("relkind"))
		}

		meta.Columns = append(meta.Columns, col)
	}

	sort.Slice(pkEntries, func(i, j int) bool {
		return pkEntries[i].position < pkEntries[j].position
	})
	for _, entry := range pkEntries {
		meta.PrimaryKey = append(meta.PrimaryKey, entry.name)
	}

	for _, col := range meta.Columns {
		if col.Writable() {
			meta.Insertable = true
			break
		}
	}
	meta.Editable = meta.Insertable && len(meta.PrimaryKey) > 0

	return &meta, nil
}

// value safely reads a row element at the given index.
func (r Row) value(idx int) interface{} {
	if idx < 0 || idx >= len(r) {
		return nil
	}
	return r[idx]
}

// resultColumns maps result column names to their position.
type resultColumns map[string]int

// get returns the column position, or -1 when the column is not present.
func (c resultColumns) get(name string) int {
	if idx, ok := c[name]; ok {
		return idx
	}
	return -1
}

func columnIndexes(columns []string) resultColumns {
	indexes := make(resultColumns, len(columns))
	for i, name := range columns {
		indexes[name] = i
	}
	return indexes
}

func stringValue(row Row, idx int) string {
	val := row.value(idx)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

func boolValue(row Row, idx int) bool {
	val := row.value(idx)
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "YES" || v == "t"
	}
	return false
}

func intValue(row Row, idx int) int {
	val := row.value(idx)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return int(v)
	case int:
		return v
	case float64:
		return int(v)
	case string:
		var num int
		fmt.Sscanf(v, "%d", &num) //nolint
		return num
	}
	return 0
}
