#!/bin/bash

if grep -q 'go/src/github.com/sosedoff/pgweb' ./pkg/data/bindata.go
then
  echo "=========================================================="
  echo "ERROR: Bindata contains development references to assets!"
  echo "Fix with 'make assets' and commit the change."
  echo "=========================================================="
  exit 1
fi