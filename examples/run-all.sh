#!/bin/bash

mkdir -p examples/build

echo "Running avg example"
go run examples/avg/avg.go -file fixtures/canyon.jpg

echo "Running embed example"
go run examples/embed/embed.go -in fixtures/canyon.jpg -out examples/build/embed-canyon.jpg

echo "Running invert example"
go run examples/invert/invert.go -in fixtures/canyon.jpg -out examples/build/invert-canyon.jpg

echo "Running resize example"
go run examples/resize/resize.go -in fixtures/canyon.jpg -out examples/build/resize-canyon.jpg
