# Example of snapshot testing for remote Helm repositories

## ingress-nginx

Docs: https://artifacthub.io/packages/helm/ingress-nginx/ingress-nginx

For example, install ingress-nginx which Service is bound to Network Load Balancer in Amazon EKS.

```yaml
# https://github.com/kubernetes/ingress-nginx/blob/main/charts/ingress-nginx/values.yaml#L518
controller:
  service:
    internal:
      enabled: true
      annotations:
        # Create internal NLB
        service.beta.kubernetes.io/aws-load-balancer-scheme: "internal"

```

Do snapshot 📸

```sh
helm chartsnap -c ingress-nginx -f ingress-nginx.values.yaml -- --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx
```

## cilium

Docs: https://docs.cilium.io/en/stable/installation/k8s-install-helm/

For example, install cilium as AWS ENI mode and enable Hubble UI.

```yaml
# https://docs.cilium.io/en/stable/installation/k8s-install-helm/
# EKS

eni:
  enabled: true

ipam:
  mode: eni

egressMasqueradeInterfaces: eth0

routingMode: native

# https://docs.cilium.io/en/stable/gettingstarted/hubble/#hubble-ui
# Enable Hubble UI
hubble:
  relay:
    enabled: true 
  ui:
    enabled: true
```

Do snapshot 📸

```sh
helm chartsnap -c cilium -f cilium.values.yaml -- --repo https://helm.cilium.io --namespace kube-system
```

However probably you will see the failure that does not match the snapshot of the certificate in Secrets.

Then add the followings to the value file and re-run the snapshot.

```yaml
# Change to fixed values
testSpec:
  dynamicFields:
    - apiVersion: v1
      kind: Secret
      name: cilium-ca
      jsonPath:
        - /data/ca.crt
        - /data/ca.key
    - apiVersion: v1
      kind: Secret
      name: hubble-relay-client-certs
      jsonPath:
        - /data/ca.crt
        - /data/tls.crt
        - /data/tls.key
    - apiVersion: v1
      kind: Secret
      name: hubble-server-certs
      jsonPath:
        - /data/ca.crt
        - /data/tls.crt
        - /data/tls.key
```