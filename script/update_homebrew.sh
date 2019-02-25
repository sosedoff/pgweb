#!/bin/bash

RELEASE_FILE="./tmp/release.json"
HOMEBREW_ROOT="/usr/local/Homebrew/Library/Taps/homebrew/homebrew-core"
export HOMEBREW_GITHUB_API_TOKEN=$(awk '/api.github.com/{getline;getline;print $2}' ~/.netrc)

# Setup directory
mkdir -p ./tmp
rm -rf ./tmp/*

# Fetch the latest published version
curl -s https://api.github.com/repos/sosedoff/pgweb/releases/latest > $RELEASE_FILE
VERSION="$(jq -r .tag_name < $RELEASE_FILE)"
URL="https://github.com/sosedoff/pgweb/archive/$VERSION.tar.gz"
URL_SHA256=$(wget -qO- $URL | shasum -a 256 | cut -d ' ' -f 1)

# Reset any changes
git -C $HOMEBREW_ROOT reset --hard

# Update formula
brew bump-formula-pr \
  --url=$URL \
  --sha256=$URL_SHA256 \
  pgweb
