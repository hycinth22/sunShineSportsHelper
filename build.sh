#!/bin/sh

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o build/jkwx

cp -i build_assets/* build/