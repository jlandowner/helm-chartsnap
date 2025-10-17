#!/bin/bash

# This simulates helm template output with post-renderer that converts \n to actual newlines
# This is a problematic case for kyaml parser
cat <<EOF
---
# Source: test-chart/templates/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
stringData:
  certificate: "-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKZ
-----END CERTIFICATE-----"
---
# Source: test-chart/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
data:
  message: "hello world"
EOF
