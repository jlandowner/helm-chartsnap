# Implementing Plugin Signing for Helm 4

This document outlines the plan to add cryptographic signatures to helm-chartsnap releases.

## Current Situation

- Helm 4 enforces plugin verification by default
- Our releases lack `.prov` signature files
- Users must use `--verify=false` or alternative installation methods
- This bypasses an important security feature

## Implementation Plan

### Step 1: GPG Key Setup

Create a dedicated GPG key for signing releases:

```bash
gpg --full-gen-key
# Choose: RSA and RSA, 4096 bits
# Name: helm-chartsnap-bot (or project name)
# Email: maintainer email
```

Export and store the key securely:
```bash
gpg --armor --export <KEY_ID> > helm-chartsnap-public-key.asc
gpg --export-secret-keys <KEY_ID> > helm-chartsnap-private-key.gpg
```

Store private key in GitHub Secrets as `GPG_PRIVATE_KEY` and passphrase as `GPG_PASSPHRASE`.

### Step 2: Update GoReleaser Configuration

Modify `.goreleaser.yml` to enable signing:

```yaml
# Add to the file
signs:
  - cmd: gpg
    args:
      - --batch
      - --local-user
      - "{{ .Env.GPG_FINGERPRINT }}"
      - --output
      - "${signature}"
      - --detach-sign
      - "${artifact}"
    artifacts: all
    signature: "${artifact}.sig"
```

For Helm plugin-specific provenance, we need to generate `.prov` files that follow Helm's format.

### Step 3: Modify Release Workflow

Update `.github/workflows/release.yaml`:

```yaml
- name: Import GPG key
  run: |
    echo "${{ secrets.GPG_PRIVATE_KEY }}" | gpg --batch --import
    
- name: Configure Git for signing
  run: |
    git config --global user.signingkey ${{ secrets.GPG_FINGERPRINT }}
    git config --global commit.gpgsign true
```

### Step 4: Generate Helm-Compatible Provenance

After GoReleaser creates the archives, generate Helm provenance files:

```bash
for file in dist/*.tar.gz; do
  gpg --armor --detach-sign --local-user <KEY_ID> "$file"
  mv "$file.asc" "$file.prov"
done
```

Upload `.prov` files alongside release artifacts.

### Step 5: Publish Public Key

Add the public key to the repository:
- Commit `helm-chartsnap-public-key.asc` to repo root
- Document the key fingerprint in README
- Add verification instructions

### Step 6: Update Documentation

Modify installation docs to show verification:

```bash
# Import public key
curl -fsSL https://raw.githubusercontent.com/jlandowner/helm-chartsnap/main/helm-chartsnap-public-key.asc | gpg --import

# Install with verification (no --verify=false needed!)
helm plugin install https://github.com/jlandowner/helm-chartsnap/releases/download/v0.7.0/chartsnap_v0.7.0_linux_amd64.tar.gz
```

## Testing Plan

1. Create test GPG key in development environment
2. Run release process locally with `goreleaser build --snapshot`
3. Verify signature files are created correctly
4. Test installation with and without verification
5. Confirm error messages when signature is invalid

## Rollout Strategy

1. Implement signing in a beta release first
2. Document the new process thoroughly  
3. Keep `--verify=false` instructions for older releases
4. Announce signing support in release notes
5. Update README to prefer signed installation after 1-2 releases
6. Eventually deprecate `--verify=false` examples

## Alternative: Sigstore/Cosign

As an alternative to GPG, consider Sigstore's cosign for keyless signing:

```bash
# Sign with cosign (uses OIDC, no key management)
cosign sign-blob --bundle chartsnap.bundle chartsnap_v0.7.0_linux_amd64.tar.gz

# Verify
cosign verify-blob --bundle chartsnap.bundle --certificate-identity=... --certificate-oidc-issuer=...
```

However, Helm's native verification only supports GPG/PGP signatures, so this would require custom verification scripts.

## Timeline

- **Immediate**: Document current state (this PR)
- **Next release (0.7.0)**: Implement GPG signing
- **Following release (0.8.0)**: Make signed installation the default
- **Future**: Consider Sigstore integration if Helm adds support

## References

- [Helm Provenance Documentation](https://helm.sh/docs/topics/provenance/)
- [GoReleaser Signing Documentation](https://goreleaser.com/customization/sign/)
- [GitHub Actions GPG Signing](https://docs.github.com/en/authentication/managing-commit-signature-verification)
