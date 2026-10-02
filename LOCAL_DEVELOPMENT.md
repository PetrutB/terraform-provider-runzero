# Local Development

Quick guide to building and using the runZero Terraform provider from source.

Released versions are distributed through the [Terraform Registry](https://registry.terraform.io/providers/petrutb/runzero) and installed with a normal `terraform init`. Follow this guide only to test unreleased changes or contribute to the provider.

> **Note:** This provider was built with the help of AI coding tools. Test changes against a non-production runZero organization.

## Prerequisites

- [Go](https://go.dev/dl/) — version per [`go.mod`](go.mod) (currently 1.25.8+)
- [Terraform 1.0+](https://developer.hashicorp.com/terraform/install)

## Build

```bash
git clone https://github.com/petrutb/terraform-provider-runzero.git
cd terraform-provider-runzero
go mod tidy
go build -o terraform-provider-runzero        # Linux / macOS
go build -o terraform-provider-runzero.exe    # Windows
```

Or just (Linux / macOS):

```bash
make build
```

## Use Locally

Add dev overrides to your Terraform CLI config so it picks up the local binary instead of downloading from the registry.

**Linux / macOS** — `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "petrutb/runzero" = "/full/path/to/terraform-provider-runzero"
  }
  direct {}
}
```

**Windows** — `%APPDATA%\terraform.rc`:

```hcl
provider_installation {
  dev_overrides {
    "petrutb/runzero" = "C:\\Users\\you\\terraform-provider-runzero"
  }
  direct {}
}
```

Set the path to the **directory** containing the built binary, not the binary itself.

With dev overrides active you **skip `terraform init`** — Terraform loads the provider directly. You'll see a warning about the override, which is expected.

To switch back to the released provider from the Terraform Registry, remove the `dev_overrides` block and run `terraform init`.

### Alternative: local plugin directory

`make install` builds the binary and copies it to `~/.terraform.d/plugins/registry.terraform.io/petrutb/runzero/0.1.0/<OS>_<ARCH>/`. Terraform then resolves the provider from that directory during a normal `terraform init`, so pin `version = "0.1.0"` in `required_providers`.

## Write a Config

```hcl
terraform {
  required_providers {
    runzero = {
      source = "petrutb/runzero"
    }
  }
}

provider "runzero" {
  org_token = "OT_your_token_here"
  # or set RUNZERO_ORG_TOKEN env var
}

resource "runzero_site" "dev" {
  name  = "Dev Network"
  scope = "192.168.1.0/24"
}
```

Then:

```bash
terraform plan
terraform apply
```

## Environment Variables

| Variable                | Purpose                          |
|-------------------------|----------------------------------|
| `RUNZERO_API_URL`       | API base URL (default: `https://console.runzero.com/api/v1.0`) |
| `RUNZERO_ORG_TOKEN`     | Organization token (OT prefix)   |
| `RUNZERO_ACCOUNT_TOKEN` | Account token (CT prefix)        |
| `RUNZERO_CLIENT_ID`     | OAuth2 client ID                 |
| `RUNZERO_CLIENT_SECRET` | OAuth2 client secret             |

## Run Tests

Unit tests (no credentials needed):

```bash
go test ./...
# or
make test
```

Acceptance tests (real API calls — use a test org):

```bash
export RUNZERO_ORG_TOKEN="OT_your_token"
export RUNZERO_ACCOUNT_TOKEN="CT_your_token"
export TF_ACC=1
go test ./internal/provider/acceptance/ -v -timeout 60m
# or
make testacc
```

Tests that need objects Terraform cannot create (agents, explorers, assets, services, hosted zones, wireless LANs) are skipped unless their IDs are provided, e.g. `RUNZERO_TEST_EXPLORER_ID`. See [`internal/provider/acceptance/helpers_test.go`](internal/provider/acceptance/helpers_test.go) for the full list.

## Release Build

Cross-platform binaries via goreleaser (local snapshot, nothing is published):

```bash
goreleaser release --snapshot --clean --skip=sign --config .goreleaser.github.yml
```

Official releases are published to the Terraform Registry via GitHub ([GITHUB_SETUP.md](GITHUB_SETUP.md)), or internally via GitLab ([GITLAB_SETUP.md](GITLAB_SETUP.md)).

## Makefile Targets

| Target     | What it does                    |
|------------|----------------------------------|
| `build`    | Build the provider binary        |
| `install`  | Build + install to `~/.terraform.d/plugins` (local plugin directory) |
| `test`     | Unit tests                       |
| `testacc`  | Acceptance tests (`TF_ACC=1`)    |
| `lint`     | Run linters                      |
| `generate` | `go generate` for docs/codegen   |
