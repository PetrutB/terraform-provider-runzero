package client

import (
	"context"
	"fmt"
	"net/http"
)

type Group struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	RoleSummary    string            `json:"role_summary"`
	UserCount      int64             `json:"user_count"`
	CreatedByEmail string            `json:"created_by_email"`
	ExpiresAt      int64             `json:"expires_at"`
	OrgDefaultRole string            `json:"org_default_role"`
	OrgRoles       map[string]string `json:"org_roles"`
	CreatedAt      int64             `json:"created_at"`
	UpdatedAt      int64             `json:"updated_at"`
}

type CreateGroupRequest struct {
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	ExpiresAt      int64             `json:"expires_at,omitempty"`
	OrgDefaultRole string            `json:"org_default_role,omitempty"`
	OrgRoles       map[string]string `json:"org_roles,omitempty"`
}

type UpdateGroupRequest struct {
	ID             string            `json:"id"`
	Name           string            `json:"name,omitempty"`
	Description    string            `json:"description,omitempty"`
	ExpiresAt      int64             `json:"expires_at,omitempty"`
	OrgDefaultRole string            `json:"org_default_role,omitempty"`
	OrgRoles       map[string]string `json:"org_roles,omitempty"`
}

func (c *Client) CreateGroup(ctx context.Context, req CreateGroupRequest) (*Group, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, "/account/groups", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var group Group
	if err := HandleResponse(resp, &group); err != nil {
		return nil, fmt.Errorf("create group: %w", err)
	}
	return &group, nil
}

func (c *Client) GetGroup(ctx context.Context, id string) (*Group, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/groups/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var group Group
	if err := HandleResponse(resp, &group); err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}
	return &group, nil
}

func (c *Client) UpdateGroup(ctx context.Context, req UpdateGroupRequest) (*Group, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/groups", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var group Group
	if err := HandleResponse(resp, &group); err != nil {
		return nil, fmt.Errorf("update group: %w", err)
	}
	return &group, nil
}

func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/groups/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListGroups(ctx context.Context) ([]Group, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/groups", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var groups []Group
	if err := HandleResponse(resp, &groups); err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	return groups, nil
}
