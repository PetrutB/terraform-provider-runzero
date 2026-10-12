package client

import (
	"context"
	"fmt"
	"net/http"
)

type Agent struct {
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
	SystemInfo     map[string]any    `json:"system_info"`
}

type UpdateAgentRequest struct {
	SiteID   string              `json:"site_id,omitempty"`
	Settings *AgentSettingsInput `json:"settings,omitempty"`
}

type AgentSettingsInput struct {
	MaxConcurrentScans int64 `json:"max_concurrent_scans,omitempty"`
}

func (c *Client) ListAgents(ctx context.Context) ([]Agent, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org/agents", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var agents []Agent
	if err := HandleResponse(resp, &agents); err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	return agents, nil
}

func (c *Client) GetAgent(ctx context.Context, id string) (*Agent, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/agents/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := HandleResponse(resp, &agent); err != nil {
		return nil, fmt.Errorf("get agent: %w", err)
	}
	return &agent, nil
}

func (c *Client) UpdateAgent(ctx context.Context, id string, req UpdateAgentRequest) (*Agent, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/agents/%s", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := HandleResponse(resp, &agent); err != nil {
		return nil, fmt.Errorf("update agent: %w", err)
	}
	return &agent, nil
}

func (c *Client) DeleteAgent(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/org/agents/%s", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) UpgradeAgent(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/org/agents/%s/update", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}
