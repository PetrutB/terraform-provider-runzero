# Changelog

All notable changes to the RunZero Terraform provider are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.3] - 2026-10-01

### Fixed

- Updating docs through tfplugindocs

## [0.2.2] - 2026-10-01

### Fixed

- `ScanTemplate.grace_period`, `ScanTemplate.source_id` — changed Go type from `string`
  to `int64` to match the API's numeric JSON response, fixing `json: cannot unmarshal
  number into Go struct field` errors on `runzero_scan_template` and
  `runzero_scan_templates` data sources.

## [0.2.1] - 2026-10-01

### Fixed

- Fixed gofmt failing on GitHub Actions

## [0.2.0] - 2026-10-01

### Fixed

- `Credential.ACL`, `Asset.credentials`, `Agent.system_info`, `Explorer.system_info`,
  `Group.org_roles`, `User.org_roles` — changed Go types from `map[string]string` to
  `map[string]any` to match the Swagger `additionalProperties: true` definitions and
  prevent unmarshal errors when the API returns non-string values.
- `Task.source_id`, `WirelessLAN.signal` — changed Go type from `int` to `int64` for
  consistency with all other integer fields and correctness on 32-bit platforms.

## [0.1.0] - 2026-10-01

Initial release of the runZero Terraform provider, built on the
[terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
(plugin protocol 6) and validated against the runZero API OpenAPI spec `4.0.251021.0`.

This provider was developed with the assistance of AI coding tools and reviewed against
the runZero API specification.

### Provider

- Authentication with an organization token (`OT`), an account token (`CT`), or OAuth2
  client credentials; at least one must be configured.
- Every provider attribute can be set via environment variables: `RUNZERO_ORG_TOKEN`,
  `RUNZERO_ACCOUNT_TOKEN`, `RUNZERO_CLIENT_ID`, `RUNZERO_CLIENT_SECRET`, `RUNZERO_API_URL`.
- Configurable API base URL (default `https://console.runzero.com/api/v1.0`).

### Resources

Account and identity:

- `runzero_organization` — manages organizations.
- `runzero_user` — manages users.
- `runzero_user_invite` — sends an email invitation to a new user (create-only).
- `runzero_group` — manages groups.
- `runzero_sso_group_mapping` — maps SSO groups to runZero groups.
- `runzero_api_key` — manages API keys.
- `runzero_org_api_key` — manages the API key of the organization tied to the configured
  organization token. The key is adopted on create, rotated when the `rotate` value
  changes, and removed on destroy.
- `runzero_export_token` — manages export tokens.

Sites, infrastructure, and scanning:

- `runzero_site` — manages sites.
- `runzero_agent` — adopts an externally deployed agent by ID.
- `runzero_explorer` — adopts an externally deployed explorer by ID.
- `runzero_scan_template` — manages reusable scan templates.
- `runzero_scan` — triggers a scan on a site. Supports the full set of scan options,
  including scheduling (`scan_start`, `scan_frequency`), explorer and hosted-zone
  targeting, rate and probe tuning, subnet/host ping, screenshots, and custom
  nameservers. Any change forces a new scan.

Assets and integrations:

- `runzero_asset` — adopts a discovered asset by ID and manages its comments and
  criticality.
- `runzero_asset_ownership_type` — manages asset ownership types.
- `runzero_credential` — manages credentials, with a write-only, sensitive `secret` map
  whose keys depend on the credential `type` (for example `aws_access_secret` →
  `access_key`/`secret_access_key`, `snmpv2_community` → `community`).
- `runzero_custom_integration` — manages custom integrations.

All resources except `runzero_scan` and `runzero_org_api_key` support
`terraform import`.

### Data Sources

Single-item and list lookups:

| Object                 | Single                           | List                              |
|------------------------|----------------------------------|-----------------------------------|
| Agents                 | `runzero_agent`                  | `runzero_agents`                  |
| API keys               | `runzero_api_key`                | `runzero_api_keys`                |
| Assets                 | `runzero_asset`                  | `runzero_assets`                  |
| Asset ownership types  | `runzero_asset_ownership_type`   | `runzero_asset_ownership_types`   |
| Credentials            | `runzero_credential`             | `runzero_credentials`             |
| Custom integrations    | `runzero_custom_integration`     | `runzero_custom_integrations`     |
| Explorers              | `runzero_explorer`               | `runzero_explorers`               |
| Export tokens          | `runzero_export_token`           | `runzero_export_tokens`           |
| Groups                 | `runzero_group`                  | `runzero_groups`                  |
| Hosted zones           | `runzero_hosted_zone`            | `runzero_hosted_zones`            |
| Organizations          | `runzero_organization`           | `runzero_organizations`           |
| Scan templates         | `runzero_scan_template`          | `runzero_scan_templates`          |
| Services               | `runzero_service`                | `runzero_services`                |
| Sites                  | `runzero_site`                   | `runzero_sites`                   |
| SSO group mappings     | `runzero_sso_group_mapping`      | `runzero_sso_group_mappings`      |
| Tasks                  | `runzero_task`                   | `runzero_tasks`                   |
| Users                  | `runzero_user`                   | `runzero_users`                   |
| Wireless LANs          | `runzero_wireless_lan`           | `runzero_wireless_lans`           |

Additional data sources:

- `runzero_org` — the organization tied to the configured organization token.
- `runzero_org_api_key` — the API key of that organization.
- `runzero_account_sites` and `runzero_account_tasks` — sites and tasks across all
  organizations in the account.
- `runzero_license` — account license details.
- `runzero_health` — API health status.

### Documentation

- Registry documentation for every resource and data source under `docs/`.
- Runnable example configurations under `examples/`.
- Guides for local development (`LOCAL_DEVELOPMENT.md`) and for publishing via GitHub
  (`GITHUB_SETUP.md`) and GitLab (`GITLAB_SETUP.md`).

### Distribution

- Published to the [Terraform Registry](https://registry.terraform.io/providers/petrutb/runzero)
  as `petrutb/runzero`.
- Signed release binaries for Linux, macOS, and Windows (`amd64`, `arm`, `arm64`) and
  FreeBSD (`amd64`, `arm`).

### Known Limitations

- Agents, explorers, and assets cannot be created by Terraform; they are deployed or
  discovered outside Terraform and adopted by ID.
- Bulk exports, cross-cloud integrations (ServiceNow, Splunk, Cisco), asset import, and
  imperative actions (bulk asset operations, user resets, task logs and data) are out of
  scope for this provider.

[Unreleased]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/petrutb/terraform-provider-runzero/releases/tag/v0.1.0
