SELECT
  session_user,
  current_user,
  current_database(),
  current_schemas(false),
  version()
