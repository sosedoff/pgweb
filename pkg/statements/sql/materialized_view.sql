SELECT
  attname AS column_name,
  atttypid::regtype AS data_type,
  (CASE WHEN attnotnull IS TRUE THEN 'NO' ELSE 'YES' END) AS is_nullable,
  NULL AS character_maximum_length,
  NULL AS character_set_catalog,
  NULL AS column_default
FROM
  pg_attribute
WHERE
  attrelid = $1::regclass
  AND attnum > 0
  AND NOT attisdropped
