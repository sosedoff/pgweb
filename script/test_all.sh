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
export PGPASSWORD=""
export PGDATABASE="booktown"
export PGPORT="15432"

versions="9.1 9.2 9.3 9.4 9.5 9.6 10.0 10.1 10.2 10.3"

for i in $versions
do
  export PGVERSION="$i"
  echo "----------------------- PostgreSQL v$PGVERSION ------------------------"
  echo "Removing existing container"
  docker rm -f postgres || true 
  echo "Starting contaienr for $PGVERSION"
  docker run -p $PGPORT:5432 --name postgres -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres:$PGVERSION
  sleep 5
  make test
done