package statements

const (
	PG_DATABASES = `SELECT datname FROM pg_database WHERE NOT datistemplate ORDER BY datname ASC`

	PG_SCHEMAS = `SELECT schema_name FROM information_schema.schemata ORDER BY schema_name ASC`

	PG_INFO = `SELECT
  session_user
, current_user
, current_database()
, current_schemas(false)
, inet_client_addr()
, inet_client_port()
, inet_server_addr()
, inet_server_port()
, version()`

	PG_TABLE_INDEXES = `SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1`

	PG_TABLE_INFO = `SELECT
  pg_size_pretty(pg_table_size($1)) AS data_size
, pg_size_pretty(pg_indexes_size($1)) AS index_size
, pg_size_pretty(pg_total_relation_size($1)) AS total_size
, (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count`

	PG_TABLE_SCHEMA = `SELECT
column_name, data_type, is_nullable, character_maximum_length, character_set_catalog, column_default
FROM information_schema.columns
WHERE table_name = $1`

	PG_TABLES = `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_schema,table_name`

	PG_ACTIVITY = `SELECT
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
  FROM pg_stat_activity
  WHERE state IS NOT NULL`
)
