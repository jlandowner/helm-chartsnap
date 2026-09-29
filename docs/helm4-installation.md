# Helm 4 Installation Guide

## Plugin Signature Support

**Starting from version 0.7.0**, helm-chartsnap releases include cryptographic signatures to support Helm 4's verification feature. This means you can install the plugin without using `--verify=false`.

### For Version 0.7.0 and Later

Install directly with signature verification enabled:

```bash
VERSION="0.7.0"  # Or any version >= 0.7.0
SYSTEM_OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
SYSTEM_ARCH="$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"

helm plugin install \
  "https://github.com/jlandowner/helm-chartsnap/releases/download/v${VERSION}/chartsnap_v${VERSION}_${SYSTEM_OS}_${SYSTEM_ARCH}.tar.gz"
```

The `.prov` provenance file is automatically used by Helm to verify the plugin's authenticity.

### For Versions Before 0.7.0

If you need to install older versions, please use one of the workarounds below.

## Understanding the Verification Challenge (Pre-0.7.0)

Helm 4 changed how plugin installation works by enabling signature checks automatically. For helm-chartsnap, this creates a challenge because:

- Our releases don't yet ship with `.prov` signature files
- Installing from our GitHub repo directly will fail verification
- The `--verify=false` workaround works but bypasses security checks

## Your Installation Options

### Option A: Use Release Archives (Our Recommendation)

Download specific release archives instead of pointing to the repo:

```bash
VERSION="0.6.0"  # Change to desired version
SYSTEM_OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
SYSTEM_ARCH="$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"

helm plugin install \
  "https://github.com/jlandowner/helm-chartsnap/releases/download/v${VERSION}/chartsnap_v${VERSION}_${SYSTEM_OS}_${SYSTEM_ARCH}.tar.gz"
```

**Why this is better:**
- Points to immutable GitHub release artifacts
- Specific version is explicit and traceable  
- GitHub's HTTPS provides transport security
- No code execution during download phase

### Option B: Clone Then Install (Best for Contributors)

If you're developing or want to inspect the code:

```bash
git clone https://github.com/jlandowner/helm-chartsnap.git /tmp/helm-chartsnap
cd /tmp/helm-chartsnap
git checkout v0.6.0  # or desired version/branch
helm plugin install "$(pwd)"
```

**Why this works:**
- Local paths skip Helm's remote verification entirely
- You can audit the source before installation
- Install hooks run from code you've reviewed
- Perfect for development workflows

### Option C: Direct Install with Verification Disabled

Quick installation that skips the verification step:

```bash
helm plugin install https://github.com/jlandowner/helm-chartsnap --verify=false
```

**Trade-offs:**
- Fastest method, one command
- Explicitly disables Helm's security feature
- Pulls from main branch by default (less stable)
- Acceptable for testing but not ideal for production

## What We've Done About This

### Version 0.7.0 Release
- ✅ Configured GoReleaser to produce `.prov` files for all release archives
- ✅ Updated release workflow to generate GPG signatures
- ✅ Added automated script to create Helm-compatible provenance files
- ✅ Updated documentation to reflect signature support

### How It Works
1. During release, GoReleaser builds the plugin archives
2. Each archive is signed with GPG to create a `.sig` file
3. A script generates Helm-compatible `.prov` provenance files containing:
   - Plugin metadata (name, version, description)
   - SHA256 hash of the archive
   - GPG signature
4. Both `.sig` and `.prov` files are uploaded to GitHub releases
5. Helm 4 automatically verifies the `.prov` file during installation

### For Repository Maintainers

To enable signing for releases, set these GitHub repository secrets:
- `GPG_PRIVATE_KEY`: The GPG private key (exported with `gpg --export-secret-keys`)
- `GPG_FINGERPRINT`: The key fingerprint (e.g., `1234567890ABCDEF`)
- `GPG_PASSPHRASE`: The key passphrase (optional if key has no passphrase)

If these secrets are not set, the release will still work but without signatures.

## Why Signatures Matter

Helm 4's verification checks that:
1. The plugin hasn't been tampered with in transit
2. It actually comes from the claimed source
3. The content matches what was officially released

Without signatures, users must trust:
- GitHub's HTTPS implementation
- That our repository hasn't been compromised
- That release artifacts are legitimate

Adding signatures will improve this trust model significantly.

## FAQ

**Q: Are signatures available now?**

A: Yes! Starting from version 0.7.0, all releases include `.prov` signature files for Helm 4 verification.

**Q: Do I still need --verify=false?**

A: Not for version 0.7.0 and later. Older versions (< 0.7.0) still require `--verify=false` or alternative installation methods.

**Q: Is --verify=false safe enough for my use case?**

A: It depends on your threat model:
- Testing/development: Probably fine
- Production/sensitive environments: Use version 0.7.0+ with signatures, or release archives with pinned versions
- High security requirements: Clone locally and audit before installing

**Q: Can I verify releases today?**

A: Yes! For version 0.7.0+, Helm automatically verifies the `.prov` file. The signatures are created using GPG and follow Helm's standard provenance format.

**Q: What about older Helm versions?**

A: Helm 3.x doesn't enforce verification by default, so existing installation methods work fine. The `--verify=false` flag is ignored in Helm 3.
