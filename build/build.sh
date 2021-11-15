#!/bin/sh

docker build -t govips-ubuntu:18.04 -f Dockerfile-ubuntu-18.04 ./ || (echo build failed 1>&2 && exit 1)
docker build -t govips-ubuntu:20.04 -f Dockerfile-ubuntu-20.04 ./ || (echo build failed 1>&2 && exit 1)
docker build -t govips-ubuntu:20.10 -f Dockerfile-ubuntu-20.10 ./ || (echo build failed 1>&2 && exit 1)

mkdir -p /tmp/volume

echo To run the container:
echo docker run --rm -it -v /tmp/volume:/volume govips-ubuntu:18.04
echo docker run --rm -it -v /tmp/volume:/volume govips-ubuntu:20.04
echo docker run --rm -it -v /tmp/volume:/volume govips-ubuntu:20.10
