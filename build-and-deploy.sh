#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o sqrt/sqrt sqrt/sqrt.go
zip sqrt/sqrt.zip sqrt/sqrt
sls deploy