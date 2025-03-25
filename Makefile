all: build

GO_VERSION ?= $(shell grep '^go ' go.mod | awk '{print $$2}')
GO ?= go$(GO_VERSION)

go:
	-go install golang.org/dl/go$(GO_VERSION)@latest
	go$(GO_VERSION) download

helm:
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

goreleaser:
	$(GO) install github.com/goreleaser/goreleaser@latest

.PHONY: build
build: goreleaser
	goreleaser build --single-target --snapshot --clean --skip=before

.PHONY: test
test:
	$(GO) test -race -coverprofile=coverage.txt -covermode=atomic `$(GO) list ./... | grep -v /hack`
	$(GO) tool cover -func=coverage.txt -o=coverage.out
	tail -1 coverage.out | awk '{gsub("%",""); print $$3}'

HELM_PLUGIN_PATH := $(shell helm env | grep HELM_PLUGINS | cut -d= -f2)

.PHONY: integ-test
integ-test: install-dev-bin
	helm chartsnap --chart example/app1 -f example/app1/test_latest/test_ingress_enabled.yaml --namespace default $(ARGS)
	helm chartsnap --chart example/app1 -f example/app1/test_latest/ --namespace default $(ARGS)
	helm chartsnap --chart example/app1 -f example/app1/test_v1/ --namespace default $(ARGS)
	helm chartsnap --chart example/app1 -f example/app1/test_v2/ --namespace default $(ARGS)
	helm chartsnap --chart example/app1 -f example/app1/test_v3/ --namespace default $(ARGS)
	helm chartsnap --chart oci://ghcr.io/nginxinc/charts/nginx-gateway-fabric -f example/remote/nginx-gateway-fabric.values.yaml $(ARGS) -- --namespace nginx-gateway $(EXTRA_ARGS)
	helm chartsnap --chart cilium -f example/remote/cilium.values.yaml $(ARGS) -- --namespace kube-system --repo https://helm.cilium.io $(EXTRA_ARGS)
	helm chartsnap --chart ingress-nginx -f example/remote/ingress-nginx.values.yaml $(ARGS) -- --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx --skip-tests $(EXTRA_ARGS)
	helm chartsnap --chart example/app2 --namespace default $(ARGS)
	helm chartsnap --chart example/app3 --namespace default $(ARGS)
	helm chartsnap --chart example/app3 --namespace default --snapshot-file-ext yaml $(ARGS)
	helm chartsnap --chart example/app3 --namespace default $(ARGS) -f example/app3/test/ok.yaml

.PHONY: integ-test-kong
integ-test-kong:
	cd hack/test; bash test-kong-chart.sh

.PHONY: integ-test-fail
integ-test-fail: install-dev-bin
	helm chartsnap --chart example/app1 --namespace default $(ARGS) && echo "should fail" && exit 1 || (echo "--- fail is expected ---"; true)
	helm chartsnap --chart example/app1 --namespace default -f example/app1/testfail/test_ingress_enabled.yaml $(ARGS) && echo "should fail" && exit 1 || (echo "--- fail is expected ---"; true)
	helm chartsnap --chart example/app1 --namespace default -f example/app1/testfail/ $(ARGS) && echo "should fail" && exit 1 || (echo "--- fail is expected ---"; true)

.PHONY: update-versions
update-versions:
	sed -i.bk 's/version: .*/version: $(VERSION)/' plugin.yaml

.PHONY: install-dev-bin
install-dev-bin: build
	-helm plugin install https://github.com/jlandowner/helm-chartsnap
	cp ./dist/chartsnap_*/chartsnap $(HELM_PLUGIN_PATH)/helm-chartsnap/bin/
	helm chartsnap --version

.PHONY: helm-template-help-snapshot
helm-template-help-snapshot:
	-rm hack/helm-template-help-snapshot/helm-template.snap
	cd hack/helm-template-help-snapshot; $(GO) run main.go

.PHONY: helm-template-diff
helm-template-diff:
	cd hack/helm-template-diff; $(GO) run main.go

.PHONY: helm-template-diff.update
helm-template-diff.update:
	-rm hack/helm-template-diff/helm-template.snap
	make helm-template-diff

kubectl-validate:
	$(GO) install sigs.k8s.io/kubectl-validate@latest

.PHONY: validate
validate: kubectl-validate
	kubectl validate example/remote/__snapshots__/
	kubectl validate example/app3/__snapshots__/ --local-crds hack/crd/
