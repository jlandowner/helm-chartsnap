#!/bin/bash
set -e

# This script generates Helm-compatible .prov (provenance) files for plugin archives
# It must be run after goreleaser creates the archives and signatures

DIST_DIR="${1:-dist}"
GPG_FINGERPRINT="${2:-${GPG_FINGERPRINT}}"
GPG_PASSPHRASE="${3:-${GPG_PASSPHRASE}}"

if [ -z "$GPG_FINGERPRINT" ]; then
    echo "Error: GPG_FINGERPRINT is required"
    echo "Usage: $0 <dist_dir> <gpg_fingerprint> <gpg_passphrase>"
    exit 1
fi

echo "Generating .prov files for archives in $DIST_DIR"

# Find all tar.gz archives (excluding checksums file)
for archive in "$DIST_DIR"/*.tar.gz; do
    # Skip if no archives found
    [ -e "$archive" ] || continue
    
    # Skip checksums file
    if [[ "$archive" == *"checksums"* ]]; then
        continue
    fi
    
    filename=$(basename "$archive")
    prov_file="${archive}.prov"
    
    echo "Processing $filename..."
    
    # Calculate SHA256 hash
    if command -v sha256sum >/dev/null 2>&1; then
        hash=$(sha256sum "$archive" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        hash=$(shasum -a 256 "$archive" | awk '{print $1}')
    else
        echo "Error: Neither sha256sum nor shasum found"
        exit 1
    fi
    
    # Create the provenance file content
    cat > "${prov_file}.unsigned" <<EOF
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

name: chartsnap
version: ${GITHUB_REF_NAME#v}
description: Snapshot testing for Helm charts
home: https://github.com/jlandowner/helm-chartsnap

...
files:
  $filename: sha256:$hash
EOF
    
    # Sign the provenance file with GPG
    if [ -n "$GPG_PASSPHRASE" ]; then
        gpg --batch --yes \
            --passphrase "$GPG_PASSPHRASE" \
            --pinentry-mode loopback \
            --local-user "$GPG_FINGERPRINT" \
            --clearsign \
            --output "$prov_file" \
            "${prov_file}.unsigned"
    else
        gpg --batch --yes \
            --local-user "$GPG_FINGERPRINT" \
            --clearsign \
            --output "$prov_file" \
            "${prov_file}.unsigned"
    fi
    
    # Clean up unsigned file
    rm -f "${prov_file}.unsigned"
    
    echo "Created $prov_file"
done

echo "All .prov files generated successfully"
