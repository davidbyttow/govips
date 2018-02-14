#!/bin/bash

set -e

go generate ./pkg/vips
go build ./pkg/vips
