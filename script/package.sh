#!/bin/bash

set -e

DIR="./bin"
rm -f $DIR/*.zip

for file in $(ls $DIR)
do
  fin=$DIR/$file
  fout=$DIR/$file.zip
  shasum -a 256 $fin
  zip -9 -q -j $fout $fin
  shasum -a 256 $fout
done