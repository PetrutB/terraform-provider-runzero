package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Service struct {
	ID                    string            `json:"id"`
	OrganizationID        string            `json:"organization_id"`
	SiteID                string            `json:"site_id"`
	CreatedAt             int64             `json:"created_at"`
	UpdatedAt             int64             `json:"updated_at"`
	ServiceID             string            `json:"service_id"`
	ServiceAssetID        string            `json:"service_asset_id"`
	ServiceCreatedAt      int64             `json:"service_created_at"`
	ServiceUpdatedAt      int64             `json:"service_updated_at"`
	ServiceAddress        string            `json:"service_address"`
	ServiceTransport      string            `json:"service_transport"`
	ServiceVhost          string            `json:"service_vhost"`
	ServicePort           string            `json:"service_port"`
	ServiceProtocol       string            `json:"service_protocol"`
	ServiceSummary        string            `json:"service_summary"`
	ServiceScreenshotLink string            `json:"service_screenshot_link"`
	ServiceLink           string            `json:"service_link"`
	ServiceData           map[string]string `json:"service_data"`
	Alive                 bool              `json:"alive"`
	FirstSeen             int64             `json:"first_seen"`
	LastSeen              int64             `json:"last_seen"`
	DetectedBy            string            `json:"detected_by"`
	Type                  string            `json:"type"`
	OS                    string            `json:"os"`
	OSVersion             string            `json:"os_version"`
	HW                    string            `json:"hw"`
	Addresses             []string          `json:"addresses"`
	AddressesExtra        []string          `json:"addresses_extra"`
	MACs                  []string          `json:"macs"`
	MACVendors            []string          `json:"mac_vendors"`
	Names                 []string          `json:"names"`
	Domains               []string          `json:"domains"`
	ServiceCount          int64             `json:"service_count"`
	ServiceCountTCP       int64             `json:"service_count_tcp"`
	ServiceCountUDP       int64             `json:"service_count_udp"`
	ServiceCountARP       int64             `json:"service_count_arp"`
	ServiceCountICMP      int64             `json:"service_count_icmp"`
	LowestTTL             int64             `json:"lowest_ttl"`
	LowestRTT             int64             `json:"lowest_rtt"`
	LastAgentID           string            `json:"last_agent_id"`
	LastTaskID            string            `json:"last_task_id"`
	NewestMAC             string            `json:"newest_mac"`
	NewestMACVendor       string            `json:"newest_mac_vendor"`
	NewestMACAge          int64             `json:"newest_mac_age"`
	Comments              string            `json:"comments"`
	ServicePortsTCP       []string          `json:"service_ports_tcp"`
	ServicePortsUDP       []string          `json:"service_ports_udp"`
	ServicePortsProtocols []string          `json:"service_ports_protocols"`
	ServicePortsProducts  []string          `json:"service_ports_products"`
	OrgName               string            `json:"org_name"`
	SiteName              string            `json:"site_name"`
	AgentName             string            `json:"agent_name"`
	Tags                  map[string]string `json:"tags"`
	RTTs                  map[string]any    `json:"rtts"`
	Credentials           map[string]string `json:"credentials"`
	Services              map[string]any    `json:"services"`
	Attributes            map[string]any    `json:"attributes"`
}

func (c *Client) ListServices(ctx context.Context, search string) ([]Service, error) {
	path := "/org/services"
	if search != "" {
		path = fmt.Sprintf("/org/services?search=%s", url.QueryEscape(search))
	}
	resp, err := c.DoRequest(ctx, http.MethodGet, path, AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var services []Service
	if err := HandleResponse(resp, &services); err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	return services, nil
}

func (c *Client) GetService(ctx context.Context, id string) (*Service, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/services/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var service Service
	if err := HandleResponse(resp, &service); err != nil {
		return nil, fmt.Errorf("get service: %w", err)
	}
	return &service, nil
}

func (c *Client) DeleteService(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/org/services/%s", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}
