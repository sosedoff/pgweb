#!/usr/bin/env bash

if grep -H -r "<<<<<<< HEAD" ./pkg/
then
  echo "Merge conflict markers detected"
  exit 1
fi