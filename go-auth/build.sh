#!/bin/sh

docker run --rm -it \
    -v $PWD:/go/src/github.com/kwangyoung/go-messenger/go-auth/ \
    -e SOURCE_PATH=github.com/kwangyoung/go-messenger/go-auth/ \
    usman/go-builder:1.4

docker build -t usman/go-auth .
