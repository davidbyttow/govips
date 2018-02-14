#!/bin/bash
if [ -f pkg/vips/operators.go ]; then
  rm pkg/vips/operators.go
fi
python scripts/gen-operators.py >> pkg/vips/operators.go
gofmt -s -w pkg/vips/operators.go
