package client

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/sosedoff/pgweb/pkg/command"
)

var (
	// ErrReadOnlyMode is returned when a write is attempted on a read-only connection.
	ErrReadOnlyMode = errors.New("connection is in read-only mode")

	// ErrNoValues is returned when a write request does not contain any values.
	ErrNoValues = errors.New("no column values provided")

	// ErrRowNotFound is returned when a row operation did not match any row.
	ErrRowNotFound = errors.New("row not found, it might have been modified by another session")

	// ErrMultipleRows is returned when a row operation matched more than one row.
	ErrMultipleRows = errors.New("operation would affect multiple rows, aborted")
)

// writeAllowed returns an error when the connection can't perform writes.
func (client *Client) writeAllowed() error {
	if client.readonly || command.Opts.ReadOnly {
		return ErrReadOnlyMode
	}
	return nil
}

// InsertRow inserts a single row and returns the resulting record.
//
// Column names are validated against the table metadata, so only existing and
// writable columns can be referenced. Values are always passed as bind
// parameters.
func (client *Client) InsertRow(table string, values map[string]interface{}) (*Result, error) {
	if err := client.writeAllowed(); err != nil {
		return nil, err
	}

	meta, err := client.TableMetadata(table)
	if err != nil {
		return nil, err
	}
	if !meta.Insertable {
		return nil, ErrNotEditable
	}

	columns, err := writableColumns(meta, values)
	if err != nil {
		return nil, err
	}

	// An insert with no values still has to produce a valid statement so that
	// tables consisting only of defaults can get a new record.
	if len(columns) == 0 {
		return client.query(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING *", meta.quotedName()))
	}

	names := make([]string, len(columns))
	marks := make([]string, len(columns))
	args := make([]interface{}, len(columns))

	for i, column := range columns {
		names[i] = quoteIdent(column)
		marks[i] = fmt.Sprintf("$%d", i+1)
		args[i] = values[column]
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		meta.quotedName(),
		strings.Join(names, ", "),
		strings.Join(marks, ", "),
	)

	return client.query(sql, args...)
}

// UpdateRow updates the given columns of a single row identified by its
// primary key and returns the updated record.
func (client *Client) UpdateRow(table string, key map[string]interface{}, values map[string]interface{}) (*Result, error) {
	if err := client.writeAllowed(); err != nil {
		return nil, err
	}

	meta, err := client.TableMetadata(table)
	if err != nil {
		return nil, err
	}
	if !meta.Editable {
		if len(meta.PrimaryKey) == 0 {
			return nil, ErrNoPrimaryKey
		}
		return nil, ErrNotEditable
	}

	columns, err := writableColumns(meta, values)
	if err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, ErrNoValues
	}

	assignments := make([]string, len(columns))
	args := make([]interface{}, 0, len(columns)+len(meta.PrimaryKey))

	for i, column := range columns {
		assignments[i] = fmt.Sprintf("%s = $%d", quoteIdent(column), i+1)
		args = append(args, values[column])
	}

	condition, keyArgs, err := primaryKeyCondition(meta, key, len(args))
	if err != nil {
		return nil, err
	}
	args = append(args, keyArgs...)

	sql := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s RETURNING *",
		meta.quotedName(),
		strings.Join(assignments, ", "),
		condition,
	)

	res, err := client.query(sql, args...)
	if err != nil {
		return nil, err
	}

	switch len(res.Rows) {
	case 1:
		return res, nil
	case 0:
		return nil, ErrRowNotFound
	default:
		return nil, ErrMultipleRows
	}
}

// DeleteRow removes a single row identified by its primary key.
func (client *Client) DeleteRow(table string, key map[string]interface{}) (*Result, error) {
	if err := client.writeAllowed(); err != nil {
		return nil, err
	}

	meta, err := client.TableMetadata(table)
	if err != nil {
		return nil, err
	}
	if len(meta.PrimaryKey) == 0 {
		return nil, ErrNoPrimaryKey
	}

	condition, args, err := primaryKeyCondition(meta, key, 0)
	if err != nil {
		return nil, err
	}

	// RETURNING makes the delete go through the regular query path so that the
	// number of affected rows can be verified before reporting success.
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s RETURNING *", meta.quotedName(), condition)

	res, err := client.query(sql, args...)
	if err != nil {
		return nil, err
	}

	switch len(res.Rows) {
	case 1:
		return res, nil
	case 0:
		return nil, ErrRowNotFound
	default:
		return nil, ErrMultipleRows
	}
}

// writableColumns validates the given values against the table metadata and
// returns the sorted list of column names to write.
func writableColumns(meta *TableMetadata, values map[string]interface{}) ([]string, error) {
	columns := make([]string, 0, len(values))

	for name := range values {
		col := meta.Column(name)
		if col == nil {
			return nil, fmt.Errorf("unknown column: %s", name)
		}
		if !col.Writable() {
			return nil, fmt.Errorf("column %s is not editable", name)
		}
		columns = append(columns, name)
	}

	// Sorting keeps generated statements stable and testable.
	sort.Strings(columns)

	return columns, nil
}

// primaryKeyCondition builds a parameterized WHERE clause matching a single
// row by its full primary key. Partial keys are rejected so that an operation
// can never affect more rows than intended.
func primaryKeyCondition(meta *TableMetadata, key map[string]interface{}, startIdx int) (string, []interface{}, error) {
	if len(meta.PrimaryKey) == 0 {
		return "", nil, ErrNoPrimaryKey
	}

	conditions := make([]string, 0, len(meta.PrimaryKey))
	args := make([]interface{}, 0, len(meta.PrimaryKey))

	for _, name := range meta.PrimaryKey {
		value, ok := key[name]
		if !ok {
			return "", nil, fmt.Errorf("primary key column %s is required", name)
		}
		if value == nil {
			return "", nil, fmt.Errorf("primary key column %s can't be null", name)
		}

		startIdx++
		conditions = append(conditions, fmt.Sprintf("%s = $%d", quoteIdent(name), startIdx))
		args = append(args, value)
	}

	if len(key) != len(meta.PrimaryKey) {
		return "", nil, fmt.Errorf("primary key must contain exactly %d column(s)", len(meta.PrimaryKey))
	}

	return strings.Join(conditions, " AND "), args, nil
}
