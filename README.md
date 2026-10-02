# Terraform Provider for runZero

A [Terraform](https://www.terraform.io) provider for the [runZero](https://www.runzero.com/) cyber asset attack surface management (CAASM) platform. It provides infrastructure-as-code management of organizations, sites, users, groups, agents, explorers, assets, scans, credentials, and related configuration through the runZero API.

- **Module:** `github.com/petrutb/terraform-provider-runzero`
- **Provider source:** `registry.terraform.io/petrutb/runzero`
- **Built with:** [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) (Terraform wire protocol 6)
- **Validated against:** runZero API **v4.0.251021.0** (OpenAPI 3.0)
- **Distribution:** [Terraform Registry](https://registry.terraform.io/providers/petrutb/runzero)

> [!NOTE]
> **Built with AI.** This provider was developed with the assistance of AI coding tools. The API client, resources, data sources, tests, and documentation were generated and iterated on with AI, then reviewed and validated against the runZero OpenAPI specification. Review the code and test against a non-production runZero organization before relying on it in production. Bug reports and Bug fixes are welcomed.

---

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.0
- [Go](https://go.dev/dl/) — version per [`go.mod`](go.mod) (build from source only)
- A runZero account with an API token (Organization `OT`, Account `CT`) or OAuth2 client credentials

## Authentication

Configure at least one credential on the provider block (or via environment variables):

```hcl
provider "runzero" {
  org_token = "OT_your_organization_token" # /org/* endpoints
  # account_token = "CT_your_account_token" # /account/* endpoints
  # client_id     = "..."                   # OAuth2 client credentials
  # client_secret = "..."
  # api_url       = "https://console.runzero.com/api/v1.0" # optional override
}
```

| Attribute       | Environment variable    | Purpose                                             |
|-----------------|-------------------------|-----------------------------------------------------|
| `org_token`     | `RUNZERO_ORG_TOKEN`     | Organization token (`OT` prefix) — `/org/*`         |
| `account_token` | `RUNZERO_ACCOUNT_TOKEN` | Account token (`CT` prefix) — `/account/*`          |
| `client_id`     | `RUNZERO_CLIENT_ID`     | OAuth2 client ID                                    |
| `client_secret` | `RUNZERO_CLIENT_SECRET` | OAuth2 client secret                                |
| `api_url`       | `RUNZERO_API_URL`       | API base URL (default `https://console.runzero.com/api/v1.0`) |

Account-level resources (organizations, users, groups, API keys, credentials, …) require an account token or OAuth2; organization-level resources (sites, agents, assets, scans, …) accept an organization token.

## Example Usage

```hcl
terraform {
  required_providers {
    runzero = {
      source = "registry.terraform.io/petrutb/runzero"
    }
  }
}

provider "runzero" {
  org_token = var.runzero_org_token
}

resource "runzero_site" "production" {
  name  = "Production Network"
  scope = "10.0.0.0/8"
}

resource "runzero_scan" "weekly" {
  site_id        = runzero_site.production.id
  name           = "Weekly Network Scan"
  targets        = "10.0.0.0/24"
  scan_frequency = "weekly"
}
```

## Resources (17)

`runzero_organization`, `runzero_site`, `runzero_user`, `runzero_user_invite`, `runzero_group`, `runzero_api_key`, `runzero_credential`, `runzero_custom_integration`, `runzero_scan_template`, `runzero_export_token`, `runzero_asset_ownership_type`, `runzero_sso_group_mapping`, `runzero_agent`, `runzero_explorer`, `runzero_asset`, `runzero_scan`, `runzero_org_api_key`

## Data Sources (42)

Single + list lookups for: assets, agents, services, explorers, tasks, hosted zones, wireless LANs, organizations (incl. `runzero_org` self), sites, users, groups, API keys, export tokens, scan templates, SSO group mappings, asset ownership types, credentials, custom integrations — plus `runzero_org_api_key`, `runzero_license`, `runzero_health`, `runzero_account_sites`, and `runzero_account_tasks`.

Full attribute references for every resource and data source live in [`docs/`](docs/), with runnable configurations in [`examples/`](examples/).

## Installation

### From the Terraform Registry

The provider is distributed through the public [Terraform Registry](https://registry.terraform.io/providers/petrutb/runzero). Declare it in `required_providers` and run `terraform init` — Terraform downloads and verifies the signed binary automatically:

```hcl
terraform {
  required_providers {
    runzero = {
      source  = "petrutb/runzero"
      version = "~> 0.1.0"
    }
  }
}
```

```bash
terraform init
```

### Using the provider locally (from source)

To run an unreleased build or test local changes, build the binary and point Terraform at it instead of the registry.

**1. Build the binary**

```bash
git clone https://github.com/petrutb/terraform-provider-runzero.git
cd terraform-provider-runzero
go mod tidy
go build -o terraform-provider-runzero        # Linux / macOS
go build -o terraform-provider-runzero.exe    # Windows
```

**2. Configure a `dev_overrides` block** in your Terraform CLI config file so Terraform loads the local binary:

- Linux / macOS: `~/.terraformrc`
- Windows: `%APPDATA%\terraform.rc`

```hcl
provider_installation {
  dev_overrides {
    # Path to the DIRECTORY containing the built binary, not the binary itself
    "petrutb/runzero" = "/full/path/to/terraform-provider-runzero"
    # Windows example: "C:\\Users\\you\\terraform-provider-runzero"
  }
  direct {}
}
```

**3. Write a configuration** using the same source address as the registry:

```hcl
terraform {
  required_providers {
    runzero = {
      source = "petrutb/runzero"
    }
  }
}

provider "runzero" {
  # or set RUNZERO_ORG_TOKEN in the environment
  org_token = "OT_your_token_here"
}

resource "runzero_site" "dev" {
  name  = "Dev Network"
  scope = "192.168.1.0/24"
}
```

**4. Plan and apply** — with `dev_overrides` active, **skip `terraform init`**; Terraform loads the provider directly and prints a warning about the override, which is expected:

```bash
terraform plan
terraform apply
```

Remove the `dev_overrides` block to switch back to the Terraform Registry release.

**Alternative: local plugin directory.** Instead of `dev_overrides`, `make install` copies the binary into `~/.terraform.d/plugins/registry.terraform.io/petrutb/runzero/0.1.0/<OS>_<ARCH>/`, after which a normal `terraform init` resolves the provider from that directory.

See [LOCAL_DEVELOPMENT.md](LOCAL_DEVELOPMENT.md) for additional details.

## Development

```bash
git clone https://github.com/petrutb/terraform-provider-runzero.git
cd terraform-provider-runzero
go mod tidy
go build ./...          # or: make build
go test ./...           # unit + mock tests; or: make test
```

### Testing

| Suite | Location | Needs | What it covers |
|-------|----------|-------|----------------|
| Client unit tests | `internal/client/` | nothing | Every API client method: HTTP method, path, query, auth token type, request/response encoding, OAuth2 token caching, rate-limit retry |
| Schema tests | `internal/provider/` | nothing | All 17 resources and 42 data sources: valid schemas, naming, sensitive credentials, import support |
| Mock lifecycle tests | `internal/provider/mocktest/` | Terraform CLI | Real `terraform plan/apply/import/destroy` against an in-memory fake API for sites, groups, scans, and the org API key |
| Acceptance tests | `internal/provider/acceptance/` | `TF_ACC=1` + real credentials | Every resource and data source against the real runZero API |

`go test ./...` runs the first three suites; acceptance tests are skipped unless `TF_ACC=1` is set.

#### Acceptance tests

Acceptance tests create, update, import, and delete real objects, so run them against a dedicated test account and organization:

```bash
export RUNZERO_API_URL="https://console.runzero.com/api/v1.0"   # optional
export RUNZERO_ORG_TOKEN="OT_..."
export RUNZERO_ACCOUNT_TOKEN="CT_..."
export TF_ACC=1
go test ./internal/provider/acceptance/ -v -timeout 60m   # or: make testacc
```

Each resource test creates the object, reads it back through its single and list data sources, verifies `terraform import`, applies an update, and confirms via the API that the object is gone after destroy. Rotation of API keys and export tokens is also verified.

Some objects cannot be created through the API. Tests for them run only when an existing ID is provided and are skipped otherwise:

| Variable | Used by |
|----------|---------|
| `RUNZERO_TEST_EXPLORER_ID` | `runzero_explorer`, `runzero_scan`, task data sources |
| `RUNZERO_TEST_AGENT_ID` | `runzero_agent` |
| `RUNZERO_TEST_ASSET_ID` | `runzero_asset` (its `comments` are modified) |
| `RUNZERO_TEST_SERVICE_ID` | `runzero_service` data source |
| `RUNZERO_TEST_TASK_ID` | `runzero_task` data source |
| `RUNZERO_TEST_HOSTED_ZONE_ID` | `runzero_hosted_zone` data source |
| `RUNZERO_TEST_WIRELESS_LAN_ID` | `runzero_wireless_lan` data source |
| `RUNZERO_TEST_INVITE_EMAIL` | `runzero_user_invite` (sends a real invitation email) |
| `RUNZERO_TEST_EMAIL_DOMAIN` | Domain for generated user emails (default `example.com`) |

Safeguards:

- Agents, explorers, assets, and the organization API key are adopted and then released from state with a `removed` block (Terraform 1.7+), so they are never deleted.
- The organization API key is never rotated, because that would invalidate the token used by the tests.
- Scans only target `192.0.2.0/30` (TEST-NET-1, reserved for documentation).

> [!NOTE]
> On Windows, Terraform's test harness symlinks subdirectories of the test package, which requires Developer Mode. The mock and acceptance tests therefore live in directories without subdirectories.

### Project layout

```
main.go                     Provider entry point (plugin server)
internal/client/            Hand-written runZero API client (structs, endpoints)
internal/provider/          Provider definition
  resources/                Managed resources (create/read/update/delete)
  datasources/              Read-only data sources
  mocktest/                 Lifecycle tests against an in-memory fake API
  acceptance/               Acceptance tests against the real runZero API
docs/                       Terraform Registry documentation
examples/                   Example Terraform configurations
```

## Releasing

Releases are built with [GoReleaser](https://goreleaser.com/) on tags matching `vX.Y.Z` and produce cross-platform binaries, a GPG-signed `SHA256SUMS`, and the [`terraform-registry-manifest.json`](terraform-registry-manifest.json). There are two independent publishing paths:

| Path   | Files | Result | Guide |
|--------|-------|--------|-------|
| GitHub | [`.github/workflows/test.yml`](.github/workflows/test.yml), [`.github/workflows/release.yml`](.github/workflows/release.yml), [`.goreleaser.github.yml`](.goreleaser.github.yml) | GitHub release ingested by the public Terraform Registry | [GITHUB_SETUP.md](GITHUB_SETUP.md) |
| GitLab | [`.gitlab-ci.yml`](.gitlab-ci.yml), [`.goreleaser.gitlab.yml`](.goreleaser.gitlab.yml) | GitLab release + Generic Package Registry (internal / air-gapped use) | [GITLAB_SETUP.md](GITLAB_SETUP.md) |

See [CHANGELOG.md](CHANGELOG.md) for version history.

## License

See the repository license terms. runZero is a trademark of runZero, Inc.; this project is an independent Terraform provider.
