SELECT
  column_name,
  data_type,
  is_nullable,
  character_maximum_length,
  character_set_catalog,
  column_default,
  pg_catalog.col_description(('"' || $1::text || '"."' || $2::text || '"')::regclass::oid, ordinal_position) as comment
FROM
  information_schema.columns
WHERE
  table_schema = $1
  AND table_name = $2
