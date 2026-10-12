package client

import (
	"context"
	"fmt"
	"net/http"
)

type Credential struct {
	ID             string            `json:"id"`
	ClientID       string            `json:"client_id"`
	Name           string            `json:"name"`
	Type           string            `json:"type"`
	CreatedAt      int64             `json:"created_at"`
	CreatedByID    string            `json:"created_by_id"`
	CreatedByEmail string            `json:"created_by_email"`
	ACL            map[string]any    `json:"acl"`
	Global         bool              `json:"global"`
	CIDRs          []string          `json:"cidrs"`
	LastUsedAt     int64             `json:"last_used_at"`
	LastUsedByID   string            `json:"last_used_by_id"`
}

type CreateCredentialRequest struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	ACL    map[string]string `json:"acl,omitempty"`
	Global bool              `json:"global,omitempty"`
	CIDRs  string            `json:"cidrs,omitempty"`
	// Type-specific secret material; keys depend on Type (e.g. AWS: access_key/secret_access_key).
	Secret map[string]string `json:"secret,omitempty"`
}

func (c *Client) CreateCredential(ctx context.Context, req CreateCredentialRequest) (*Credential, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/credentials", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var credential Credential
	if err := HandleResponse(resp, &credential); err != nil {
		return nil, fmt.Errorf("create credential: %w", err)
	}
	return &credential, nil
}

func (c *Client) GetCredential(ctx context.Context, id string) (*Credential, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/credentials/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var credential Credential
	if err := HandleResponse(resp, &credential); err != nil {
		return nil, fmt.Errorf("get credential: %w", err)
	}
	return &credential, nil
}

func (c *Client) DeleteCredential(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/credentials/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListCredentials(ctx context.Context) ([]Credential, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/credentials", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var credentials []Credential
	if err := HandleResponse(resp, &credentials); err != nil {
		return nil, fmt.Errorf("list credentials: %w", err)
	}
	return credentials, nil
}
