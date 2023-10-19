#!/bin/bash

# Wait for PostgreSQL
while ! psql -Xq "$PGWEB_DATABASE_URL" -c "select 1" &>/dev/null
do
    sleep 5s
done

exec "$@"
