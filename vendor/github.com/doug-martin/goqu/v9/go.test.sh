#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -race -coverprofile=coverage.txt -coverpkg=./... ./...