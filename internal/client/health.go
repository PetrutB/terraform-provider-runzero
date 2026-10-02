package client

import (
	"context"
	"fmt"
	"net/http"
)

type HealthCheck struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (c *Client) GetHealth(ctx context.Context) (*HealthCheck, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/health", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var health HealthCheck
	if err := HandleResponse(resp, &health); err != nil {
		return nil, fmt.Errorf("get health: %w", err)
	}
	return &health, nil
}
