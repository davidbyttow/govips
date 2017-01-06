#!/bin/bash

mkdir build

echo "Running avg example"
go run avg/avg.go -file ../fixtures/canyon.jpg

echo "Running buffer example"
go run buffer/buffer.go -in ../fixtures/canyon.jpg -out build/buffer-canyon.jpg

echo "Running embed example"
go run embed/embed.go -in ../fixtures/canyon.jpg -out build/embed-canyon.jpg

echo "Running invert example"
go run invert/invert.go -in ../fixtures/canyon.jpg -out build/invert-canyon.jpg

echo "Running resize example"
go run resize/resize.go -in ../fixtures/canyon.jpg -out build/resize-canyon.jpg
