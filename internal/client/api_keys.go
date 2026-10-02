package client

import (
	"context"
	"fmt"
	"net/http"
)

type APIKey struct {
	ID             string `json:"id"`
	ClientID       string `json:"client_id"`
	OrganizationID string `json:"organization_id"`
	CreatedAt      int64  `json:"created_at"`
	CreatedBy      string `json:"created_by"`
	Comment        string `json:"comment"`
	LastUsedAt     int64  `json:"last_used_at"`
	LastUsedIP     string `json:"last_used_ip"`
	LastUsedUA     string `json:"last_used_ua"`
	Counter        int64  `json:"counter"`
	UsageToday     int64  `json:"usage_today"`
	UsageLimit     int64  `json:"usage_limit"`
	Token          string `json:"token"`
	Inactive       bool   `json:"inactive"`
	Type           string `json:"type"`
}

type CreateAPIKeyRequest struct {
	Comment        string `json:"comment,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*APIKey, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/keys", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var apiKey APIKey
	if err := HandleResponse(resp, &apiKey); err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}
	return &apiKey, nil
}

func (c *Client) GetAPIKey(ctx context.Context, id string) (*APIKey, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/keys/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var apiKey APIKey
	if err := HandleResponse(resp, &apiKey); err != nil {
		return nil, fmt.Errorf("get api key: %w", err)
	}
	return &apiKey, nil
}

func (c *Client) DeleteAPIKey(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/keys/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/keys", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var apiKeys []APIKey
	if err := HandleResponse(resp, &apiKeys); err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	return apiKeys, nil
}

func (c *Client) RotateAPIKey(ctx context.Context, id string) (*APIKey, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/keys/%s/rotate", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var apiKey APIKey
	if err := HandleResponse(resp, &apiKey); err != nil {
		return nil, fmt.Errorf("rotate api key: %w", err)
	}
	return &apiKey, nil
}
