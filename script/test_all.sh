#!/bin/bash

set -e

export PGHOST=${PGHOST:-192.168.99.100}
export PGUSER="postgres"
export PGPASSWORD=""
export PGDATABASE="booktown"
export PGPORT="15432"

for i in {1..5}
do
  export PGVERSION="9.$i"
  echo "Running tests against PostgreSQL v$PGVERSION"
  docker rm -f postgres || true
  docker run -p $PGPORT:5432 --name postgres -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres:$PGVERSION
  sleep 5
  make test
  echo "----------"
done