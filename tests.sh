#!/usr/bin/env bash

# We run the tests separately because when running them all together gives these kinds of errors:
#
# And there is one such error for each symbol defined in bridge.c
# Possibly this issue: https://github.com/golang/go/issues/32150
# Last tested with go 1.12.7, still happens.

export CGO_CFLAGS_ALLOW=-Xpreprocessor

for file in ./vips/*_test.go
do
  go test -v ${file}
done
