package statements

import (
	_ "embed"
)

var (
	//go:embed sql/databases.sql
	Databases string

	//go:embed sql/schemas.sql
	Schemas string

	//go:embed sql/info.sql
	Info string

	//go:embed sql/info_simple.sql
	InfoSimple string

	//go:embed sql/estimated_row_count.sql
	EstimatedTableRowCount string

	//go:embed sql/table_indexes.sql
	TableIndexes string

	//go:embed sql/table_constraints.sql
	TableConstraints string

	//go:embed sql/table_info.sql
	TableInfo string

	//go:embed sql/table_info_cockroach.sql
	TableInfoCockroach string

	//go:embed sql/table_schema.sql
	TableSchema string

	//go:embed sql/materialized_view.sql
	MaterializedView string

	//go:embed sql/objects.sql
	Objects string

	//go:embed sql/tables_stats.sql
	TablesStats string

	//go:embed sql/function.sql
	Function string

	// Activity queries for specific PG versions
	Activity = map[string]string{
		"default": "SELECT * FROM pg_stat_activity WHERE datname = current_database()",
		"9.1":     "SELECT datname, current_query, waiting, query_start, procpid as pid, datid, application_name, client_addr FROM pg_stat_activity WHERE datname = current_database()",
		"9.2":     "SELECT datname, query, state, waiting, query_start, state_change, pid, datid, application_name, client_addr FROM pg_stat_activity WHERE datname = current_database()",
		"9.3":     "SELECT datname, query, state, waiting, query_start, state_change, pid, datid, application_name, client_addr FROM pg_stat_activity WHERE datname = current_database()",
		"9.4":     "SELECT datname, query, state, waiting, query_start, state_change, pid, datid, application_name, client_addr FROM pg_stat_activity WHERE datname = current_database()",
		"9.5":     "SELECT datname, query, state, waiting, query_start, state_change, pid, datid, application_name, client_addr FROM pg_stat_activity WHERE datname = current_database()",
		"9.6":     "SELECT datname, query, state, wait_event, wait_event_type, query_start, state_change, pid, datid, application_name, client_addr FROM pg_stat_activity WHERE datname = current_database()",
	}
)
