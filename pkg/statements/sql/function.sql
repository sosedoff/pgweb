SELECT
  p.oid,
  p.proname,
  p.pronamespace,
  p.proowner,
  pg_get_functiondef(oid) AS functiondef
FROM
  pg_catalog.pg_proc p
WHERE
  oid = $1::oid
