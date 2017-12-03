#!/bin/bash
#
# Integartion testing with dockerized Postgres servers
#
# Boot2Docker is deprecated and no longer supported.
# Requires Docker for Mac to run on OSX.
# Install: https://docs.docker.com/engine/installation/mac/
#

set -e

export PGWEB_VERSION=0.9.10
export PGWEB_PORT=8081
export PGHOST=${PGHOST:-localhost}
export PGUSER="postgres"
export PGPASSWORD="postgres"
export PGDATABASE="booktown"
export PGPORT="15432"


export PGVERSION="9.6"

echo "---------------- BEGIN TEST ----------------"
echo "Running acceptance tests against PostgreSQL v$PGVERSION"

docker rm -f postgres || true
docker run -p $PGPORT:5432 --name postgres -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres:$PGVERSION
echo " done"

sleep 5


killall pgweb || true
./pgweb -s 2>&1 > /dev/null &

sleep 5

ginkgo ./spec/...
echo "---------------- END TEST ------------------"
