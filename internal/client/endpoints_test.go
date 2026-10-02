package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testOrgToken     = "OT_org"
	testAccountToken = "CT_acct"
)

type endpointCase struct {
	name     string
	method   string
	path     string
	query    string
	token    string
	response string
	call     func(ctx context.Context, c *Client) error
}

func ignore[T any](_ T, err error) error { return err }

func TestEndpoints(t *testing.T) {
	const obj = `{"id":"x1"}`
	const list = `[{"id":"x1"}]`

	cases := []endpointCase{
		// Account
		{"ListAccountSites", "GET", "/account/sites", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAccountSites(ctx, "")) }},
		{"ListAccountSites search", "GET", "/account/sites", "search=name%3Aprod+site", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAccountSites(ctx, "name:prod site")) }},
		{"ListAccountTasks", "GET", "/account/tasks", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAccountTasks(ctx, "")) }},
		{"InviteUser", "PUT", "/account/users/invite", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.InviteUser(ctx, UserInviteRequest{})) }},
		{"ResetUserMFA", "PATCH", "/account/users/u1/resetMFA", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.ResetUserMFA(ctx, "u1")) }},
		{"ResetUserLockout", "PATCH", "/account/users/u1/resetLockout", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.ResetUserLockout(ctx, "u1")) }},
		{"ResetUserPassword", "PATCH", "/account/users/u1/resetPassword", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.ResetUserPassword(ctx, "u1")) }},

		// Agents
		{"ListAgents", "GET", "/org/agents", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAgents(ctx)) }},
		{"GetAgent", "GET", "/org/agents/a1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetAgent(ctx, "a1")) }},
		{"UpdateAgent", "PATCH", "/org/agents/a1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateAgent(ctx, "a1", UpdateAgentRequest{}))
			}},
		{"DeleteAgent", "DELETE", "/org/agents/a1", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteAgent(ctx, "a1") }},
		{"UpgradeAgent", "POST", "/org/agents/a1/update", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.UpgradeAgent(ctx, "a1") }},

		// API keys
		{"CreateAPIKey", "PUT", "/account/keys", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.CreateAPIKey(ctx, CreateAPIKeyRequest{})) }},
		{"GetAPIKey", "GET", "/account/keys/k1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetAPIKey(ctx, "k1")) }},
		{"DeleteAPIKey", "DELETE", "/account/keys/k1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteAPIKey(ctx, "k1") }},
		{"ListAPIKeys", "GET", "/account/keys", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAPIKeys(ctx)) }},
		{"RotateAPIKey", "PATCH", "/account/keys/k1/rotate", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.RotateAPIKey(ctx, "k1")) }},

		// Asset ownership types
		{"UpdateAssetOwnershipType", "PATCH", "/account/assets/ownership-types/o1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateAssetOwnershipType(ctx, "o1", UpdateAssetOwnershipTypeRequest{}))
			}},
		{"DeleteAssetOwnershipType", "DELETE", "/account/assets/ownership-types/o1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteAssetOwnershipType(ctx, "o1") }},
		{"ListAssetOwnershipTypes", "GET", "/account/assets/ownership-types", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAssetOwnershipTypes(ctx)) }},

		// Assets
		{"ListAssets", "GET", "/org/assets", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAssets(ctx, "")) }},
		{"ListAssets search", "GET", "/org/assets", "search=alive%3At", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListAssets(ctx, "alive:t")) }},
		{"GetAsset", "GET", "/org/assets/as1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetAsset(ctx, "as1")) }},
		{"DeleteAsset", "DELETE", "/org/assets/as1", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteAsset(ctx, "as1") }},
		{"UpdateAssetComments", "PATCH", "/org/assets/as1/comments", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateAssetComments(ctx, "as1", UpdateAssetCommentsRequest{}))
			}},
		{"UpdateAssetTags", "PATCH", "/org/assets/as1/tags", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateAssetTags(ctx, "as1", UpdateAssetTagsRequest{}))
			}},
		{"UpdateAssetCriticality", "PATCH", "/org/assets/as1/criticality", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateAssetCriticality(ctx, "as1", UpdateAssetCriticalityRequest{}))
			}},
		{"UpdateAssetOwnerships", "PATCH", "/org/assets/as1/owners", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateAssetOwnerships(ctx, "as1", UpdateAssetOwnershipsRequest{}))
			}},

		// Credentials
		{"CreateCredential", "PUT", "/account/credentials", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateCredential(ctx, CreateCredentialRequest{}))
			}},
		{"GetCredential", "GET", "/account/credentials/c1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetCredential(ctx, "c1")) }},
		{"DeleteCredential", "DELETE", "/account/credentials/c1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteCredential(ctx, "c1") }},
		{"ListCredentials", "GET", "/account/credentials", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListCredentials(ctx)) }},

		// Custom integrations
		{"CreateCustomIntegration", "POST", "/account/custom-integrations", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateCustomIntegration(ctx, CreateCustomIntegrationRequest{}))
			}},
		{"GetCustomIntegration", "GET", "/account/custom-integrations/ci1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetCustomIntegration(ctx, "ci1")) }},
		{"UpdateCustomIntegration", "PATCH", "/account/custom-integrations/ci1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateCustomIntegration(ctx, "ci1", UpdateCustomIntegrationRequest{}))
			}},
		{"DeleteCustomIntegration", "DELETE", "/account/custom-integrations/ci1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteCustomIntegration(ctx, "ci1") }},
		{"ListCustomIntegrations", "GET", "/account/custom-integrations", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListCustomIntegrations(ctx)) }},

		// Explorers
		{"ListExplorers", "GET", "/org/explorers", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListExplorers(ctx)) }},
		{"GetExplorer", "GET", "/org/explorers/e1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetExplorer(ctx, "e1")) }},
		{"UpdateExplorer", "PATCH", "/org/explorers/e1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateExplorer(ctx, "e1", UpdateExplorerRequest{}))
			}},
		{"DeleteExplorer", "DELETE", "/org/explorers/e1", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteExplorer(ctx, "e1") }},
		{"UpgradeExplorer", "POST", "/org/explorers/e1/update", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.UpgradeExplorer(ctx, "e1") }},

		// Export tokens
		{"CreateExportToken", "POST", "/account/orgs/org1/exportTokens", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateExportToken(ctx, "org1", CreateExportTokenRequest{}))
			}},
		{"GetExportToken", "GET", "/account/orgs/org1/exportTokens/t1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetExportToken(ctx, "org1", "t1")) }},
		{"DeleteExportToken", "DELETE", "/account/orgs/org1/exportTokens/t1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteExportToken(ctx, "org1", "t1") }},
		{"ListExportTokens", "GET", "/account/orgs/org1/exportTokens", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListExportTokens(ctx, "org1")) }},
		{"RotateExportToken", "PATCH", "/account/orgs/org1/exportTokens/t1/rotate", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.RotateExportToken(ctx, "org1", "t1")) }},

		// Groups
		{"CreateGroup", "POST", "/account/groups", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.CreateGroup(ctx, CreateGroupRequest{})) }},
		{"GetGroup", "GET", "/account/groups/g1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetGroup(ctx, "g1")) }},
		{"UpdateGroup", "PUT", "/account/groups", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.UpdateGroup(ctx, UpdateGroupRequest{})) }},
		{"DeleteGroup", "DELETE", "/account/groups/g1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteGroup(ctx, "g1") }},
		{"ListGroups", "GET", "/account/groups", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListGroups(ctx)) }},

		// Health, hosted zones, license
		{"GetHealth", "GET", "/health", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetHealth(ctx)) }},
		{"GetHostedZone", "GET", "/org/hosted-zones/h1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetHostedZone(ctx, "h1")) }},
		{"ListHostedZones", "GET", "/org/hosted-zones", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListHostedZones(ctx)) }},
		{"GetLicense", "GET", "/account/license", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetLicense(ctx)) }},

		// Org key
		{"GetOrgKey", "GET", "/org/key", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetOrgKey(ctx)) }},
		{"RotateOrgKey", "PATCH", "/org/key/rotate", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.RotateOrgKey(ctx)) }},
		{"DeleteOrgKey", "DELETE", "/org/key", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteOrgKey(ctx) }},

		// Organizations
		{"CreateOrganization", "PUT", "/account/orgs", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateOrganization(ctx, CreateOrganizationRequest{}))
			}},
		{"GetOrganization", "GET", "/account/orgs/org1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetOrganization(ctx, "org1")) }},
		{"UpdateOrganization", "PATCH", "/account/orgs/org1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateOrganization(ctx, "org1", UpdateOrganizationRequest{}))
			}},
		{"DeleteOrganization", "DELETE", "/account/orgs/org1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteOrganization(ctx, "org1") }},
		{"ListOrganizations", "GET", "/account/orgs", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListOrganizations(ctx)) }},
		{"GetOrg", "GET", "/org", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetOrg(ctx)) }},
		{"UpdateOrg", "PATCH", "/org", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateOrg(ctx, UpdateOrganizationRequest{}))
			}},

		// Scan templates
		{"CreateScanTemplate", "POST", "/account/tasks/templates", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateScanTemplate(ctx, CreateScanTemplateRequest{}))
			}},
		{"GetScanTemplate", "GET", "/account/tasks/templates/st1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetScanTemplate(ctx, "st1")) }},
		{"UpdateScanTemplate", "PUT", "/account/tasks/templates", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateScanTemplate(ctx, UpdateScanTemplateRequest{}))
			}},
		{"DeleteScanTemplate", "DELETE", "/account/tasks/templates/st1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteScanTemplate(ctx, "st1") }},
		{"ListScanTemplates", "GET", "/account/tasks/templates", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListScanTemplates(ctx)) }},

		// Scans
		{"CreateScan", "PUT", "/org/sites/s1/scan", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateScan(ctx, "s1", CreateScanRequest{}))
			}},

		// Services
		{"ListServices", "GET", "/org/services", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListServices(ctx, "")) }},
		{"ListServices search", "GET", "/org/services", "search=port%3A22", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListServices(ctx, "port:22")) }},
		{"GetService", "GET", "/org/services/sv1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetService(ctx, "sv1")) }},
		{"DeleteService", "DELETE", "/org/services/sv1", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteService(ctx, "sv1") }},

		// Sites
		{"CreateSite", "PUT", "/org/sites", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.CreateSite(ctx, CreateSiteRequest{})) }},
		{"GetSite", "GET", "/org/sites/s1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetSite(ctx, "s1")) }},
		{"UpdateSite", "PATCH", "/org/sites/s1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateSite(ctx, "s1", UpdateSiteRequest{}))
			}},
		{"DeleteSite", "DELETE", "/org/sites/s1", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteSite(ctx, "s1") }},
		{"ListSites", "GET", "/org/sites", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListSites(ctx)) }},

		// SSO group mappings
		{"CreateSSOGroupMapping", "POST", "/account/sso/groups", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.CreateSSOGroupMapping(ctx, CreateSSOGroupMappingRequest{}))
			}},
		{"GetSSOGroupMapping", "GET", "/account/sso/groups/m1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetSSOGroupMapping(ctx, "m1")) }},
		{"UpdateSSOGroupMapping", "PUT", "/account/sso/groups", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateSSOGroupMapping(ctx, UpdateSSOGroupMappingRequest{}))
			}},
		{"DeleteSSOGroupMapping", "DELETE", "/account/sso/groups/m1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteSSOGroupMapping(ctx, "m1") }},
		{"ListSSOGroupMappings", "GET", "/account/sso/groups", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListSSOGroupMappings(ctx)) }},

		// Tasks
		{"ListTasks", "GET", "/org/tasks", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListTasks(ctx, "")) }},
		{"ListTasks status", "GET", "/org/tasks", "status=processed", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListTasks(ctx, "processed")) }},
		{"GetTask", "GET", "/org/tasks/t1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetTask(ctx, "t1")) }},
		{"UpdateTask", "PATCH", "/org/tasks/t1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateTask(ctx, "t1", UpdateTaskRequest{}))
			}},
		{"StopTask", "POST", "/org/tasks/t1/stop", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.StopTask(ctx, "t1") }},
		{"HideTask", "POST", "/org/tasks/t1/hide", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.HideTask(ctx, "t1") }},

		// Users
		{"CreateUser", "PUT", "/account/users", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.CreateUser(ctx, CreateUserRequest{})) }},
		{"GetUser", "GET", "/account/users/u1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetUser(ctx, "u1")) }},
		{"UpdateUser", "PATCH", "/account/users/u1", "", testAccountToken, obj,
			func(ctx context.Context, c *Client) error {
				return ignore(c.UpdateUser(ctx, "u1", UpdateUserRequest{}))
			}},
		{"DeleteUser", "DELETE", "/account/users/u1", "", testAccountToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteUser(ctx, "u1") }},
		{"ListUsers", "GET", "/account/users", "", testAccountToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListUsers(ctx)) }},

		// Wireless
		{"ListWirelessLANs", "GET", "/org/wireless", "", testOrgToken, list,
			func(ctx context.Context, c *Client) error { return ignore(c.ListWirelessLANs(ctx)) }},
		{"GetWirelessLAN", "GET", "/org/wireless/w1", "", testOrgToken, obj,
			func(ctx context.Context, c *Client) error { return ignore(c.GetWirelessLAN(ctx, "w1")) }},
		{"DeleteWirelessLAN", "DELETE", "/org/wireless/w1", "", testOrgToken, "",
			func(ctx context.Context, c *Client) error { return c.DeleteWirelessLAN(ctx, "w1") }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				if r.Method != tc.method {
					t.Errorf("method: want %s, got %s", tc.method, r.Method)
				}
				if r.URL.Path != tc.path {
					t.Errorf("path: want %s, got %s", tc.path, r.URL.Path)
				}
				if r.URL.RawQuery != tc.query {
					t.Errorf("query: want %q, got %q", tc.query, r.URL.RawQuery)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer "+tc.token {
					t.Errorf("auth: want Bearer %s, got %s", tc.token, got)
				}
				_, _ = io.WriteString(w, tc.response)
			}))
			defer srv.Close()

			c, err := NewClient(Config{BaseURL: srv.URL, OrgToken: testOrgToken, AccountToken: testAccountToken})
			if err != nil {
				t.Fatal(err)
			}
			if err := tc.call(context.Background(), c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if calls != 1 {
				t.Errorf("expected 1 request, got %d", calls)
			}
		})
	}
}

