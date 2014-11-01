package main

const (
	PG_INFO          = "SELECT version(), user, current_database(), inet_client_addr(), inet_client_port(), inet_server_addr(), inet_server_port()"
	PG_DATABASES     = "SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname ASC;"
	PG_TABLES        = "SELECT table_schema, table_name FROM information_schema.tables ORDER BY table_schema,table_name;"
	PG_TABLE_SCHEMA  = "SELECT column_name, data_type, is_nullable, character_maximum_length, character_set_catalog, column_default FROM information_schema.columns where table_schema = '%s' AND table_name = '%s';"
	PG_TABLE_INDEXES = "SELECT indexname, indexdef FROM pg_indexes WHERE schemaname = '%s' AND tablename = '%s';"
	PG_TABLE_INFO    = "SELECT pg_size_pretty(pg_table_size('%s')) AS data_size, pg_size_pretty(pg_indexes_size('%s')) AS index_size, pg_size_pretty(pg_total_relation_size('%s')) AS total_size, (SELECT COUNT(*) FROM %s) AS rows_count"
)
