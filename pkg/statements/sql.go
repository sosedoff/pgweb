package statements

const (
	// ---------------------------------------------------------------------------

	PG_DATABASES = `
SELECT
  datname
FROM
  pg_database
WHERE
  NOT datistemplate
ORDER BY
  datname ASC`

	// ---------------------------------------------------------------------------

	PG_SCHEMAS = `
SELECT
  schema_name
FROM
  information_schema.schemata
ORDER BY
  schema_name ASC`

	// ---------------------------------------------------------------------------

	PG_INFO = `
SELECT
  session_user,
  current_user,
  current_database(),
  current_schemas(false),
  inet_client_addr(),
  inet_client_port(),
  inet_server_addr(),
  inet_server_port(),
  version()`

	// ---------------------------------------------------------------------------

	PG_TABLE_INDEXES = `
SELECT
  indexname, indexdef
FROM
  pg_indexes
WHERE
  schemaname = $1 AND
  tablename = $2`

	// ---------------------------------------------------------------------------

	PG_TABLE_CONSTRAINTS = `
SELECT
  pg_get_constraintdef(c.oid, true) as condef
FROM
  pg_constraint c
JOIN
  pg_namespace n ON n.oid = c.connamespace
JOIN
  pg_class cl ON cl.oid = c.conrelid
WHERE
  n.nspname = $1 AND
  relname = $2
ORDER BY
  contype desc`

	// ---------------------------------------------------------------------------

	PG_TABLE_INFO = `
SELECT
  pg_size_pretty(pg_table_size($1)) AS data_size,
  pg_size_pretty(pg_indexes_size($1)) AS index_size,
  pg_size_pretty(pg_total_relation_size($1)) AS total_size,
  (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count`

	// ---------------------------------------------------------------------------

	PG_TABLE_SCHEMA = `
SELECT
  column_name,
  data_type,
  is_nullable,
  character_maximum_length,
  character_set_catalog,
  column_default
FROM
  information_schema.columns
WHERE
  table_schema = $1 AND
  table_name = $2`

	// ---------------------------------------------------------------------------

	PG_MATERIALIZED_VIEW_SCHEMA = `
SELECT 
  attname as column_name, 
  atttypid::regtype AS data_type,
  (case when attnotnull IS TRUE then 'NO' else 'YES' end) as is_nullable,
  null as character_maximum_length,
  null as character_set_catalog,
  null as column_default
FROM
  pg_attribute
WHERE
  attrelid = $1::regclass AND
  attnum > 0 AND
  NOT attisdropped`

	// ---------------------------------------------------------------------------

	PG_ACTIVITY = `
SELECT
  datname,
  query,
  state,
  waiting,
  query_start,
  state_change,
  pid,
  datid,
  application_name,
  client_addr
FROM
  pg_stat_activity
WHERE
  state IS NOT NULL`

	// ---------------------------------------------------------------------------

	PG_OBJECTS = `
SELECT
  n.nspname as "schema",
  c.relname as "name",
  CASE c.relkind
    WHEN 'r' THEN 'table'
    WHEN 'v' THEN 'view'
    WHEN 'm' THEN 'materialized_view'
    WHEN 'i' THEN 'index'
    WHEN 'S' THEN 'sequence'
    WHEN 's' THEN 'special'
    WHEN 'f' THEN 'foreign_table'
  END as "type",
  pg_catalog.pg_get_userbyid(c.relowner) as "owner"
FROM
  pg_catalog.pg_class c
LEFT JOIN
  pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE
  c.relkind IN ('r','v','m','S','s','') AND
  n.nspname !~ '^pg_toast' AND 
  n.nspname NOT IN ('information_schema', 'pg_catalog')
ORDER BY 1, 2`
)
