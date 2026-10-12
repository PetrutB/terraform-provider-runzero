package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Task struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	TemplateID          string            `json:"template_id"`
	ClientID            string            `json:"client_id"`
	OrganizationID      string            `json:"organization_id"`
	AgentID             string            `json:"agent_id"`
	HostedZoneID        string            `json:"hosted_zone_id"`
	SiteID              string            `json:"site_id"`
	CruncherID          string            `json:"cruncher_id"`
	CreatedAt           int64             `json:"created_at"`
	CreatedBy           string            `json:"created_by"`
	CreatedByUserID     string            `json:"created_by_user_id"`
	CustomIntegrationID string            `json:"custom_integration_id"`
	SourceID            int64             `json:"source_id"`
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
}

func (c *Client) ListTasks(ctx context.Context, status string) ([]Task, error) {
	path := "/org/tasks"
	if status != "" {
		path = fmt.Sprintf("/org/tasks?status=%s", url.QueryEscape(status))
	}
	resp, err := c.DoRequest(ctx, http.MethodGet, path, AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := HandleResponse(resp, &tasks); err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}

func (c *Client) GetTask(ctx context.Context, id string) (*Task, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/org/tasks/%s", id), AuthOrg, nil)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := HandleResponse(resp, &task); err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	return &task, nil
}

type UpdateTaskRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *Client) UpdateTask(ctx context.Context, id string, req UpdateTaskRequest) (*Task, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/org/tasks/%s", id), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := HandleResponse(resp, &task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	return &task, nil
}

func (c *Client) StopTask(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/org/tasks/%s/stop", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) HideTask(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/org/tasks/%s/hide", id), AuthOrg, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}
