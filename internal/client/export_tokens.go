package client

import (
	"context"
	"fmt"
	"net/http"
)

type ExportToken struct {
	ID         string `json:"id"`
	Token      string `json:"token"`
	Comment    string `json:"comment"`
	CreatedAt  int64  `json:"created_at"`
	LastUsedAt int64  `json:"last_used_at"`
	LastUsedBy string `json:"last_used_by"`
	Counter    int64  `json:"counter"`
}

type CreateExportTokenRequest struct {
	Comment string `json:"comment,omitempty"`
}

func (c *Client) CreateExportToken(ctx context.Context, orgID string, req CreateExportTokenRequest) (*ExportToken, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/account/orgs/%s/exportTokens", orgID), AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var token ExportToken
	if err := HandleResponse(resp, &token); err != nil {
		return nil, fmt.Errorf("create export token: %w", err)
	}
	return &token, nil
}

func (c *Client) GetExportToken(ctx context.Context, orgID, id string) (*ExportToken, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/orgs/%s/exportTokens/%s", orgID, id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var token ExportToken
	if err := HandleResponse(resp, &token); err != nil {
		return nil, fmt.Errorf("get export token: %w", err)
	}
	return &token, nil
}

func (c *Client) DeleteExportToken(ctx context.Context, orgID, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/orgs/%s/exportTokens/%s", orgID, id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListExportTokens(ctx context.Context, orgID string) ([]ExportToken, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/orgs/%s/exportTokens", orgID), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var tokens []ExportToken
	if err := HandleResponse(resp, &tokens); err != nil {
		return nil, fmt.Errorf("list export tokens: %w", err)
	}
	return tokens, nil
}

func (c *Client) RotateExportToken(ctx context.Context, orgID, id string) (*ExportToken, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/orgs/%s/exportTokens/%s/rotate", orgID, id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var token ExportToken
	if err := HandleResponse(resp, &token); err != nil {
		return nil, fmt.Errorf("rotate export token: %w", err)
	}
	return &token, nil
}
