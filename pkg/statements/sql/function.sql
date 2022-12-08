SELECT
  p.*,
  pg_get_functiondef(oid) AS functiondef
FROM
  pg_catalog.pg_proc p
WHERE
  oid = $1::oid
