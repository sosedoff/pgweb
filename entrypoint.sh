#!/bin/bash

# Wait for PostgreSQL
for i in {1..5}
do
    psql -Xq "$PGWEB_DATABASE_URL" -c "select 1" &>/dev/null && break
    sleep 5s
done

exec "$@"
