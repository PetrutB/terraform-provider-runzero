package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	maxRetries = 3
	// maxRetryWait caps how long a single Retry-After value can stall a request.
	maxRetryWait = 2 * time.Minute
)

// retryBaseDelay is the backoff unit used when the API does not send Retry-After.
var retryBaseDelay = 5 * time.Second

type AuthType int

const (
	AuthOrg AuthType = iota
	AuthAccount
	AuthExport
)

type Config struct {
	BaseURL      string
	OrgToken     string
	AccountToken string
	ClientID     string
	ClientSecret string
	UserAgent    string
}

type Client struct {
	httpClient   *http.Client
	baseURL      string
	orgToken     string
	accountToken string
	clientID     string
	clientSecret string
	userAgent    string

	mu          sync.Mutex
	bearerToken string
	tokenExpiry time.Time
}

func NewClient(cfg Config) (*Client, error) {
	parsedURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:      parsedURL.String(),
		orgToken:     cfg.OrgToken,
		accountToken: cfg.AccountToken,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		userAgent:    cfg.UserAgent,
	}, nil
}

func (c *Client) tokenForAuth(ctx context.Context, authType AuthType) (string, error) {
	switch authType {
	case AuthOrg:
		if c.orgToken != "" {
			return c.orgToken, nil
		}
		return c.getOAuthToken(ctx)
	case AuthAccount:
		if c.accountToken != "" {
			return c.accountToken, nil
		}
		return c.getOAuthToken(ctx)
	case AuthExport:
		if c.orgToken != "" {
			return c.orgToken, nil
		}
		return c.getOAuthToken(ctx)
	default:
		return c.getOAuthToken(ctx)
	}
}

func (c *Client) getOAuthToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bearerToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.bearerToken, nil
	}

	if c.clientID == "" || c.clientSecret == "" {
		return "", fmt.Errorf("no valid credentials: provide org_token, account_token, or client_id+client_secret")
	}

	tokenURL := fmt.Sprintf("%s/account/api/token", c.baseURL)
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("OAuth2 token request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OAuth2 token request returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	c.bearerToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)
	return c.bearerToken, nil
}

func (c *Client) DoRequest(ctx context.Context, method, path string, authType AuthType, body interface{}) (*http.Response, error) {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	token, err := c.tokenForAuth(ctx, authType)
	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s%s", c.baseURL, path)
	for attempt := 0; ; attempt++ {
		// A request body can only be read once, so each attempt needs a fresh request.
		var bodyReader io.Reader
		if payload != nil {
			bodyReader = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if c.userAgent != "" {
			req.Header.Set("User-Agent", c.userAgent)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		wait := retryDelay(resp.Header.Get("Retry-After"), attempt)
		_ = resp.Body.Close()
		if attempt >= maxRetries {
			return nil, fmt.Errorf("rate limited after %d retries", maxRetries)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}
}

// retryDelay honors a Retry-After header (seconds or HTTP date) and otherwise backs off linearly.
func retryDelay(retryAfter string, attempt int) time.Duration {
	if retryAfter != "" {
		if secs, err := strconv.Atoi(retryAfter); err == nil && secs >= 0 {
			return min(time.Duration(secs)*time.Second, maxRetryWait)
		}
		if at, err := http.ParseTime(retryAfter); err == nil {
			return min(max(time.Until(at), 0), maxRetryWait)
		}
	}
	return retryBaseDelay * time.Duration(attempt+1)
}

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("RunZero API error (HTTP %d): %s — %s", e.StatusCode, e.Message, e.Body)
}

func HandleResponse(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
			Body:       string(bodyBytes),
		}
	}

	if target != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}
