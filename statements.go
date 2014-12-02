package main

const (
	PG_DATABASES = `SELECT datname FROM pg_database WHERE NOT datistemplate ORDER BY datname ASC`

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

	PG_TABLE_INDEXES = `SELECT
  indexname
, indexdef
FROM pg_class, pg_indexes, pg_namespace
WHERE pg_class.oid = $1::regclass
  AND pg_class.relname = pg_indexes.tablename
  AND pg_class.relnamespace = pg_namespace.oid
  AND pg_namespace.nspname = pg_indexes.schemaname`

	PG_TABLE_INFO = `SELECT
  pg_size_pretty(pg_table_size($1)) AS data_size
, pg_size_pretty(pg_indexes_size($1)) AS index_size
, pg_size_pretty(pg_total_relation_size($1)) AS total_size
, (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count`

	PG_TABLE_SCHEMA = `SELECT
  column_name
, data_type
, is_nullable
, character_maximum_length
, character_set_catalog
, column_default
FROM information_schema.columns, pg_class, pg_namespace
WHERE pg_class.oid = $1::regclass
  AND pg_class.relname = information_schema.columns.table_name
  AND pg_class.relnamespace = pg_namespace.oid
  AND pg_namespace.nspname = information_schema.columns.table_schema`

	PG_TABLES = `SELECT
  pg_class.oid::regclass
FROM pg_class, pg_namespace
WHERE pg_class.relkind IN ('m', 'r', 'v')
  AND pg_class.relnamespace = pg_namespace.oid
  AND pg_namespace.nspname = ANY (current_schemas(false))`
)
