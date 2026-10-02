package client

import (
	"context"
	"fmt"
	"net/http"
)

type HostedZone struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Enabled              bool   `json:"enabled"`
	UpdatedAt            string `json:"updated_at"`
	ProcessorID          string `json:"processor_id"`
	ExplorersConcurrency int64  `json:"explorers_concurrency"`
	ExplorersTotal       int64  `json:"explorers_total"`
	TasksActive          int64  `json:"tasks_active"`
	TasksWaiting         int64  `json:"tasks_waiting"`
	OrganizationID       string `json:"organization_id"`
}

func (c *Client) GetHostedZone(ctx context.Context, id string) (*HostedZone, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/hosted-zones/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var zone HostedZone
	if err := HandleResponse(resp, &zone); err != nil {
		return nil, fmt.Errorf("get hosted zone: %w", err)
	}
	return &zone, nil
}

func (c *Client) ListHostedZones(ctx context.Context) ([]HostedZone, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org/hosted-zones", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var zones []HostedZone
	if err := HandleResponse(resp, &zones); err != nil {
		return nil, fmt.Errorf("list hosted zones: %w", err)
	}
	return zones, nil
}
