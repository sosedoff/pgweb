#!/bin/bash

files="$(ls ./bin/ | grep zip | grep amd64)"
bucket="s3://pgweb-github-builds"

for file in $files; do
  echo "Uploading $file"
  aws s3 cp ./bin/$file $bucket/$TRAVIS_BRANCH/$file
done