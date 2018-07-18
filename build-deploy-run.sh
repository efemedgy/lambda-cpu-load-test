#!/bin/bash
mSize=$1

GOOS=linux GOARCH=amd64 go build -o sqrt/sqrt sqrt/sqrt.go
GOOS=linux GOARCH=amd64 go build -o sleep/sleep sleep/sleep.go
GOOS=linux GOARCH=amd64 go build -o goroutine/goroutine goroutine/goroutine.go

zip cpuload.zip sqrt/sqrt sleep/sleep goroutine/goroutine
sls deploy --mSize $mSize

go build -o invoker-test invoker/invoker.go
./invoker-test sqrt $mSize
./invoker-test sleep $mSize
./invoker-test goroutine $mSize