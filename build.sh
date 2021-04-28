#!/bin/bash
rm -rf main
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./src/main
