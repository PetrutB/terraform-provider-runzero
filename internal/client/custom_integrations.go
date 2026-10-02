package client

import (
	"context"
	"fmt"
	"net/http"
)

type CustomIntegration struct {
	ID          string `json:"id"`
	ClientID    string `json:"clientId"`
	CreatedByID string `json:"createdById"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	DestroyedAt string `json:"destroyedAt"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type CreateCustomIntegrationRequest struct {
	Name        string `json:"name"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateCustomIntegrationRequest struct {
	Name        string `json:"name,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *Client) CreateCustomIntegration(ctx context.Context, req CreateCustomIntegrationRequest) (*CustomIntegration, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, "/account/custom-integrations", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var integration CustomIntegration
	if err := HandleResponse(resp, &integration); err != nil {
		return nil, fmt.Errorf("create custom integration: %w", err)
	}
	return &integration, nil
}

func (c *Client) GetCustomIntegration(ctx context.Context, id string) (*CustomIntegration, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/custom-integrations/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var integration CustomIntegration
	if err := HandleResponse(resp, &integration); err != nil {
		return nil, fmt.Errorf("get custom integration: %w", err)
	}
	return &integration, nil
}

func (c *Client) UpdateCustomIntegration(ctx context.Context, id string, req UpdateCustomIntegrationRequest) (*CustomIntegration, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/custom-integrations/%s", id), AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var integration CustomIntegration
	if err := HandleResponse(resp, &integration); err != nil {
		return nil, fmt.Errorf("update custom integration: %w", err)
	}
	return &integration, nil
}

func (c *Client) DeleteCustomIntegration(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/custom-integrations/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListCustomIntegrations(ctx context.Context) ([]CustomIntegration, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/custom-integrations", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var integrations []CustomIntegration
	if err := HandleResponse(resp, &integrations); err != nil {
		return nil, fmt.Errorf("list custom integrations: %w", err)
	}
	return integrations, nil
}
