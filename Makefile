all: build

GO ?= go1.22.1

go:
	-go install golang.org/dl/$(GO)@latest
	$(GO) download
	rm -f $$(dirname $$(which $(GO)))/go
	ln -s $$(which $(GO)) $$(dirname $$(which $(GO)))/go
	go version

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

.PHONY: integ-test
integ-test: debug-plugin
	-helm chartsnap --chart example/app1 $(ARGS)
	-helm chartsnap --chart example/app1 -f example/app1/test/test_ingress_enabled.yaml $(ARGS)
	-helm chartsnap --chart example/app1 -f example/app1/test/ $(ARGS)

.PHONY: update-versions
update-versions:
	sed -i.bk 's/version: .*/version: $(VERSION)/' plugin.yaml

.PHONY: debug-plugin
debug-plugin: build
	-helm plugin install https://github.com/jlandowner/helm-chartsnap
	cp ./dist/chartsnap_*/chartsnap ~/.local/share/helm/plugins/helm-chartsnap/bin/

.PHONY: snap-helm-template-help
snap-helm-template-help:
	cd hack/helm-template-help-snapshot; $(GO) run main.go
