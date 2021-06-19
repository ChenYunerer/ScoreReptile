#!/bin/bash
rm -rf soupu_admin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o soupu_admin ./src/main
