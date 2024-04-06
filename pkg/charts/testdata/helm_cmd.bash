#!/bin/bash

# Output arguments
echo "Arguments for helm: $@"

# Output environment variables starting with "HELM_"
echo "Environment variables starting with HELM_:"
env | grep '^HELM_'