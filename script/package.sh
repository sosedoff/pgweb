#!/bin/bash

set -e

DIR="./bin"
rm $DIR/*.zip

for file in $(ls $DIR)
do
  fin=$DIR/$file
  fout=$DIR/$file.zip
  shasum -a 256 $fin
  zip -9 -q $fout $fin
done