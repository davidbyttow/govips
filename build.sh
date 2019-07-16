#!/bin/bash

set -e

CGO_CFLAGS_ALLOW=-Xpreprocessor go build ./vips
