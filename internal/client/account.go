package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListAccountSites(ctx context.Context, search string) ([]Site, error) {
	path := "/account/sites"
	if search != "" {
		path = fmt.Sprintf("/account/sites?search=%s", url.QueryEscape(search))
	}
	resp, err := c.DoRequest(ctx, http.MethodGet, path, AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var sites []Site
	if err := HandleResponse(resp, &sites); err != nil {
		return nil, fmt.Errorf("list account sites: %w", err)
	}
	return sites, nil
}

func (c *Client) ListAccountTasks(ctx context.Context, search string) ([]Task, error) {
	path := "/account/tasks"
	if search != "" {
		path = fmt.Sprintf("/account/tasks?search=%s", url.QueryEscape(search))
	}
	resp, err := c.DoRequest(ctx, http.MethodGet, path, AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := HandleResponse(resp, &tasks); err != nil {
		return nil, fmt.Errorf("list account tasks: %w", err)
	}
	return tasks, nil
}

type UserInviteRequest struct {
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	Email          string            `json:"email"`
	ClientAdmin    bool              `json:"client_admin"`
	OrgDefaultRole string            `json:"org_default_role"`
	OrgRoles       map[string]string `json:"org_roles,omitempty"`
	Subject        string            `json:"subject,omitempty"`
	Message        string            `json:"message,omitempty"`
}

func (c *Client) InviteUser(ctx context.Context, req UserInviteRequest) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/users/invite", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("invite user: %w", err)
	}
	return &user, nil
}

func (c *Client) ResetUserMFA(ctx context.Context, userID string) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/users/%s/resetMFA", userID), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("reset user MFA: %w", err)
	}
	return &user, nil
}

func (c *Client) ResetUserLockout(ctx context.Context, userID string) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/users/%s/resetLockout", userID), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("reset user lockout: %w", err)
	}
	return &user, nil
}

func (c *Client) ResetUserPassword(ctx context.Context, userID string) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/users/%s/resetPassword", userID), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("reset user password: %w", err)
	}
	return &user, nil
}
