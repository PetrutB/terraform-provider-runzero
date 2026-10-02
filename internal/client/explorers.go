package client

import (
	"context"
	"fmt"
	"net/http"
)

type Explorer struct {
	ID             string            `json:"id"`
	ClientID       string            `json:"client_id"`
	OrganizationID string            `json:"organization_id"`
	SiteID         string            `json:"site_id"`
	Name           string            `json:"name"`
	HubID          string            `json:"hub_id"`
	HostID         string            `json:"host_id"`
	OS             string            `json:"os"`
	Arch           string            `json:"arch"`
	Version        string            `json:"version"`
	ExternalIP     string            `json:"external_ip"`
	InternalIP     string            `json:"internal_ip"`
	Connected      bool              `json:"connected"`
	Inactive       bool              `json:"inactive"`
	LastCheckin    int64             `json:"last_checkin"`
	CreatedAt      int64             `json:"created_at"`
	UpdatedAt      int64             `json:"updated_at"`
	DeactivatedAt  int64             `json:"deactivated_at"`
	SystemInfo     map[string]string `json:"system_info"`
}

type UpdateExplorerRequest struct {
	SiteID   string              `json:"site_id,omitempty"`
	Settings *AgentSettingsInput `json:"settings,omitempty"`
}

func (c *Client) ListExplorers(ctx context.Context) ([]Explorer, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org/explorers", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var explorers []Explorer
	if err := HandleResponse(resp, &explorers); err != nil {
		return nil, fmt.Errorf("list explorers: %w", err)
	}
	return explorers, nil
}

func (c *Client) GetExplorer(ctx context.Context, id string) (*Explorer, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/explorers/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var explorer Explorer
	if err := HandleResponse(resp, &explorer); err != nil {
		return nil, fmt.Errorf("get explorer: %w", err)
	}
	return &explorer, nil
}

func (c *Client) UpdateExplorer(ctx context.Context, id string, req UpdateExplorerRequest) (*Explorer, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/explorers/%s", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var explorer Explorer
	if err := HandleResponse(resp, &explorer); err != nil {
		return nil, fmt.Errorf("update explorer: %w", err)
	}
	return &explorer, nil
}

func (c *Client) DeleteExplorer(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/org/explorers/%s", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) UpgradeExplorer(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/org/explorers/%s/update", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}
