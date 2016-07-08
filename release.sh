#!/bin/bash

rm -rf dist/*
docker build -t maestrod_build -f Dockerfile_build .
docker run --rm -v `pwd`/dist/:/opt/bin/ maestrod_build
