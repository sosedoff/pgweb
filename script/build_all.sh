#!/usr/bin/env bash

TARGETS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"
ARM_TARGETS="arm/v5 arm64/v7"

for target in $TARGETS; do
  echo "-> target: $target"

  parts=(${target//\// })
  os=${parts[0]}
  arch=${parts[1]}

  GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o "./bin/pgweb_${os}_${arch}"
done

for target in $ARM_TARGETS; do
  echo "-> target: $target"

  parts=(${target//\// })
  arch=${parts[0]}
  arm=$(echo ${parts[1]} | sed s/v//g)

  GOOS=linux GOARCH=$arch GOARM=$arm go build -ldflags "$LDFLAGS" -o "./bin/pgweb_linux_${arch}_v${arm}"
done
