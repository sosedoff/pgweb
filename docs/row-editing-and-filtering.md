# pgweb-plus v1.0 — row editing and multi-column filtering

Migration notes and performance assessment for the three v1.0 features:
inline cell editing, insert/edit/delete rows, and multi-column filtering.

## Migration notes

### Upgrading

No migration is required. The features are additive and need no schema changes,
no new configuration, and no new runtime dependencies. Existing bookmarks,
sessions, query history and connection strings are unaffected.

### Behaviour changes for existing users

| Area | Before | After |
| --- | --- | --- |
| Double-click on a cell | Always opened the content viewer | Opens the inline editor on editable cells, content viewer everywhere else |
| Rows filter toolbar | One condition, rebuilt on every apply | Conditions accumulate into a chip list shown above the grid |
| Rows filter SQL | Built as a raw `where` string in JavaScript | Built server-side as parameterized conditions |
| Rows toolbar | Filters only | Filters plus `New` / `Edit` / `Delete` buttons |

The legacy `where` form parameter of `GET /api/tables/:table/rows` is still
accepted and is combined with the structured `filters` parameter using `AND`.
Any external client that uses it keeps working unchanged.

### Disabling the write features

Row editing is on by default and turns off automatically in `--readonly` mode.
To disable it explicitly:

```
pgweb --no-row-editing
# or
PGWEB_NO_ROW_EDITING=1 pgweb
```

When disabled, `POST`, `PUT` and `DELETE /api/tables/:table/rows` return `403`,
and `GET /api/tables/:table/metadata` reports `insertable: false` and
`editable: false` so the UI hides every write affordance. Multi-column
filtering is read-only and stays available.

`GET /api/info` reports the state under `features.row_editing`.

### New API surface

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/api/tables/:table/metadata` | Column metadata, primary key, insertable/editable flags |
| `POST` | `/api/tables/:table/rows` | Body `{"values": {...}}` |
| `PUT` | `/api/tables/:table/rows` | Body `{"key": {...}, "values": {...}}` |
| `DELETE` | `/api/tables/:table/rows` | Body `{"key": {...}}` |

`GET /api/tables/:table/rows` accepts an additional `filters` parameter holding
a JSON array of `{"column": ..., "op": ..., "value": ...}` objects.

### Known limitations

- Editing or deleting a row requires a detected primary key. Composite keys are
  supported; partial keys are rejected.
- `bytea` columns are not editable as text.
- Generated and non-updatable columns are read-only.
- Only `AND` is supported between filters. `OR` groups and nested conditions
  are out of scope for v1.0.
- Inline editing treats an empty input as `NULL` on nullable columns. Use the
  row editor to store an empty string.

## Performance impact assessment

### Startup and binary size

No new Go or JavaScript dependencies were added. The feature ships as five new
Go files, one embedded SQL file and additions to the existing `app.js` and
`app.css`. Startup performs no extra work: the metadata cache is an empty map
created with the client and is only populated on demand.

### Memory

Per connection, the only new allocation is the table metadata cache: a
`map[string]*TableMetadata` holding at most 256 entries, cleared wholesale when
the limit is exceeded. A cached entry is a small struct plus one `ColumnMeta`
per column, so a typical 30-column table costs a few kilobytes. The cache is
guarded by a `sync.RWMutex` and read-locked on the hot path.

No result set is buffered beyond what pgweb already does. Every write returns
exactly one row via `RETURNING *`.

### Query behaviour

- **Browsing without filters** issues exactly the same two queries as before
  (rows plus count). `rowsCondition()` only loads metadata when filters are
  present, so there is no extra round trip on the common path.
- **Browsing with filters** adds one metadata lookup the first time a table is
  opened, served from cache afterwards. Filters are compiled into a `WHERE`
  clause with `$n` placeholders and executed by PostgreSQL, so index usage,
  `LIMIT`/`OFFSET` pagination and the row count all behave normally. No rows
  are ever loaded into memory for client-side filtering.
- **Writes** are single-statement DML, which PostgreSQL already executes
  atomically, so no explicit transaction is opened.

### Frontend

- The metadata request is cached per table in a plain object, so switching tabs
  or paging does not re-query the schema.
- Inline editing and the edit/delete modals patch the affected `<tr>` in place
  instead of reloading the grid. Only `New Row` triggers a reload, because the
  position of the new row is not known.
- Filter chips are rendered with jQuery into a single container; no state
  management library, templating engine or UI framework was introduced.

### Verification

Verified against PostgreSQL 16.8 and 10.6 with the Go unit and integration
suites, and end-to-end in headless Chrome against a live database: filtered
queries on a 31-column table returned correct counts and pages, inline edits
updated a single row without a reload, tables without a primary key and tables
with composite primary keys behaved as specified, and no JavaScript errors were
raised.
