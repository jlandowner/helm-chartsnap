all: ingress-nginx cilium

.PHONY: ingress-nginx
ingress-nginx:
	helm chartsnap -c ingress-nginx -f ingress-nginx.values.yaml -- --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx

.PHONY: cilium
cilium:
	helm chartsnap -c cilium -f cilium.values.yaml -- --repo https://helm.cilium.io --namespace kube-system

helm:
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

helm-chartsnap:
	helm plugin install https://github.com/jlandowner/helm-chartsnap
