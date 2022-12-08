SELECT
  pg_size_pretty(pg_table_size($1)) AS data_size,
  pg_size_pretty(pg_indexes_size($1)) AS index_size,
  pg_size_pretty(pg_total_relation_size($1)) AS total_size,
  (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count
