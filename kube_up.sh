#!/bin/bash

OS=`uname -s`
ARCH=`uname -m`

GET=`which curl`

if [ -z "$GET" ]; then
    GET=`which wget`
else
    GET="$GET -Lo ./minikube"
fi

get_kube() {
    echo "get"
    echo $GET
    if [ -z $OS ] || [ -z $ARCH  ]; then
        echo "could not determine OS or architecture"
        exit 1
    fi
    if [ "$ARCH" == "x86_64" ]; then
        if [ "$OS" == "Darwin" ]; then
            $GET https://storage.googleapis.com/minikube/releases/v0.7.1/minikube-darwin-amd64
        elif [ "$OS" == "Linux" ]; then
            $GET https://storage.googleapis.com/minikube/releases/v0.7.1/minikube-linux-amd64
        fi
    fi
}

mkdir -p ./dev_deps/ && cd ./dev_deps/ 

stat -q ./minikube
if [ $? -eq 1 ]; then
    get_kube && chmod 755 ./minikube
fi

./minikube $1 $2
cd -
