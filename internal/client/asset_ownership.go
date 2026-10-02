package client

import (
	"context"
	"fmt"
	"net/http"
)

type AssetOwnershipType struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Reference int64  `json:"reference"`
	Order     int64  `json:"order"`
	Hidden    bool   `json:"hidden"`
}

type CreateAssetOwnershipTypeRequest struct {
	Name      string `json:"name"`
	Reference int64  `json:"reference,omitempty"`
	Order     int64  `json:"order,omitempty"`
	Hidden    bool   `json:"hidden,omitempty"`
}

type UpdateAssetOwnershipTypeRequest struct {
	Name      string `json:"name,omitempty"`
	Reference int64  `json:"reference,omitempty"`
	Order     int64  `json:"order,omitempty"`
	Hidden    bool   `json:"hidden,omitempty"`
}

func (c *Client) CreateAssetOwnershipType(ctx context.Context, req CreateAssetOwnershipTypeRequest) (*AssetOwnershipType, error) {
	// The API expects and returns a JSON array of ownership types for this bulk-create endpoint.
	resp, err := c.DoRequest(ctx, http.MethodPost, "/account/assets/ownership-types", AuthAccount, []CreateAssetOwnershipTypeRequest{req})
	if err != nil {
		return nil, err
	}
	var created []AssetOwnershipType
	if err := HandleResponse(resp, &created); err != nil {
		return nil, fmt.Errorf("create asset ownership type: %w", err)
	}
	if len(created) == 0 {
		return nil, fmt.Errorf("create asset ownership type: empty response")
	}
	for i := range created {
		if created[i].Name == req.Name {
			return &created[i], nil
		}
	}
	return &created[len(created)-1], nil
}

func (c *Client) GetAssetOwnershipType(ctx context.Context, id string) (*AssetOwnershipType, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/assets/ownership-types", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var types []AssetOwnershipType
	if err := HandleResponse(resp, &types); err != nil {
		return nil, fmt.Errorf("list asset ownership types: %w", err)
	}
	for _, t := range types {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, &APIError{
		StatusCode: 404,
		Message:    "Not Found",
		Body:       fmt.Sprintf("asset ownership type %s not found", id),
	}
}

func (c *Client) UpdateAssetOwnershipType(ctx context.Context, id string, req UpdateAssetOwnershipTypeRequest) (*AssetOwnershipType, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/assets/ownership-types/%s", id), AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var ownershipType AssetOwnershipType
	if err := HandleResponse(resp, &ownershipType); err != nil {
		return nil, fmt.Errorf("update asset ownership type: %w", err)
	}
	return &ownershipType, nil
}

func (c *Client) DeleteAssetOwnershipType(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/assets/ownership-types/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListAssetOwnershipTypes(ctx context.Context) ([]AssetOwnershipType, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/assets/ownership-types", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var types []AssetOwnershipType
	if err := HandleResponse(resp, &types); err != nil {
		return nil, fmt.Errorf("list asset ownership types: %w", err)
	}
	return types, nil
}
