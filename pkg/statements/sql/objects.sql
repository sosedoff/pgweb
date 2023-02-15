WITH all_objects AS (
  SELECT
    c.oid,
    n.nspname AS schema,
    c.relname AS name,
    CASE c.relkind
      WHEN 'r' THEN 'table'
      WHEN 'v' THEN 'view'
      WHEN 'm' THEN 'materialized_view'
      WHEN 'i' THEN 'index'
      WHEN 'S' THEN 'sequence'
      WHEN 's' THEN 'special'
      WHEN 'f' THEN 'foreign_table'
    END AS type,
    pg_catalog.pg_get_userbyid(c.relowner) AS owner,
    pg_catalog.obj_description(c.oid) AS comment
  FROM
    pg_catalog.pg_class c
  LEFT JOIN
    pg_catalog.pg_namespace n ON n.oid = c.relnamespace
  WHERE
    c.relkind IN ('r','v','m','S','s','')
    AND n.nspname !~ '^pg_(toast|temp)'
    AND n.nspname NOT IN ('information_schema', 'pg_catalog')
    AND has_schema_privilege(n.nspname, 'USAGE')

  UNION

  SELECT
    p.oid,
    n.nspname AS schema,
    p.proname AS name,
    'function' AS function,
    pg_catalog.pg_get_userbyid(p.proowner) AS owner,
    NULL AS comment
  FROM
    pg_catalog.pg_namespace n
  JOIN
    pg_catalog.pg_proc p ON p.pronamespace = n.oid
  WHERE
    n.nspname !~ '^pg_(toast|temp)'
    AND n.nspname NOT IN ('information_schema', 'pg_catalog')
)
SELECT * FROM all_objects
ORDER BY 2, 3
