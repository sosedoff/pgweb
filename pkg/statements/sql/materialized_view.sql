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
  NOT attisdropped
