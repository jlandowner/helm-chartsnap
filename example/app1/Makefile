all: dir

OPT ?= # --debug

.PHONY: file
file:
	-helm chartsnap $(OPT) -c . -f test/test_ingress_enabled.yaml

.PHONY: dir
dir:
	-helm chartsnap $(OPT) -c . -f test/

.PHONY: update
update:
	-helm chartsnap $(OPT) -c . -f test/ -u

helm:
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

helm-chartsnap:
	helm plugin install https://github.com/jlandowner/helm-chartsnap --verify=false
