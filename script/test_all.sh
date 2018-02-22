#!/bin/bash
#
# Integartion testing with dockerized Postgres servers
#
# Boot2Docker is deprecated and no longer supported.
# Requires Docker for Mac to run on OSX.
# Install: https://docs.docker.com/engine/installation/mac/
#

set -e

export PGHOST=${PGHOST:-localhost}
export PGUSER="postgres"
export PGPASSWORD=""
export PGDATABASE="booktown"
export PGPORT="15432"

versions="9.1 9.2 9.3 9.4 9.5 9.6 10 10.1 10.2"

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