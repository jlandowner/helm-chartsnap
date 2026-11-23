# Plugin Signing Implementation for Helm 4

This document describes the implementation of cryptographic signatures for helm-chartsnap releases.

## Implementation Status

✅ **COMPLETED** - Plugin signing is now fully implemented as of version 0.7.0.

## Overview

Helm 4 enforces plugin verification by default, requiring `.prov` provenance files alongside plugin archives. We have implemented automated signing in our release process to support this security feature.

## What Was Implemented

### 1. GoReleaser Configuration

Modified `.goreleaser.yml` to:
- Include `plugin.yaml` and `scripts/*` in release archives
- Configure GPG signing for all archive artifacts
- Generate `.sig` signature files
- Include `.prov` files in release assets

### 2. Provenance File Generation

Created `scripts/generate_prov_files.sh` to:
- Calculate SHA256 hashes of release archives
- Generate Helm-compatible `.prov` provenance files
- Sign provenance files with GPG
- Format files according to Helm's requirements

### 3. Release Workflow Updates

Modified `.github/workflows/release.yaml` to:
- Import GPG keys from GitHub secrets
- Run GoReleaser with signing enabled
- Generate `.prov` files after building
- Upload all signature files to the release

### 4. Documentation Updates

Updated documentation to:
- Explain signature support in README.md
- Document installation with signatures in docs/helm4-installation.md
- Provide setup instructions for maintainers

## How It Works

## How It Works

1. **Build Phase**: GoReleaser creates platform-specific archives containing:
   - Compiled `chartsnap` binary
   - `plugin.yaml` manifest
   - Installation scripts

2. **Signing Phase**: GPG signs each archive to create `.sig` files

3. **Provenance Generation**: A custom script generates `.prov` files containing:
   - Plugin metadata (name, version, description)
   - SHA256 hash of the archive
   - GPG signature (clearsigned format)

4. **Release Phase**: All artifacts are uploaded to GitHub:
   - `chartsnap_v*.tar.gz` (plugin archive)
   - `chartsnap_v*.tar.gz.sig` (detached signature)
   - `chartsnap_v*.tar.gz.prov` (Helm provenance file)

5. **Installation**: Helm 4 automatically:
   - Downloads both `.tar.gz` and `.tar.gz.prov`
   - Verifies the SHA256 hash matches
   - Validates the GPG signature
   - Installs only if verification succeeds

## Setup Requirements

### For Repository Maintainers

To enable signing, configure these GitHub repository secrets:

- **GPG_PRIVATE_KEY**: The GPG private key in ASCII armor format
  ```bash
  gpg --armor --export-secret-keys KEY_ID > private-key.asc
  ```

- **GPG_FINGERPRINT**: The full key fingerprint
  ```bash
  gpg --fingerprint KEY_ID
  ```

- **GPG_PASSPHRASE**: The key's passphrase (optional if no passphrase)

### Key Generation (For New Projects)

```bash
# Generate a new GPG key
gpg --full-gen-key
# Choose: RSA and RSA, 4096 bits
# Name: helm-chartsnap-bot
# Email: maintainer's email

# Export the key
gpg --armor --export KEY_ID > helm-chartsnap-public-key.asc
gpg --armor --export-secret-keys KEY_ID > helm-chartsnap-private-key.asc

# Get fingerprint
gpg --fingerprint KEY_ID
```

## Provenance File Format

The `.prov` files follow Helm's standard format:

```
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

name: chartsnap
version: 0.7.0
description: Snapshot testing for Helm charts
home: https://github.com/jlandowner/helm-chartsnap

...
files:
  chartsnap_v0.7.0_linux_amd64.tar.gz: sha256:abc123...

-----BEGIN PGP SIGNATURE-----
[GPG signature block]
-----END PGP SIGNATURE-----
```

This format allows Helm to:
- Verify the archive hasn't been tampered with
- Confirm it comes from a trusted source
- Ensure the archive matches the signed hash

## Backward Compatibility

- Versions < 0.7.0: No signatures, requires `--verify=false`
- Versions >= 0.7.0: Signed, works with Helm 4's default verification
- Helm 3.x: Ignores verification by default, works with all versions

## Testing Locally

To test the signing process without a release:

```bash
# Set required environment variables
export GPG_FINGERPRINT="YOUR_KEY_FINGERPRINT"
export GPG_PASSPHRASE="your_passphrase"  # optional
export GITHUB_REF_NAME="v0.7.0"

# Build with goreleaser in snapshot mode
goreleaser build --snapshot --clean

# Generate .prov files
bash scripts/generate_prov_files.sh dist "$GPG_FINGERPRINT" "$GPG_PASSPHRASE"

# Verify a .prov file
gpg --verify dist/chartsnap_v0.7.0_linux_amd64.tar.gz.prov
```

## References

- [Helm Provenance Documentation](https://helm.sh/docs/topics/provenance/)
- [GoReleaser Signing](https://goreleaser.com/customization/sign/)
- [GPG Signing in GitHub Actions](https://docs.github.com/en/authentication/managing-commit-signature-verification)
