#!/bin/bash
if [ -f operators.go ]; then
  rm operators.go
fi
python scripts/gen-operators.py >> operators.go
gofmt -s -w operators.go
