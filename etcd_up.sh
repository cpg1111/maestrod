#!/bin/bash

if [ -z "$GOPATH" ]; then
    echo "no gopath set"
    exit 1
fi

cd $GOPATH

go get -u github.com/mattn/goreman
go get -u github.com/coreos/etcd

cd $GOPATH/src/github.com/coreos/etcd

PORT_PREFIX=''

if [ ! -z "$1" ]; then
    git checkout $1
    PORT_PREFIX=2
else
    git checkout d53923c636e0e4ab7f00cb75681b97a8f11f5a9d
    PORT_PREFIX=3
fi

./build
./bin/etcd --listen-peer-urls="0.0.0.0:`echo $PORT_PREFIX`2379,0.0.0.0:`echo $PORT_PREFIX`4001" --listen-client-urls="0.0.0.0:`echo $PORT_PREFIX`2380,0.0.0.0:`echo $PORT_PREFIX`7001"