func TestEndpoints_APIErrorPropagates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"error":"forbidden"}`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: testOrgToken, AccountToken: testAccountToken})
	_, err := c.GetGroup(context.Background(), "g1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "forbidden") {
		t.Errorf("error should include status and body, got: %v", err)
	}
}

func TestCreateAssetOwnershipType_SendsArrayAndPicksByName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []CreateAssetOwnershipTypeRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("request body is not a JSON array: %v", err)
			return
		}
		if len(body) != 1 || body[0].Name != "Owner" {
			t.Errorf("unexpected body: %+v", body)
		}
		_, _ = io.WriteString(w, `[{"id":"other","name":"Other"},{"id":"o1","name":"Owner"}]`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, AccountToken: testAccountToken})
	got, err := c.CreateAssetOwnershipType(context.Background(), CreateAssetOwnershipTypeRequest{Name: "Owner"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "o1" {
		t.Errorf("want id o1, got %s", got.ID)
	}
}

func TestCreateAssetOwnershipType_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `[]`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, AccountToken: testAccountToken})
	if _, err := c.CreateAssetOwnershipType(context.Background(), CreateAssetOwnershipTypeRequest{Name: "x"}); err == nil {
		t.Fatal("expected error for empty response")
	}
}

func TestGetAssetOwnershipType_FiltersListByID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/account/assets/ownership-types" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `[{"id":"a","name":"A"},{"id":"b","name":"B"}]`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, AccountToken: testAccountToken})
	got, err := c.GetAssetOwnershipType(context.Background(), "b")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "B" {
		t.Errorf("want B, got %s", got.Name)
	}

	_, err = c.GetAssetOwnershipType(context.Background(), "missing")
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404 APIError, got %v", err)
	}
}

func TestCreateScan_SendsHyphenatedOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
			return
		}
		want := map[string]string{
			"targets":        "10.0.0.0/24",
			"scan-name":      "weekly",
			"scan-frequency": "weekly",
			"subnet-ping":    "true",
		}
		for k, v := range want {
			if body[k] != v {
				t.Errorf("%s: want %q, got %v", k, v, body[k])
			}
		}
		if _, ok := body["tos"]; ok {
			t.Error("empty optional fields should be omitted")
		}
		_, _ = io.WriteString(w, `{"id":"task1","site_id":"s1"}`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: testOrgToken})
	scan, err := c.CreateScan(context.Background(), "s1", CreateScanRequest{
		Targets:       "10.0.0.0/24",
		ScanName:      "weekly",
		ScanFrequency: "weekly",
		SubnetPing:    "true",
	})
	if err != nil {
		t.Fatal(err)
	}
	if scan.ID != "task1" {
		t.Errorf("want task1, got %s", scan.ID)
	}
}

func TestOAuth_TokenFetchedOnceAndCached(t *testing.T) {
	var tokenCalls, apiCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/account/api/token":
			atomic.AddInt32(&tokenCalls, 1)
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse form: %v", err)
				return
			}
			if r.Form.Get("grant_type") != "client_credentials" ||
				r.Form.Get("client_id") != "id" || r.Form.Get("client_secret") != "secret" {
				t.Errorf("unexpected token form: %v", r.Form)
			}
			_, _ = io.WriteString(w, `{"access_token":"bearer-xyz","expires_in":3600}`)
		default:
			atomic.AddInt32(&apiCalls, 1)
			if got := r.Header.Get("Authorization"); got != "Bearer bearer-xyz" {
				t.Errorf("want OAuth bearer, got %s", got)
			}
			_, _ = io.WriteString(w, `[]`)
		}
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, ClientID: "id", ClientSecret: "secret"})
	ctx := context.Background()
	if _, err := c.ListSites(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListUsers(ctx); err != nil {
		t.Fatal(err)
	}
	if tokenCalls != 1 {
		t.Errorf("want 1 token request, got %d", tokenCalls)
	}
	if apiCalls != 2 {
		t.Errorf("want 2 API requests, got %d", apiCalls)
	}
}

func TestOAuth_TokenEndpointError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, "bad client")
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, ClientID: "id", ClientSecret: "wrong"})
	_, err := c.ListSites(context.Background())
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("want 401 token error, got %v", err)
	}
}

func TestDoRequest_RateLimitHonorsContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: testOrgToken})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := c.DoRequest(ctx, http.MethodGet, "/org/sites", AuthOrg, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want context deadline error, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("cancellation was not honored while waiting, took %s", elapsed)
	}
}

func TestDoRequest_RetryResendsBody(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"name":"retry-site"}` {
			t.Errorf("attempt %d: unexpected body %q", atomic.LoadInt32(&attempts)+1, body)
		}
		if atomic.AddInt32(&attempts, 1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = io.WriteString(w, `{"id":"s1","name":"retry-site"}`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: testOrgToken})
	site, err := c.CreateSite(context.Background(), CreateSiteRequest{Name: "retry-site"})
	if err != nil {
		t.Fatal(err)
	}
	if site.ID != "s1" || attempts != 2 {
		t.Errorf("want site s1 after 2 attempts, got %q after %d", site.ID, attempts)
	}
}

