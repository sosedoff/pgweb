WITH rel AS (
  SELECT format('%I.%I', $1::text, $2::text)::regclass AS oid
),
pk AS (
  SELECT
    a.attname AS column_name,
    k.ord     AS key_position
  FROM pg_index i
  CROSS JOIN LATERAL unnest(i.indkey) WITH ORDINALITY AS k(attnum, ord)
  JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = k.attnum
  WHERE i.indrelid = (SELECT oid FROM rel)
    AND i.indisprimary
)
SELECT
  c.column_name,
  c.data_type,
  c.udt_name,
  c.is_nullable::text = 'YES'  AS is_nullable,
  c.column_default::text,
  c.ordinal_position,
  c.is_identity::text = 'YES'  AS is_identity,
  c.is_generated::text <> 'NEVER' AS is_generated,
  c.is_updatable::text = 'YES' AS is_updatable,
  COALESCE(pk.key_position, 0) AS key_position,
  (SELECT relkind::text FROM pg_class WHERE oid = (SELECT oid FROM rel)) AS relkind
FROM information_schema.columns c
LEFT JOIN pk ON pk.column_name = c.column_name
WHERE c.table_schema = $1
  AND c.table_name = $2
ORDER BY c.ordinal_position
