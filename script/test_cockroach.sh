#!/bin/bash

set -e

function killproc() {
  if [[ $(lsof -i tcp:8888) ]]; then
    lsof -i tcp:8888 | grep pgweb | awk '{print $2}' | xargs kill
  fi
}

# Nuke the old container if exists.
docker rm -f cockroach || true

# Start cockroach on 26258 so we dont mess with local server.
docker run \
  --name=cockroach \
  -d \
  -t \
  -p 26258:26257 \
  cockroachdb/cockroach:v20.2.5 \
  start-single-node --insecure

sleep 3

# Load the demo database.
docker exec -i cockroach ./cockroach sql --insecure < ./data/roach.sql

# Find and destroy the existing pgweb process.
# Would be great if pgweb had --pid option.
killproc

# Start pgweb and connect to cockroach.
make build

./pgweb \
  --url=postgres://root@localhost:26258/roach?sslmode=disable \
  --listen=8888 \
  --skip-open &

sleep 1

# Run smoke tests
base="-w \"\n\" -f http://localhost:8888/api"
table="product_information"

curl $base/info
curl $base/connection
curl $base/schemas
curl $base/objects
curl $base/query -F query='select * from product_information;'
curl $base/tables/$table
curl $base/tables/$table/rows
curl $base/tables/$table/info
curl $base/tables/$table/indexes
curl $base/tables/$table/constraints

# Cleanup
killproc