func TestDoRequest_RetriesExhausted(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{BaseURL: srv.URL, OrgToken: testOrgToken})
	_, err := c.DoRequest(context.Background(), http.MethodGet, "/org/sites", AuthOrg, nil)
	if err == nil || !strings.Contains(err.Error(), "rate limited") {
		t.Fatalf("want rate limit error, got %v", err)
	}
	if attempts != maxRetries+1 {
		t.Errorf("want %d attempts, got %d", maxRetries+1, attempts)
	}
}

func TestRetryDelay(t *testing.T) {
	future := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)
	cases := []struct {
		name       string
		retryAfter string
		attempt    int
		min, max   time.Duration
	}{
		{"seconds", "7", 0, 7 * time.Second, 7 * time.Second},
		{"zero", "0", 2, 0, 0},
		{"capped", "86400", 0, maxRetryWait, maxRetryWait},
		{"http date", future, 0, 25 * time.Second, 30 * time.Second},
		{"past date", "Mon, 02 Jan 2006 15:04:05 GMT", 0, 0, 0},
		{"missing, first attempt", "", 0, retryBaseDelay, retryBaseDelay},
		{"missing, third attempt", "", 2, 3 * retryBaseDelay, 3 * retryBaseDelay},
		{"garbage", "soon", 1, 2 * retryBaseDelay, 2 * retryBaseDelay},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := retryDelay(tc.retryAfter, tc.attempt)
			if got < tc.min || got > tc.max {
				t.Errorf("retryDelay(%q, %d) = %s, want between %s and %s", tc.retryAfter, tc.attempt, got, tc.min, tc.max)
			}
		})
	}
}

func TestUserAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "terraform-provider-runzero/1.2.3" {
			t.Errorf("unexpected User-Agent %q", got)
		}
		if r.URL.Path == "/account/api/token" {
			_, _ = io.WriteString(w, `{"access_token":"b","expires_in":3600}`)
			return
		}
		_, _ = io.WriteString(w, `[]`)
	}))
	defer srv.Close()

	c, _ := NewClient(Config{
		BaseURL:      srv.URL,
		ClientID:     "id",
		ClientSecret: "secret",
		UserAgent:    "terraform-provider-runzero/1.2.3",
	})
	if _, err := c.ListSites(context.Background()); err != nil {
		t.Fatal(err)
	}
}
