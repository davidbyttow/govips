all: deps build test

deps:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go get ./...

test:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go test -v ./...

build:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build ./vips