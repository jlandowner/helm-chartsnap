# Helm 4 Installation Guide

## Understanding the Verification Challenge

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

## What We're Doing About This

We recognize that `--verify=false` isn't ideal. Here's our roadmap:

### Near-term (Next Release)
- Generate GPG key for signing releases
- Configure goreleaser to produce `.prov` files
- Publish our public key in the repository
- Document the verification process

### Long-term
- Automatic signature generation in CI
- Verification instructions in main README
- Remove `--verify=false` from all examples
- Consider Sigstore/cosign integration

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

**Q: Is --verify=false safe enough for my use case?**

A: It depends on your threat model:
- Testing/development: Probably fine
- Production/sensitive environments: Use release archives with pinned versions
- High security requirements: Clone locally and audit before installing

**Q: Can I verify releases today without signatures?**

A: Not cryptographically, but you can:
- Check release checksums if provided
- Compare against known-good versions
- Build from source and compare binaries
- Monitor GitHub for unexpected changes

**Q: When will signatures be available?**

A: We're targeting the next minor release. Check the [repository issues](https://github.com/jlandowner/helm-chartsnap/issues) for updates on implementation progress.

**Q: What about older Helm versions?**

A: Helm 3.x doesn't enforce verification by default, so existing installation methods work fine. The `--verify=false` flag is ignored in Helm 3.
