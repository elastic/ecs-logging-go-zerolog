SHELL=/usr/bin/env bash

# Directory to dump build tools into
GOBIN=$(shell go env GOPATH)/bin/

.PHONY: check
check:  ## Check copyright headers, format, and tidy the mod file
	@env GOBIN=${GOBIN} go install github.com/elastic/go-licenser@latest
	@env PATH="${GOBIN}:${PATH}" go-licenser
	@env go fmt
	@env go mod tidy

.PHONY: test
test:
	go test -v -race ./...

