#! /bin/sh
set -euo pipefail

export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0
export GO111MODULE="on"

mkdir -p dist

echo "Downloading Go modules…"
go mod download

echo "Building linux binary…"
go build -v -o dist/go-mysql-crud

echo "Building docker image…"
docker build -t go-mysql-crud .
