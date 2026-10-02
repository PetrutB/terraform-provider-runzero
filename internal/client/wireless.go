package client

import (
	"context"
	"fmt"
	"net/http"
)

type WirelessLAN struct {
	ID             string            `json:"id"`
	OrganizationID string            `json:"organization_id"`
	SiteID         string            `json:"site_id"`
	CreatedAt      int64             `json:"created_at"`
	LastSeen       int64             `json:"last_seen"`
	LastAgentID    string            `json:"last_agent_id"`
	LastTaskID     string            `json:"last_task_id"`
	ESSID          string            `json:"essid"`
	BSSID          string            `json:"bssid"`
	Type           string            `json:"type"`
	Authentication string            `json:"authentication"`
	Encryption     string            `json:"encryption"`
	Signal         int               `json:"signal"`
	Channels       string            `json:"channels"`
	Interface      string            `json:"interface"`
	Vendor         string            `json:"vendor"`
	Family         string            `json:"family"`
	Data           map[string]string `json:"data"`
	OrgName        string            `json:"org_name"`
	SiteName       string            `json:"site_name"`
	AgentName      string            `json:"agent_name"`
}

func (c *Client) ListWirelessLANs(ctx context.Context) ([]WirelessLAN, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org/wireless", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var lans []WirelessLAN
	if err := HandleResponse(resp, &lans); err != nil {
		return nil, fmt.Errorf("list wireless LANs: %w", err)
	}
	return lans, nil
}

func (c *Client) GetWirelessLAN(ctx context.Context, id string) (*WirelessLAN, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/wireless/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var lan WirelessLAN
	if err := HandleResponse(resp, &lan); err != nil {
		return nil, fmt.Errorf("get wireless LAN: %w", err)
	}
	return &lan, nil
}

func (c *Client) DeleteWirelessLAN(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/org/wireless/%s", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}
