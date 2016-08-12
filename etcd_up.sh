#!/bin/bash

if [ -z "$GOPATH" ]; then
    echo "no gopath set"
    exit 1
fi

cd $GOPATH

go get -u github.com/mattn/goreman
go get -u github.com/coreos/etcd

cd $GOPATH/src/github.com/coreos/etcd

if [ ! -z "$1" ]; then
    git checkout $1
else
    git checkout d53923c636e0e4ab7f00cb75681b97a8f11f5a9d
fi

./build
goreman start

