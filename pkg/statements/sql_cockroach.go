package statements

var (
	
	// ---------------------------------------------------------------------------
	DatabasesCockroach = `
		SELECT
  			datname
		FROM
  			pg_database
		WHERE
  			NOT datistemplate
		ORDER BY
  			datname ASC`
  			
	// ---------------------------------------------------------------------------
	SchemasCockroach = `
		SELECT
  			schema_name
		FROM
  			information_schema.schemata
		ORDER BY
  			schema_name ASC`
	
	// ---------------------------------------------------------------------------
	InfoCockroach = `
		SELECT
  			session_user,
  			current_user,
  			current_database(),
  			current_schemas(false),
  			'',
  			'',
  			'',
  			'',
  			version()`
	// ---------------------------------------------------------------------------
	TableIndexesCockroach = "SHOW INDEX FROM $1.$2";
	// ---------------------------------------------------------------------------
	TableConstraintsCockroach = `SHOW CONSTRAINTS FROM $1.$2`
	
	// ---------------------------------------------------------------------------
	TableInfoCockroach = `
		SELECT
			0 AS data_size,
  	   		0 AS index_size,
			0 AS total_size,
			( SELECT count(*) FROM $1 AS rows_count )`
	// ---------------------------------------------------------------------------
	TableSchemaCockroach = `
		SELECT
  			column_name,
  			data_type,
  			is_nullable,
  			character_maximum_length,
  			column_default
		FROM
  			information_schema.columns
		WHERE
  			table_schema = $1 AND
  			table_name = $2`
	// ---------------------------------------------------------------------------
	MaterializedViewCockroach = `
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
	ObjectsCockroach = `
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
  			pg_catalog.pg_get_userbyid(c.relowner) as "owner",
  			pg_catalog.obj_description(c.oid) as "comment"
  		FROM
  			pg_catalog.pg_class c
		LEFT JOIN
  			pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE
  			c.relkind IN ('r','v','m','S','s','')
  			AND (
  				  ( current_user() = 'root' and n.nspname NOT IN ('system', 'crdb_internal') ) OR
  				  ( current_user() <> 'root' and n.nspname NOT IN ('information_schema', 'pg_catalog','system', 'crdb_internal') )
  				)
		ORDER BY 1, 2`
		
	// ---------------------------------------------------------------------------
	ActivityCockroach = map[string]string{
		"default": `SHOW CLUSTER QUERIES`,
	}


	CockroachDialect = &DatabaseDialect {
&DatabasesCockroach,
&SchemasCockroach,
&InfoCockroach,
&TableIndexesCockroach,
&TableConstraintsCockroach,
&TableInfoCockroach,
&TableSchemaCockroach,
&MaterializedViewCockroach,
&ObjectsCockroach,
&ActivityCockroach,
"cockroach",
}

)

func init() {
	RegisterDialect("cockroach",CockroachDialect )
}
