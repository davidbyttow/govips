#!/bin/bash

set -e

go generate ./pkg/vips
CGO_CFLAGS_ALLOW=-Xpreprocessor go build ./pkg/vips
