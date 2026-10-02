package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(Config{
		BaseURL:  "https://example.com/api/v1.0",
		OrgToken: "OT_test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.baseURL != "https://example.com/api/v1.0" {
		t.Errorf("expected base URL https://example.com/api/v1.0, got %s", c.baseURL)
	}
	if c.orgToken != "OT_test" {
		t.Errorf("expected org token OT_test, got %s", c.orgToken)
	}
}

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient(Config{
		BaseURL: "://bad-url",
	})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestHandleResponse_Success(t *testing.T) {
	body := `{"id":"abc","name":"test"}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	var result struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err := HandleResponse(resp, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "abc" {
		t.Errorf("expected id abc, got %s", result.ID)
	}
}

func TestHandleResponse_APIError(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader(`{"error":"not found"}`)),
	}

	err := HandleResponse(resp, nil)
	if err == nil {
		t.Fatal("expected error for 404")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

func TestDoRequest_Auth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer OT_test123" {
			t.Errorf("expected Bearer OT_test123, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c, _ := NewClient(Config{
		BaseURL:  srv.URL,
		OrgToken: "OT_test123",
	})

	resp, err := c.DoRequest(context.Background(), http.MethodGet, "/test", AuthOrg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestDoRequest_AccountAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer CT_acct_token" {
			t.Errorf("expected Bearer CT_acct_token, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c, _ := NewClient(Config{
		BaseURL:      srv.URL,
		AccountToken: "CT_acct_token",
	})

	resp, err := c.DoRequest(context.Background(), http.MethodGet, "/test", AuthAccount, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestDoRequest_RateLimitRetry(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c, _ := NewClient(Config{
		BaseURL:  srv.URL,
		OrgToken: "OT_test",
	})

	resp, err := c.DoRequest(context.Background(), http.MethodGet, "/test", AuthOrg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestCreateSite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/org/sites" {
			t.Errorf("expected /org/sites, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Site{
			ID:   "site-123",
			Name: "Test Site",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: "OT_test"})
	site, err := c.CreateSite(context.Background(), CreateSiteRequest{Name: "Test Site"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site.ID != "site-123" {
		t.Errorf("expected site ID site-123, got %s", site.ID)
	}
}

func TestGetSite_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: "OT_test"})
	_, err := c.GetSite(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for 404")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
}

func TestListSites(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]Site{
			{ID: "s1", Name: "Site 1"},
			{ID: "s2", Name: "Site 2"},
		})
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: "OT_test"})
	sites, err := c.ListSites(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sites) != 2 {
		t.Errorf("expected 2 sites, got %d", len(sites))
	}
}

func TestOAuthFallback_NoCredentials(t *testing.T) {
	c, _ := NewClient(Config{BaseURL: "https://example.com"})
	_, err := c.tokenForAuth(context.Background(), AuthOrg)
	if err == nil {
		t.Fatal("expected error when no credentials configured")
	}
}
