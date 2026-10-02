package client

import (
	"context"
	"fmt"
	"net/http"
)

type SSOGroupMapping struct {
	ID             string `json:"id"`
	GroupID        string `json:"group_id"`
	GroupName      string `json:"group_name"`
	SSOAttribute   string `json:"sso_attribute"`
	SSOValue       string `json:"sso_value"`
	Description    string `json:"description"`
	CreatedByEmail string `json:"created_by_email"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

type CreateSSOGroupMappingRequest struct {
	GroupID      string `json:"group_id"`
	SSOAttribute string `json:"sso_attribute"`
	SSOValue     string `json:"sso_value"`
	Description  string `json:"description,omitempty"`
}

type UpdateSSOGroupMappingRequest struct {
	ID           string `json:"id"`
	GroupID      string `json:"group_id"`
	SSOAttribute string `json:"sso_attribute"`
	SSOValue     string `json:"sso_value"`
	Description  string `json:"description,omitempty"`
}

func (c *Client) CreateSSOGroupMapping(ctx context.Context, req CreateSSOGroupMappingRequest) (*SSOGroupMapping, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, "/account/sso/groups", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var mapping SSOGroupMapping
	if err := HandleResponse(resp, &mapping); err != nil {
		return nil, fmt.Errorf("create SSO group mapping: %w", err)
	}
	return &mapping, nil
}

func (c *Client) GetSSOGroupMapping(ctx context.Context, id string) (*SSOGroupMapping, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/sso/groups/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var mapping SSOGroupMapping
	if err := HandleResponse(resp, &mapping); err != nil {
		return nil, fmt.Errorf("get SSO group mapping: %w", err)
	}
	return &mapping, nil
}

func (c *Client) UpdateSSOGroupMapping(ctx context.Context, req UpdateSSOGroupMappingRequest) (*SSOGroupMapping, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/sso/groups", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var mapping SSOGroupMapping
	if err := HandleResponse(resp, &mapping); err != nil {
		return nil, fmt.Errorf("update SSO group mapping: %w", err)
	}
	return &mapping, nil
}

func (c *Client) DeleteSSOGroupMapping(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/sso/groups/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListSSOGroupMappings(ctx context.Context) ([]SSOGroupMapping, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/sso/groups", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var mappings []SSOGroupMapping
	if err := HandleResponse(resp, &mappings); err != nil {
		return nil, fmt.Errorf("list SSO group mappings: %w", err)
	}
	return mappings, nil
}
