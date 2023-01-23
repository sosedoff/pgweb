WITH columns_counts AS (
  SELECT table_schema, table_name, COUNT(1) AS num
  FROM information_schema.columns
  GROUP BY table_schema, table_name
),
indexes_counts AS (
  SELECT schemaname, tablename, COUNT(1) AS num
  FROM pg_indexes
  GROUP BY schemaname, tablename
)
SELECT
  tables.schemaname AS schema_name,
  tables.relname AS table_name,
  pg_size_pretty(pg_total_relation_size(tables.relid)) AS total_size,
  pg_size_pretty(pg_relation_size(tables.relid)) AS data_size,
  pg_size_pretty(pg_indexes_size(tables.relid)) AS index_size,
  pg_class.reltuples AS estimated_rows_count,
  CASE
    WHEN pg_class.reltuples >= 0 AND pg_class.reltuples < 1000
      THEN pg_class.reltuples::text
    WHEN pg_class.reltuples >= 1000 AND pg_class.reltuples < 1000000
      THEN ROUND((pg_class.reltuples / 1000))::text || 'K'
    WHEN pg_class.reltuples >= 1000000
      THEN ROUND(pg_class.reltuples / 1000000)::text || 'M'
  END AS estimated_rows,
  CASE
    WHEN pg_class.reltuples > 1000
      THEN ROUND(pg_indexes_size(tables.relid)::numeric / pg_relation_size(tables.relid), 2)
  END AS index_to_data_ratio,
  indexes_counts.num AS indexes_count,
  columns_counts.num AS columns_count
FROM
  pg_catalog.pg_statio_user_tables AS tables
LEFT JOIN pg_class
  ON pg_class.oid = tables.relid
LEFT JOIN indexes_counts
  ON indexes_counts.schemaname = tables.schemaname
  AND indexes_counts.tablename = tables.relname
LEFT JOIN columns_counts
  ON columns_counts.table_schema = tables.schemaname
  AND columns_counts.table_name = tables.relname
ORDER BY
  pg_total_relation_size(tables.relid) DESC,
  pg_relation_size(tables.relid) DESC
