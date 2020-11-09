all: deps build test

deps:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go get ./...

build:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build ./vips

test:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go test -v ./...
