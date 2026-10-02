package client

import (
	"context"
	"fmt"
	"net/http"
)

type ScanTemplate struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	ClientID            string            `json:"client_id"`
	OrganizationID      string            `json:"organization_id"`
	AgentID             string            `json:"agent_id"`
	SiteID              string            `json:"site_id"`
	CruncherID          string            `json:"cruncher_id"`
	CreatedAt           int64             `json:"created_at"`
	CreatedBy           string            `json:"created_by"`
	CreatedByUserID     string            `json:"created_by_user_id"`
	UpdatedAt           int64             `json:"updated_at"`
	Type                string            `json:"type"`
	Status              string            `json:"status"`
	Error               string            `json:"error"`
	Params              map[string]string `json:"params"`
	Stats               map[string]any    `json:"stats"`
	Hidden              bool              `json:"hidden"`
	ParentID            string            `json:"parent_id"`
	Recur               bool              `json:"recur"`
	RecurFrequency      string            `json:"recur_frequency"`
	StartTime           int64             `json:"start_time"`
	RecurLast           int64             `json:"recur_last"`
	RecurNext           int64             `json:"recur_next"`
	RecurLastTaskID     string            `json:"recur_last_task_id"`
	GracePeriod         string            `json:"grace_period"`
	CustomIntegrationID string            `json:"custom_integration_id"`
	SourceID            string            `json:"source_id"`
	TemplateID          string            `json:"template_id"`
	SizeSite            int64             `json:"size_site"`
	SizeData            int64             `json:"size_data"`
	SizeResults         int64             `json:"size_results"`
	HostedZoneID        string            `json:"hosted_zone_id"`
	LinkedTaskCount     int32             `json:"linked_task_count"`
	Global              bool              `json:"global"`
	ACL                 map[string]any    `json:"acl"`
}

type CreateScanTemplateRequest struct {
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	OrganizationID string            `json:"organization_id"`
	Params         map[string]string `json:"params,omitempty"`
	Global         bool              `json:"global"`
	ACL            map[string]any    `json:"acl"`
}

type UpdateScanTemplateRequest struct {
	ID             string            `json:"id"`
	Name           string            `json:"name,omitempty"`
	Description    string            `json:"description,omitempty"`
	OrganizationID string            `json:"organization_id"`
	Params         map[string]string `json:"params,omitempty"`
	Global         bool              `json:"global"`
	ACL            map[string]any    `json:"acl"`
}

func (c *Client) CreateScanTemplate(ctx context.Context, req CreateScanTemplateRequest) (*ScanTemplate, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, "/account/tasks/templates", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var template ScanTemplate
	if err := HandleResponse(resp, &template); err != nil {
		return nil, fmt.Errorf("create scan template: %w", err)
	}
	return &template, nil
}

func (c *Client) GetScanTemplate(ctx context.Context, id string) (*ScanTemplate, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/tasks/templates/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var template ScanTemplate
	if err := HandleResponse(resp, &template); err != nil {
		return nil, fmt.Errorf("get scan template: %w", err)
	}
	return &template, nil
}

func (c *Client) UpdateScanTemplate(ctx context.Context, req UpdateScanTemplateRequest) (*ScanTemplate, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/tasks/templates", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var template ScanTemplate
	if err := HandleResponse(resp, &template); err != nil {
		return nil, fmt.Errorf("update scan template: %w", err)
	}
	return &template, nil
}

func (c *Client) DeleteScanTemplate(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/tasks/templates/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListScanTemplates(ctx context.Context) ([]ScanTemplate, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/tasks/templates", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var templates []ScanTemplate
	if err := HandleResponse(resp, &templates); err != nil {
		return nil, fmt.Errorf("list scan templates: %w", err)
	}
	return templates, nil
}
