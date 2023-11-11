all: build

GO ?= go1.21.3

go:
	go install golang.org/dl/$(GO)@latest
	rm -f $$(which $(GO)))/go
	ln -s $$(which $(GO)) $$(dirname $$(which $(GO)))/go

helm:
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

goreleaser:
	$(GO) install github.com/goreleaser/goreleaser@latest

.PHONY: build
build: goreleaser
	goreleaser build --single-target --snapshot --clean --skip=before

.PHONY: test
test:
	$(GO) test ./... -cover cover.out

.PHONY: update-versions
update-versions:
	sed -i.bk 's/version: .*/version: $(VERSION)/' plugin.yaml

.PHONY: debug-plugin
debug-plugin: build
	cp ./dist/chartsnap_*/chartsnap ~/.local/share/helm/plugins/helm-chartsnap/bin/
