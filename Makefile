all: deps build test

.PHONY: deps
deps:
	CGO_CFLAGS_ALLOW='-I.*' go get -v -t ./...

.PHONY: build
build:
	CGO_CFLAGS_ALLOW='-I.*' go build -v ./vips

.PHONY: test
test:
	CGO_CFLAGS_ALLOW='-I.*' go test -v -coverprofile=profile.cov ./...

.PHONY: clean
clean:
	go clean

.PHONY: clean-cache
clean-cache:
	# Purge build cache and test cache.
	# When something went wrong on building or testing, try this.
	-go clean -testcache
	-go clean -cache

.PHONY: distclean
distclean:
	-go clean -testcache
	-go clean -cache
	-git clean -f -x
