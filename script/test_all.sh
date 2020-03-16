#!/bin/bash
#
# Integartion testing with dockerized Postgres servers
#
# Requires Docker for Mac to run on OSX.
# Install: https://docs.docker.com/engine/installation/mac/
#

set -e

export PGHOST=${PGHOST:-localhost}
export PGUSER="postgres"
export PGPASSWORD="ci"
export PGDATABASE="booktown"
export PGPORT="15432"

# TODO: Enable the 10.x branch when it's supported on Travis.
# Local 10.x version is required so that pg_dump can properly work with older versions.
# 10.x branch is normally supported.
versions="9.1 9.2 9.3 9.4 9.5 9.6"

for i in $versions
do
  export PGVERSION="$i"

  echo "------------------------------- BEGIN TEST -------------------------------"
  echo "Running tests against PostgreSQL v$PGVERSION"
  docker rm -f postgres || true
  docker run -p $PGPORT:5432 --name postgres -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres:$PGVERSION
  sleep 5
  make test
  echo "-------------------------------- END TEST --------------------------------"
done

