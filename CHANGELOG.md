# Changelog

All notable changes to the RunZero Terraform provider are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.10] - 2026-10-13

### Fixed

- `ServerEnrichedMap()` plan modifier — removed merge logic that combined config
  and state maps. Terraform core requires that a planned value for
  `Optional + Computed` attributes must equal either the config value or the
  prior state value; a merged map matching neither was rejected with
  `planned value does not match config value nor prior value`. The modifier now
  only preserves prior state when config is null/unknown, and passes config
  through unchanged otherwise.
- `runzero_scan_template`, `runzero_group`, `runzero_user` resources — `Create`
  and `Update` functions now preserve the user's configured map values (`params`,
  `org_roles`) in state instead of the full API response. `Read` functions now
  filter API-returned maps to only include keys already in state, preventing
  server-added defaults from accumulating in state and causing drift.

## [0.2.9] - 2026-10-13

### Fixed

- `ServerEnrichedMap()` plan modifier — removed the `MapUnknown` behavior on
  create that caused `Provider produced invalid plan` errors. Terraform core
  rejects setting a planned value to unknown when the config provides a known
  value. The modifier now passes config through on create for all resources.
- `runzero_scan_template` resource — the `Create` function now preserves the
  user's configured `params` in state instead of the full API response. This
  prevents post-apply consistency errors when the API enriches the map with
  server-side defaults. On the next plan/refresh, `Read` fetches the full set
  and the merge modifier handles subsequent plans.

## [0.2.8] - 2026-10-13

### Added

- `ServerEnrichedMapUnknownOnCreate()` plan modifier variant for attributes
  where the server fully replaces the map on create (e.g., `scan_template.params`).

### Changed

- Split `ServerEnrichedMap()` plan modifier into two variants: the base version
  passes config through on create, while `UnknownOnCreate` marks the plan as
  unknown to allow the server response through.

## [0.2.7] - 2026-10-13

### Fixed

- Fixed `gofmt` import ordering in `group.go` and `user.go` — moved
  project-local `planmodifiers` import to the correct import group.

## [0.2.6] - 2026-10-13

### Added

- Custom `ServerEnrichedMap()` plan modifier to handle API-enriched map attributes.
  Merges user-configured keys on top of server-added defaults, preventing Terraform
  drift when the API returns additional keys beyond what the user specified.

### Fixed

- `runzero_scan_template` resource — `params` no longer causes spurious plan diffs
  ("element X has vanished" / "new element Y has appeared") when the API enriches the
  map with server-side default scan parameters.
- `runzero_group` resource — added `ServerEnrichedMap()` to `org_roles` to prevent
  potential drift if the API enriches the org-role map with additional entries.
- `runzero_user` resource — added `ServerEnrichedMap()` to `org_roles` for the same
  reason as `runzero_group`.

## [0.2.5] - 2026-10-13

### Added

- `runzero_credential` / `runzero_credentials` data sources — added `acl` (map of string)
  and `cidrs` (list of string) fields returned by the API but previously missing.
- `runzero_group` / `runzero_groups` data sources — added `org_roles` (map of string)
  field returned by the API but previously missing.

### Changed

- `runzero_scan_template` resource — changed `params` and `acl` from JSON-encoded
  `types.String` to `types.Map` of string for proper Terraform map semantics.
- `runzero_organization` resource — changed `expiration_assets_stale`,
  `expiration_assets_offline`, `expiration_integration_attributes`, `expiration_scans`,
  and `expiration_vulnerabilities` from `String` to `Int64` to match the API's numeric
  type. Changed `keep_latest_integration_attributes` from `String` to `Bool`.

### Fixed

- `Service` client struct — `Credentials` type changed from `map[string]string` to
  `map[string]any` to match Swagger `additionalProperties` with nested properties,
  consistent with the `Asset` struct. `Attributes` type changed from `map[string]any`
  to `map[string]string` to match Swagger `additionalProperties: { type: string }`,
  consistent with the `Asset` struct.
- `runzero_asset` resource — `risk_rank` changed from `Int64` to `Float64` to match
  the Swagger `type: number` definition and the client struct's `float64` type.

## [0.2.4] - 2026-10-13

### Added

- `runzero_scan_template` / `runzero_scan_templates` data sources — added `acl` field
  (map of string) returned by the API but previously missing from the datasource schema.

