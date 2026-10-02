# GitHub Publish

This guide covers how to publish the runZero Terraform provider to the public [Terraform Registry](https://registry.terraform.io/providers/petrutb/runzero) from GitHub. It is fully independent of the GitLab flow described in [GITLAB_SETUP.md](GITLAB_SETUP.md).

> **Note:** This provider was built with the help of AI coding tools. See the [README](README.md) for details.

The Terraform Registry only ingests providers from **public GitHub repositories** named `terraform-provider-<name>`, so this is the path used for public distribution.

## Files

| File                                                           | Purpose                                                  |
|----------------------------------------------------------------|----------------------------------------------------------|
| [`.github/workflows/test.yml`](.github/workflows/test.yml)     | Lint, `gofmt` check, build, and unit tests on PRs and `main` pushes; reused by the release workflow |
| [`.github/workflows/release.yml`](.github/workflows/release.yml) | Runs the tests, then builds, signs, and publishes a GitHub release on `v*` tags |
| [`.goreleaser.github.yml`](.goreleaser.github.yml)             | GoReleaser config (GitHub release + registry manifest)   |
| [`terraform-registry-manifest.json`](terraform-registry-manifest.json) | Declares Terraform plugin protocol version 6       |

## Prerequisites

- A public GitHub repository: `petrutb/terraform-provider-runzero`
- A GPG key for signing release checksums (RSA or DSA; the Registry does not accept ECC keys)
- A [Terraform Registry](https://registry.terraform.io) account signed in with GitHub

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

# Private key -> GitHub secret GPG_PRIVATE_KEY
gpg --armor --export-secret-keys YOUR_FINGERPRINT

# Public key -> Terraform Registry signing key
gpg --armor --export YOUR_FINGERPRINT
```

## One-time Setup

1. **Push the code to GitHub.** If the code also lives in GitLab, add GitHub as a second remote:

   ```bash
   git remote add github https://github.com/petrutb/terraform-provider-runzero.git
   git push github main
   ```

2. **Add repository secrets** under **Settings > Secrets and variables > Actions**:

   | Secret            | Description                                      |
   |-------------------|--------------------------------------------------|
   | `GPG_PRIVATE_KEY` | ASCII-armored GPG private key                    |
   | `PASSPHRASE`      | Key passphrase (leave empty if the key has none) |

   `GITHUB_TOKEN` is provided automatically; the workflow requests `contents: write` to create releases.

3. **Register the signing key.** Add the public key at **registry.terraform.io > User Settings > Signing Keys**.

4. **Publish the provider.** On [registry.terraform.io](https://registry.terraform.io), choose **Publish > Provider** and select `petrutb/terraform-provider-runzero`. The Registry installs a webhook and ingests every new GitHub release automatically.

## Creating a Release

Tag and push to GitHub:

```bash
git tag v0.1.0
git push github v0.1.0
```

The workflow will:
1. Run the test workflow (lint, `gofmt`, build, unit tests); the release is skipped if any check fails
2. Check out the tag and set up Go from [`go.mod`](go.mod)
3. Import the GPG key
4. Build cross-platform binaries with GoReleaser (`--config .goreleaser.github.yml`)
5. Sign `SHA256SUMS` and publish a GitHub release with the zips, `SHA256SUMS`, `SHA256SUMS.sig`, and `manifest.json`

The Terraform Registry then picks up the release, usually within a few minutes.

## Consuming the Provider

```hcl
terraform {
  required_providers {
    runzero = {
      source  = "petrutb/runzero"
      version = "~> 0.1.0"
    }
  }
}

provider "runzero" {
  org_token = var.runzero_org_token
}
```

```bash
terraform init
```

## Pipeline Overview

```
PR / main push:
  test.yml → lint, gofmt, build, unit tests

Tag push (v*) to GitHub:

┌──────────────┐   ┌──────────────┐   ┌──────────────┐   ┌────────────────────┐
│  test.yml    │──▶│  import GPG  │──▶│  goreleaser  │──▶│  GitHub Release    │
│  lint+tests  │   │              │   │  build+sign  │   │                    │
└──────────────┘   └──────────────┘   └──────────────┘   └─────────┬──────────┘
                                                                   │ webhook
                                                                   ▼
                                                        ┌────────────────────┐
                                                        │ Terraform Registry │
                                                        └────────────────────┘
```

## Troubleshooting

- **Release missing in the Registry:** check **registry.terraform.io > your provider > Settings** for ingestion errors; use **Resync** after fixing.
- **Signature errors:** the key that signed `SHA256SUMS` must match a signing key registered in the Registry.
- **Missing `manifest.json`:** the release must include `terraform-provider-runzero_<version>_manifest.json`; it is attached by `release.extra_files` in `.goreleaser.github.yml`.
