#!/usr/bin/env sh
set -euox pipefail

go test ./... -short $@
