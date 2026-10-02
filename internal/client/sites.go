package client

import (
	"context"
	"fmt"
	"net/http"
)

type Site struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Scope       string         `json:"scope"`
	Excludes    string         `json:"excludes"`
	Subnets     map[string]any `json:"subnets"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
}

type CreateSiteRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Excludes    string `json:"excludes,omitempty"`
}

type UpdateSiteRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Excludes    string `json:"excludes,omitempty"`
}

func (c *Client) CreateSite(ctx context.Context, req CreateSiteRequest) (*Site, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/org/sites", AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var site Site
	if err := HandleResponse(resp, &site); err != nil {
		return nil, fmt.Errorf("create site: %w", err)
	}
	return &site, nil
}

func (c *Client) GetSite(ctx context.Context, id string) (*Site, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/sites/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var site Site
	if err := HandleResponse(resp, &site); err != nil {
		return nil, fmt.Errorf("get site: %w", err)
	}
	return &site, nil
}

func (c *Client) UpdateSite(ctx context.Context, id string, req UpdateSiteRequest) (*Site, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/sites/%s", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var site Site
	if err := HandleResponse(resp, &site); err != nil {
		return nil, fmt.Errorf("update site: %w", err)
	}
	return &site, nil
}

func (c *Client) DeleteSite(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/org/sites/%s", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/org/sites", AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var sites []Site
	if err := HandleResponse(resp, &sites); err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}
	return sites, nil
}
