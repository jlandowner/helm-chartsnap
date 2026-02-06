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
	# For development: install directly from GitHub (requires --verify=false for Helm 4)
	# For production: use release tarballs from https://github.com/jlandowner/helm-chartsnap/releases
	helm plugin install https://github.com/jlandowner/helm-chartsnap --verify=false
