#!/bin/sh

MSYS_NO_PATHCONV=1 docker run --rm -e GOOS=js -e GOARCH=wasm -v "$GOPATH":/go -w /go/src/github.com/pichiw/pichiwui golang:1.11-rc-alpine3.7 go build -v -o app.wasm

MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 -v "$GOPATH":/go -w /go/src/github.com/pichiw/pichiwui golang:1.11-rc-alpine3.7 go build -v -a -installsuffix cgo -o server.bin github.com/pichiw/pichiwui/cmd/server

docker build .