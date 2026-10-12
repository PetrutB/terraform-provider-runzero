package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Asset struct {
	ID                 string            `json:"id"`
	OrganizationID     string            `json:"organization_id"`
	SiteID             string            `json:"site_id"`
	Scanned            bool              `json:"scanned"`
	Alive              bool              `json:"alive"`
	FirstSeen          int64             `json:"first_seen"`
	LastSeen           int64             `json:"last_seen"`
	DetectedBy         string            `json:"detected_by"`
	Type               string            `json:"type"`
	OS                 string            `json:"os"`
	OSVersion          string            `json:"os_version"`
	HW                 string            `json:"hw"`
	Addresses          []string          `json:"addresses"`
	AddressesExtra     []string          `json:"addresses_extra"`
	MACs               []string          `json:"macs"`
	MACVendors         []string          `json:"mac_vendors"`
	Names              []string          `json:"names"`
	Domains            []string          `json:"domains"`
	ServiceCount       int64             `json:"service_count"`
	ServiceCountTCP    int64             `json:"service_count_tcp"`
	ServiceCountUDP    int64             `json:"service_count_udp"`
	ServiceCountARP    int64             `json:"service_count_arp"`
	ServiceCountICMP   int64             `json:"service_count_icmp"`
	LowestTTL          int64             `json:"lowest_ttl"`
	LowestRTT          int64             `json:"lowest_rtt"`
	LastAgentID        string            `json:"last_agent_id"`
	LastTaskID         string            `json:"last_task_id"`
	NewestMAC          string            `json:"newest_mac"`
	NewestMACVendor    string            `json:"newest_mac_vendor"`
	NewestMACAge       int64             `json:"newest_mac_age"`
	Comments           string            `json:"comments"`
	Criticality        string            `json:"criticality"`
	ServicePortsTCP    []string          `json:"service_ports_tcp"`
	ServicePortsUDP    []string          `json:"service_ports_udp"`
	ServiceProtocols   []string          `json:"service_protocols"`
	ServiceProducts    []string          `json:"service_products"`
	OrgName            string            `json:"org_name"`
	SiteName           string            `json:"site_name"`
	AgentName          string            `json:"agent_name"`
	Tags               map[string]string `json:"tags"`
	Attributes         map[string]string `json:"attributes"`
	RTTs               map[string]any    `json:"rtts"`
	Credentials        map[string]any    `json:"credentials"`
	Services           map[string]any    `json:"services"`
	RiskRank           float64           `json:"risk_rank"`
	Risk               string            `json:"risk"`
	VulnerabilityCount float64           `json:"vulnerability_count"`
	FindingCount       float64           `json:"finding_count"`
	EolOS              float64           `json:"eol_os"`
	EolOSExt           float64           `json:"eol_os_ext"`
	Sources            []string          `json:"sources"`
	CreatedAt          int64             `json:"created_at"`
	UpdatedAt          int64             `json:"updated_at"`
}

func (c *Client) ListAssets(ctx context.Context, search string) ([]Asset, error) {
	path := "/org/assets"
	if search != "" {
		path = fmt.Sprintf("/org/assets?search=%s", url.QueryEscape(search))
	}
	resp, err := c.DoRequest(ctx, http.MethodGet, path, AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var assets []Asset
	if err := HandleResponse(resp, &assets); err != nil {
		return nil, fmt.Errorf("list assets: %w", err)
	}
	return assets, nil
}

func (c *Client) GetAsset(ctx context.Context, id string) (*Asset, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/assets/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var asset Asset
	if err := HandleResponse(resp, &asset); err != nil {
		return nil, fmt.Errorf("get asset: %w", err)
	}
	return &asset, nil
}

func (c *Client) DeleteAsset(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/org/assets/%s", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

type UpdateAssetCommentsRequest struct {
	Comments string `json:"comments"`
}

func (c *Client) UpdateAssetComments(ctx context.Context, id string, req UpdateAssetCommentsRequest) (*Asset, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/assets/%s/comments", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var asset Asset
	if err := HandleResponse(resp, &asset); err != nil {
		return nil, fmt.Errorf("update asset comments: %w", err)
	}
	return &asset, nil
}

type UpdateAssetTagsRequest struct {
	Tags string `json:"tags"`
}

func (c *Client) UpdateAssetTags(ctx context.Context, id string, req UpdateAssetTagsRequest) (*Asset, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/assets/%s/tags", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var asset Asset
	if err := HandleResponse(resp, &asset); err != nil {
		return nil, fmt.Errorf("update asset tags: %w", err)
	}
	return &asset, nil
}

type UpdateAssetCriticalityRequest struct {
	Criticality string `json:"criticality"`
}

func (c *Client) UpdateAssetCriticality(ctx context.Context, id string, req UpdateAssetCriticalityRequest) (*Asset, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/assets/%s/criticality", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var asset Asset
	if err := HandleResponse(resp, &asset); err != nil {
		return nil, fmt.Errorf("update asset criticality: %w", err)
	}
	return &asset, nil
}

type AssetOwnership struct {
	OwnershipTypeID string `json:"ownership_type_id"`
	Owner           string `json:"owner"`
}

type UpdateAssetOwnershipsRequest struct {
	Ownerships []AssetOwnership `json:"ownerships"`
}

func (c *Client) UpdateAssetOwnerships(ctx context.Context, id string, req UpdateAssetOwnershipsRequest) (*Asset, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/assets/%s/owners", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var asset Asset
	if err := HandleResponse(resp, &asset); err != nil {
		return nil, fmt.Errorf("update asset ownerships: %w", err)
	}
	return &asset, nil
}
