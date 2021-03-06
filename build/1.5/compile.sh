#!/bin/bash
set -e
# firewall block
mkdir -p $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://go.googlesource.com/crypto
# Set directory to where we expect code to be
cd /go/src/${SOURCE_PATH}
echo "Downloading dependencies"
godep restore
echo "Fix formatting"
go fmt ./...
echo "Running Tests"
go test ./... 
echo "Building source"
go build
echo "Build Successful"
