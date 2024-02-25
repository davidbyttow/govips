#!/bin/sh

docker build -t govips-ubuntu:23.10 -f Dockerfile-ubuntu-23.10 ./ || (echo build failed 1>&2 && exit 1)

mkdir -p /tmp/volume

echo To run the container:
echo docker run --rm -it -v /tmp/volume:/volume govips-ubuntu:23.10
