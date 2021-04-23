#!/usr/bin/env bash
set -e

GOOS=linux GOARCH=amd64 go build -o pp

docker build -t=pp:1.0.0 ./
rm -f pp
