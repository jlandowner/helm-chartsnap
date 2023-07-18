all: build

GO ?= go
goreleaser:
	$(GO) install github.com/goreleaser/goreleaser@latest

.PHONY: build
build: goreleaser
	goreleaser build --single-target --snapshot --clean --skip-before

.PHONY: test
test:
	$(GO) test ./... -cover cover.out
