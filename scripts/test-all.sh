#!/usr/bin/env sh
set -euox pipefail

scripts/smoke.sh
go test pkg/vips/regression_test.go
go test pkg/vips/transform_test.go
