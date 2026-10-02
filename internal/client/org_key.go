package client

import (
	"context"
	"fmt"
	"net/http"
)

// GetOrgKey returns the API key for the organization tied to the org-scoped token (GET /org/key).
func (c *Client) GetOrgKey(ctx context.Context) (*APIKey, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org/key", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var key APIKey
	if err := HandleResponse(resp, &key); err != nil {
		return nil, fmt.Errorf("get org key: %w", err)
	}
	return &key, nil
}

// RotateOrgKey rotates the organization API key and returns the new value (PATCH /org/key/rotate).
func (c *Client) RotateOrgKey(ctx context.Context) (*APIKey, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, "/org/key/rotate", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var key APIKey
	if err := HandleResponse(resp, &key); err != nil {
		return nil, fmt.Errorf("rotate org key: %w", err)
	}
	return &key, nil
}

// DeleteOrgKey removes the organization API key (DELETE /org/key).
func (c *Client) DeleteOrgKey(ctx context.Context) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, "/org/key", AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}
