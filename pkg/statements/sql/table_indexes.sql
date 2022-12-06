SELECT
  indexname AS index_name,
  pg_size_pretty(pg_table_size((schemaname || '.' || indexname)::regclass)) AS index_size,
  indexdef AS index_definition
FROM
  pg_indexes
WHERE
  schemaname = $1
  AND tablename = $2
