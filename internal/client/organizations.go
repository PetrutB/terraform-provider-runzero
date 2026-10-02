package client

import (
	"context"
	"fmt"
	"net/http"
)

type Organization struct {
	ID                              string `json:"id"`
	Name                            string `json:"name"`
	Description                     string `json:"description"`
	ClientID                        string `json:"client_id"`
	DownloadToken                   string `json:"download_token"`
	DownloadTokenCreatedAt          int64  `json:"download_token_created_at"`
	Inactive                        bool   `json:"inactive"`
	DeactivatedAt                   int64  `json:"deactivated_at"`
	ServiceCount                    int64  `json:"service_count"`
	ServiceCountTCP                 int64  `json:"service_count_tcp"`
	ServiceCountUDP                 int64  `json:"service_count_udp"`
	ServiceCountARP                 int64  `json:"service_count_arp"`
	ServiceCountICMP                int64  `json:"service_count_icmp"`
	AssetCount                      int64  `json:"asset_count"`
	Project                         bool   `json:"project"`
	ParentID                        string `json:"parent_id"`
	ExpirationAssetsStale           int64  `json:"expiration_assets_stale"`
	ExpirationAssetsOffline         int64  `json:"expiration_assets_offline"`
	ExpirationIntegrationAttributes int64  `json:"expiration_integration_attributes"`
	ExpirationScans                 int64  `json:"expiration_scans"`
	ExpirationVulnerabilities       int64  `json:"expiration_vulnerabilities"`
	KeepLatestIntegrationAttributes bool   `json:"keep_latest_integration_attributes"`
	ExportToken                     string `json:"export_token"`
	ExportTokenCreatedAt            int64  `json:"export_token_created_at"`
	ExportTokenLastUsedAt           int64  `json:"export_token_last_used_at"`
	ExportTokenLastUsedBy           string `json:"export_token_last_used_by"`
	ExportTokenCounter              int64  `json:"export_token_counter"`
	CreatedAt                       int64  `json:"created_at"`
	UpdatedAt                       int64  `json:"updated_at"`
}

type CreateOrganizationRequest struct {
	Name                            string `json:"name"`
	Description                     string `json:"description,omitempty"`
	ParentID                        string `json:"parent_id,omitempty"`
	ExpirationAssetsStale           string `json:"expiration_assets_stale,omitempty"`
	ExpirationAssetsOffline         string `json:"expiration_assets_offline,omitempty"`
	ExpirationIntegrationAttributes string `json:"expiration_integration_attributes,omitempty"`
	ExpirationScans                 string `json:"expiration_scans,omitempty"`
	ExpirationVulnerabilities       string `json:"expiration_vulnerabilities,omitempty"`
	KeepLatestIntegrationAttributes string `json:"keep_latest_integration_attributes,omitempty"`
}

type UpdateOrganizationRequest struct {
	Name                            string `json:"name,omitempty"`
	Description                     string `json:"description,omitempty"`
	ExpirationAssetsStale           string `json:"expiration_assets_stale,omitempty"`
	ExpirationAssetsOffline         string `json:"expiration_assets_offline,omitempty"`
	ExpirationIntegrationAttributes string `json:"expiration_integration_attributes,omitempty"`
	ExpirationScans                 string `json:"expiration_scans,omitempty"`
	ExpirationVulnerabilities       string `json:"expiration_vulnerabilities,omitempty"`
	KeepLatestIntegrationAttributes string `json:"keep_latest_integration_attributes,omitempty"`
}

func (c *Client) CreateOrganization(ctx context.Context, req CreateOrganizationRequest) (*Organization, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/orgs", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := HandleResponse(resp, &org); err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}
	return &org, nil
}

func (c *Client) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/orgs/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := HandleResponse(resp, &org); err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}
	return &org, nil
}

func (c *Client) UpdateOrganization(ctx context.Context, id string, req UpdateOrganizationRequest) (*Organization, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/orgs/%s", id), AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := HandleResponse(resp, &org); err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}
	return &org, nil
}

func (c *Client) DeleteOrganization(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/orgs/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListOrganizations(ctx context.Context) ([]Organization, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/orgs", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var orgs []Organization
	if err := HandleResponse(resp, &orgs); err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	return orgs, nil
}

// GetOrg returns the organization associated with the org-scoped token (GET /org).
func (c *Client) GetOrg(ctx context.Context) (*Organization, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := HandleResponse(resp, &org); err != nil {
		return nil, fmt.Errorf("get org: %w", err)
	}
	return &org, nil
}

// UpdateOrg updates the organization associated with the org-scoped token (PATCH /org).
func (c *Client) UpdateOrg(ctx context.Context, req UpdateOrganizationRequest) (*Organization, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, "/org", AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := HandleResponse(resp, &org); err != nil {
		return nil, fmt.Errorf("update org: %w", err)
	}
	return &org, nil
}
