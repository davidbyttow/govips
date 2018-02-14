#!/bin/bash

mkdir -p cmd/examples/build

echo "Running avg example"
go run cmd/examples/avg/avg.go -file assets/fixtures/canyon.jpg

echo "Running embed example"
go run cmd/examples/embed/embed.go -in assets/fixtures/canyon.jpg -out cmd/examples/build/embed-canyon.jpg

echo "Running invert example"
go run cmd/examples/invert/invert.go -in assets/fixtures/canyon.jpg -out cmd/examples/build/invert-canyon.jpg

echo "Running resize example"
go run cmd/examples/resize/resize.go -in assets/fixtures/canyon.jpg -out cmd/examples/build/resize-canyon.jpg
