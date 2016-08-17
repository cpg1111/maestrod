#!/bin/bash

INITCTL=`which systemctl`

if [ -z "$INITCTL" ]; then
    INITCTL=`which service`
fi

DOCKER_MACHINE=`which docker-machine`

if [ ! -z "$DOCKER_MACHINE" ]; then
    $DOCKER_MACHINE create -d=vmwarefusion maestrod-dev;
    eval "$($DOCKER_MACHINE env maestrod-dev)"
fi

echo $INITCTL | grep systemctl

if [ $? -eq 0 ]; then
    $INITCTL start docker
else
    $INITCTL docker start
fi

docker rm -f $(docker ps -a -q)
docker run --rm -p 27017:27017 -p 28017:28027 mongo &
docker run --rm -p 6379:6379 redis
