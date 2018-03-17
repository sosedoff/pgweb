#!/bin/bash

files="$(go fmt ./pkg/...)"
ignore="pkg/data/bindata.go"
files=${files[@]/$ignore}

if [ -n "$files" ]; then
  echo "Go code is not formatted: $files"
  for file in $files; do
    echo "----> $file"
  done
  exit 1
fi