### Fixed

- `Service` client struct — added 25 missing fields from the Swagger `Service` schema:
  `service_screenshot_link`, `detected_by`, `os_version`, `addresses_extra`,
  `mac_vendors`, `service_count`, `service_count_tcp`, `service_count_udp`,
  `service_count_arp`, `service_count_icmp`, `lowest_ttl`, `lowest_rtt`,
  `last_agent_id`, `last_task_id`, `newest_mac`, `newest_mac_vendor`, `newest_mac_age`,
  `service_ports_tcp`, `service_ports_udp`, `service_ports_protocols`,
  `service_ports_products`, `org_name`, `site_name`, `agent_name`, `attributes`.
- `Scan` client struct — added 8 missing fields from the Swagger `TaskBase` schema:
  `template_id`, `cruncher_id`, `created_by_user_id`, `custom_integration_id`,
  `source_id`, `stats`, `parent_id`, `recur_last_task_id`.

### Audit

Full Swagger spec audit (`SWAGGER_4.0.251021.0`) of all resources and data sources
identified the following remaining field gaps:

**Data sources missing fields:**

- `runzero_organization` / `runzero_organizations` — `download_token`,
  `download_token_created_at`, `deactivated_at`, `service_count_tcp`,
  `service_count_udp`, `service_count_arp`, `service_count_icmp`
- `runzero_agent` / `runzero_agents` — `client_id`, `host_id`, `hub_id`, `system_info`,
  `deactivated_at`
- `runzero_explorer` / `runzero_explorers` — `client_id`, `host_id`, `hub_id`,
  `system_info`, `deactivated_at`
- `runzero_scan_template` / `runzero_scan_templates` — `cruncher_id`, `error`, `stats`,
  `size_site`, `size_data`, `size_results`, `linked_task_count`
- `runzero_task` / `runzero_tasks` / `runzero_account_tasks` — `client_id`,
  `cruncher_id`, `created_by_user_id`, `custom_integration_id`, `source_id`, `params`,
  `stats`, `parent_id`, `recur_last_task_id`
- `runzero_site` / `runzero_sites` / `runzero_account_sites` — `subnets`
- `runzero_service` / `runzero_services` — `service_data`, `service_screenshot_link`,
  `addresses`, `macs`, `names`, `domains`, `tags`, `rtts`, `credentials`, `services`
- `runzero_wireless_lan` / `runzero_wireless_lans` — `data`
- `runzero_asset` / `runzero_assets` — `services`, `rtts`
- `runzero_license` — `settings`
- `runzero_hosted_zones` (plural only) — `processor_id`

**Resources missing read-only response fields:**

- `runzero_organization` — `download_token`, `download_token_created_at`,
  `deactivated_at`, `service_count_tcp`, `service_count_udp`, `service_count_arp`,
  `service_count_icmp`
- `runzero_agent` — `host_id`, `hub_id`, `system_info`
- `runzero_explorer` — `host_id`, `hub_id`, `system_info`
- `runzero_scan_template` — `cruncher_id`, `created_by_user_id`, `error`, `stats`,
  `parent_id`, `custom_integration_id`, `source_id`, `template_id`, `size_site`,
  `size_data`, `size_results`, `hosted_zone_id`, `recur_last_task_id`
- `runzero_site` — `subnets`
- `runzero_user` — `client_id`, `name`, `login_failures`, `last_login_ip`,
  `last_login_at`, `last_login_ua`, `last_activity_at`, `actions`, `last_action_at`,
  `invite_token_expiration`, `reset_token_expiration`, `password_enabled_at`,
  `terms_agreed_as_of`

**Fully aligned (no gaps):** API Key, Org API Key, Export Token, SSO Group Mapping,
Custom Integration, Asset Ownership Type, Health, User (datasource), Hosted Zone
(singular).

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

[Unreleased]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.10...HEAD
[0.2.10]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.9...v0.2.10
[0.2.9]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.8...v0.2.9
[0.2.8]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.7...v0.2.8
[0.2.7]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.6...v0.2.7
[0.2.6]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.5...v0.2.6
[0.2.5]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/petrutb/terraform-provider-runzero/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/petrutb/terraform-provider-runzero/releases/tag/v0.1.0
