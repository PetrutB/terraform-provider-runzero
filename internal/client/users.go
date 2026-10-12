package client

import (
	"context"
	"fmt"
	"net/http"
)

type User struct {
	ID                    string         `json:"id"`
	ClientID              string         `json:"client_id"`
	Name                  string         `json:"name"`
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	Email                 string         `json:"email"`
	ClientAdmin           bool           `json:"client_admin"`
	OrgDefaultRole        string         `json:"org_default_role"`
	OrgRoles              map[string]any `json:"org_roles"`
	SSOOnly               bool           `json:"sso_only"`
	MFAEnabled            bool           `json:"mfa_enabled"`
	LoginFailures         int64          `json:"login_failures"`
	LastLoginIP           string         `json:"last_login_ip"`
	LastLoginAt           int64          `json:"last_login_at"`
	LastLoginUA           string         `json:"last_login_ua"`
	LastActivityAt        int64          `json:"last_activity_at"`
	Actions               int64          `json:"actions"`
	LastActionAt          int64          `json:"last_action_at"`
	InviteTokenExpiration int64          `json:"invite_token_expiration"`
	ResetTokenExpiration  int64          `json:"reset_token_expiration"`
	PasswordEnabledAt     int64          `json:"password_enabled_at"`
	TermsAgreedAsOf       int64          `json:"terms_agreed_as_of"`
	CreatedAt             int64          `json:"created_at"`
	UpdatedAt             int64          `json:"updated_at"`
}

type CreateUserRequest struct {
	Email          string            `json:"email"`
	FirstName      string            `json:"first_name,omitempty"`
	LastName       string            `json:"last_name,omitempty"`
	ClientAdmin    bool              `json:"client_admin,omitempty"`
	OrgDefaultRole string            `json:"org_default_role,omitempty"`
	OrgRoles       map[string]string `json:"org_roles,omitempty"`
}

type UpdateUserRequest struct {
	FirstName      string            `json:"first_name,omitempty"`
	LastName       string            `json:"last_name,omitempty"`
	ClientAdmin    *bool             `json:"client_admin,omitempty"`
	OrgDefaultRole string            `json:"org_default_role,omitempty"`
	OrgRoles       map[string]string `json:"org_roles,omitempty"`
}

func (c *Client) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, "/account/users", AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/account/users/%s", id), AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

func (c *Client) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/account/users/%s", id), AuthAccount, req)
	if err != nil {
		return nil, err
	}
	var user User
	if err := HandleResponse(resp, &user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &user, nil
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/account/users/%s", id), AuthAccount, nil)
	if err != nil {
		return err
	}
	return HandleResponse(resp, nil)
}

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/users", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var users []User
	if err := HandleResponse(resp, &users); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}
