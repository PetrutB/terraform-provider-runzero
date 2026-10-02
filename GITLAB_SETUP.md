# GitLab Publish

This guide covers how to build, publish, and consume the runZero Terraform provider from a GitLab instance. It is fully independent of the GitHub / Terraform Registry flow described in [GITHUB_SETUP.md](GITHUB_SETUP.md).

> **Note:** This provider was built with the help of AI coding tools. See the [README](README.md) for details.

Use this path for internal, private, or air-gapped distribution. The public Terraform Registry only ingests providers from GitHub, so GitLab-published releases are installed through a local filesystem mirror rather than `terraform init` against the Registry.

## Files

| File                                             | Purpose                                                   |
|--------------------------------------------------|-----------------------------------------------------------|
| [`.gitlab-ci.yml`](.gitlab-ci.yml)               | Lint, test, build, GitLab release, package upload         |
| [`.goreleaser.gitlab.yml`](.goreleaser.gitlab.yml) | GoReleaser config (builds and signs; release disabled)  |

## Prerequisites

- GitLab 15.0+ (uses the Generic Package Registry)
- A GPG key for signing release checksums
- Go (version per [`go.mod`](go.mod)) and GoReleaser (handled by CI)

## CI/CD Variables

Set these in **Settings > CI/CD > Variables** (mask sensitive values):

| Variable           | Description                          | Protected | Masked |
|--------------------|--------------------------------------|-----------|--------|
| `GPG_PRIVATE_KEY`  | ASCII-armored GPG private key        | Yes       | Yes    |
| `GPG_FINGERPRINT`  | GPG key fingerprint for signing      | Yes       | No     |

Protected variables are only exposed to protected refs, so protect the `v*` tag pattern under **Settings > Repository > Protected tags**.

### Generate a GPG key (if you don't have one)

```bash
gpg --batch --gen-key <<EOF
Key-Type: RSA
Key-Length: 4096
Name-Real: Terraform Provider RunZero
Name-Email: your-email@example.com
Expire-Date: 0
%no-protection
EOF

# Get the fingerprint
gpg --list-keys --keyid-format long

# Export the private key (paste this into GPG_PRIVATE_KEY variable)
gpg --armor --export-secret-keys YOUR_FINGERPRINT
```

## Creating a Release

Tag and push to GitLab:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The CI pipeline will:
1. Run unit tests
2. Build cross-platform binaries with GoReleaser (`--config .goreleaser.gitlab.yml`)
3. Sign the `SHA256SUMS` with your GPG key
4. Upload the zips, `SHA256SUMS`, `SHA256SUMS.sig`, and manifest to the GitLab Generic Package Registry
5. Create a GitLab Release

## Consuming the Provider

The provider `source` address is `registry.terraform.io/petrutb/runzero`. GitLab-hosted artifacts are resolved **locally** through a filesystem mirror or `dev_overrides`, so Terraform never contacts the public Registry for this provider.

### Option 1: Install from GitLab Package Registry

Create a script to download and install the provider binary:

```bash
#!/bin/bash
set -euo pipefail

GITLAB_URL="https://gitlab.example.com"
PROJECT_ID="12345"                    # Your GitLab project ID
VERSION="0.1.0"
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
PROVIDER_NAME="runzero"

ARCHIVE="terraform-provider-${PROVIDER_NAME}_${VERSION}_${OS}_${ARCH}.zip"
INSTALL_DIR="${HOME}/.terraform.d/plugins/registry.terraform.io/petrutb/${PROVIDER_NAME}/${VERSION}/${OS}_${ARCH}"

mkdir -p "${INSTALL_DIR}"

curl --fail --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" \
  -L "${GITLAB_URL}/api/v4/projects/${PROJECT_ID}/packages/generic/${PROVIDER_NAME}/v${VERSION}/${ARCHIVE}" \
  -o "/tmp/${ARCHIVE}"

unzip -o "/tmp/${ARCHIVE}" -d "${INSTALL_DIR}"
rm "/tmp/${ARCHIVE}"

echo "Installed to ${INSTALL_DIR}"
```

Then in your Terraform config:

```hcl
terraform {
  required_providers {
    runzero = {
      source  = "registry.terraform.io/petrutb/runzero"
      version = "0.1.0"
    }
  }
}

provider "runzero" {
  org_token = var.runzero_org_token
}
```

### Option 2: Filesystem mirror

Download the provider zip and set up a local filesystem mirror:

```bash
# Directory structure Terraform expects
mkdir -p /opt/terraform/providers/registry.terraform.io/petrutb/runzero/0.1.0/linux_amd64

# Download and extract
curl --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" \
  -L "${GITLAB_URL}/api/v4/projects/${PROJECT_ID}/packages/generic/runzero/v0.1.0/terraform-provider-runzero_0.1.0_linux_amd64.zip" \
  -o /opt/terraform/providers/registry.terraform.io/petrutb/runzero/0.1.0/linux_amd64/terraform-provider-runzero_0.1.0_linux_amd64.zip
```

Configure Terraform to use the filesystem mirror in `~/.terraformrc`:

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/opt/terraform/providers"
    include = ["registry.terraform.io/petrutb/*"]
  }
  direct {
    exclude = ["registry.terraform.io/petrutb/*"]
  }
}
```

### Option 3: dev_overrides (for development)

Build locally and point Terraform at the binary directory (see [LOCAL_DEVELOPMENT.md](LOCAL_DEVELOPMENT.md) for full steps):

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/petrutb/runzero" = "/path/to/terraform-provider-runzero"
  }
  direct {}
}
```

This skips `terraform init` and loads the provider directly.

## Pipeline Overview

```
Tag push (v*) triggers:

┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│  lint    │──▶│  test    │──▶│  build   │──▶│ release  │──▶│ publish  │
│ fmt-check│   │unit-test │   │goreleaser│   │ gitlab   │   │ packages │
└──────────┘   │build-chk │   └──────────┘   └──────────┘   └──────────┘
               └──────────┘

MR / main push triggers:
  lint, fmt-check, unit-test, build-check only
```

## Verifying a Release

After the pipeline completes, verify artifacts exist:

```bash
# List packages
curl --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" \
  "${GITLAB_URL}/api/v4/projects/${PROJECT_ID}/packages?package_name=runzero"

# Download and verify checksums
curl -L -o SHA256SUMS "...url.../terraform-provider-runzero_0.1.0_SHA256SUMS"
curl -L -o SHA256SUMS.sig "...url.../terraform-provider-runzero_0.1.0_SHA256SUMS.sig"
gpg --verify SHA256SUMS.sig SHA256SUMS
sha256sum -c SHA256SUMS
```
