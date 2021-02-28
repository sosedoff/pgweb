#!/bin/bash

# Get list of offending files
files="$(go fmt ./pkg/...)"

if [ -n "$files" ]; then
  echo "Go code is not formatted:"
  for file in $files; do
    echo "----> $file"
  done
  exit 1
fi