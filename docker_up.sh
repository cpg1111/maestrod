#!/bin/bash

INITCTL=`which systemctl`

if [ -z "$INITCTL" ]; then
    INITCTL=`which service`
fi

DOCKER_MACHINE=`which docker-machine`

if [ ! -z "$DOCKER_MACHINE" ]; then
    $DOCKER_MACHINE create $1 maestrod-dev
    eval "$($DOCKER_MACHINE env maestrod-dev)"
    exit 0
fi

echo $INITCTL | grep systemctl

if [ $? -eq 0 ]; then
    $INITCTL start docker
else
    $INITCTL docker start
fi

