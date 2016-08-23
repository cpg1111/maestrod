#!/bin/bash

PORT_PREFIX=''

if [ ! -z "$1" ]; then
    mkdir -p /tmp/etcd2/src/ /tmp/etcd2/pkg/ /tmp/etcd2/bin/
    GOPATH=/tmp/etcd2/
    cd $GOPATH
    go get -u github.com/coreos/etcd
    cd $GOPATH/src/github.com/coreos/etcd
    git checkout $1
    PORT_PREFIX=2
else
    mkdir -p /tmp/etcd3/src/ /tmp/etcd3/pkg/ /tmp/etcd3/bin/
    GOPATH=/tmp/etcd3/
    cd $GOPATH
    go get -u github.com/coreos/etcd
    cd $GOPATH/src/github.com/coreos/etcd
    git checkout d53923c636e0e4ab7f00cb75681b97a8f11f5a9d
    PORT_PREFIX=3
fi

./build
./bin/etcd --listen-peer-urls="http://0.0.0.0:`echo $PORT_PREFIX`2380,http://0.0.0.0:`echo $PORT_PREFIX`7001" --listen-client-urls="http://0.0.0.0:`echo $PORT_PREFIX`2379,http://0.0.0.0:`echo $PORT_PREFIX`4001" --advertise-client-urls="http://0.0.0.0:`echo $PORT_PREFIX`2379,http://0.0.0.0:`echo $PORT_PREFIX`4001"
