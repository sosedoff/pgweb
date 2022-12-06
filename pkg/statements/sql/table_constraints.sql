SELECT
  conname AS name,
  pg_get_constraintdef(c.oid, true) AS definition
FROM
  pg_constraint c
JOIN
  pg_namespace n ON n.oid = c.connamespace
JOIN
  pg_class cl ON cl.oid = c.conrelid
WHERE
  n.nspname = $1
  AND relname = $2
ORDER BY
  contype DESC
