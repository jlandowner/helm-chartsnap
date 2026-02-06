all: ingress-nginx cilium nginx-gateway-fabric

OPT ?= # -u

.PHONY: ingress-nginx
ingress-nginx:
	-helm chartsnap $(OPT) -c ingress-nginx -f ingress-nginx.values.yaml -- --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx --skip-tests

.PHONY: cilium
cilium:
	-helm chartsnap $(OPT) -c cilium -f cilium.values.yaml -- --repo https://helm.cilium.io --namespace kube-system

.PHONY: nginx-gateway-fabric
nginx-gateway-fabric:
	-helm chartsnap $(OPT) -c oci://ghcr.io/nginxinc/charts/nginx-gateway-fabric -n nginx-gateway -f nginx-gateway-fabric.values.yaml

helm:
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

helm-chartsnap:
	# For development: install directly from GitHub (requires --verify=false for Helm 4)
	# For production: use release tarballs from https://github.com/jlandowner/helm-chartsnap/releases
	helm plugin install https://github.com/jlandowner/helm-chartsnap --verify=false
