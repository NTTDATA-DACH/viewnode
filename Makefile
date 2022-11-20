VERSION = v0.8.3

.PHONY: clean build test run install all release

clean:
	@go clean

build:
	@go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=`git rev-parse HEAD`"

test:
	# use 'atomic' cover mode for safe concurrent test (default is 'set', we can use also 'count')
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...

run:
	@go run -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=`git rev-parse HEAD`" main.go $(cmd)

install:
	@go install -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=`git rev-parse HEAD`"

all: clean install

release:
	@goreleaser --snapshot --skip-publish --rm-dist