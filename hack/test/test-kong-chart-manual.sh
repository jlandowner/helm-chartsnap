#!/bin/bash

# Update submodule to the latest commit
git submodule update --init --recursive --remote --merge

# Update charts repo to the latest main
cd charts
git pull origin main --rebase

# ----------------------------------
# Same test in kong golden tests
# ----------------------------------
helm chartsnap \
                -c ./charts/kong \
                -f ./charts/kong/ci/ \
                 \
                -- \
                --api-versions cert-manager.io/v1 \
                --api-versions gateway.networking.k8s.io/v1 \
                --api-versions gateway.networking.k8s.io/v1beta1 \
                --api-versions gateway.networking.k8s.io/v1alpha2 \
                --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicy \
                --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicyBinding

helm chartsnap \
                -c ./charts/ingress \
                -f ./charts/ingress/ci/ \
                 \
                -- \
                --api-versions cert-manager.io/v1 \
                --api-versions gateway.networking.k8s.io/v1 \
                --api-versions gateway.networking.k8s.io/v1beta1 \
                --api-versions gateway.networking.k8s.io/v1alpha2 \
                --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicy \
                --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicyBinding

helm chartsnap \
                -c ./charts/gateway-operator \
                -f ./charts/gateway-operator/ci/ \
                 \
                -- \
                --api-versions cert-manager.io/v1 \
                --api-versions gateway.networking.k8s.io/v1 \
                --api-versions gateway.networking.k8s.io/v1beta1 \
                --api-versions gateway.networking.k8s.io/v1alpha2 \
                --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicy \
                --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicyBinding