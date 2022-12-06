SELECT
  reltuples
FROM
  pg_class
WHERE
  oid = ('"' || $1::text || '"."' || $2::text || '"')::regclass
