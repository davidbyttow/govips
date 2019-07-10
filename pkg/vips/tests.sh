#!/usr/bin/env bash

# We run the tests separately because when running them all together gives these kinds of errors:
#
# duplicate symbol _gobject_set_property in:
#    /var/folders/jn/vwmxflgd3cg4fc88wvxphmh86qgkln/T/go-link-092501829/000010.o
#    /var/folders/jn/vwmxflgd3cg4fc88wvxphmh86qgkln/T/go-link-092501829/000030.o
#
# And there is one such error for each symbol defined in bridge.c
# Possibly this issue: https://github.com/golang/go/issues/32150
# Last tested with go 1.12.7, still happens.

export CGO_CFLAGS_ALLOW=-Xpreprocessor
for file in *_test.go
do
  go test -v ${file}
done
