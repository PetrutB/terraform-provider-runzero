package acceptance

// Acceptance tests create, update, and delete real objects. Run them against a
// dedicated test account/organization:
//
//   TF_ACC=1 go test ./internal/provider/acceptance/ -v -timeout 60m
//
// Credentials (same variables the provider reads):
//
//   RUNZERO_API_URL        optional, defaults to https://console.runzero.com/api/v1.0
//   RUNZERO_ORG_TOKEN      required for organization-level tests (OT...)
//   RUNZERO_ACCOUNT_TOKEN  required for account-level tests (CT...)
//   RUNZERO_CLIENT_ID / RUNZERO_CLIENT_SECRET  may replace either token
//
// Optional fixtures. Tests that need them are skipped when unset:
//
//   RUNZERO_TEST_EMAIL_DOMAIN       domain for generated user emails (default example.com)
//   RUNZERO_TEST_INVITE_EMAIL       address that receives a real invitation email
//   RUNZERO_TEST_EXPLORER_ID        explorer that runs runzero_scan and runzero_explorer tests
//   RUNZERO_TEST_AGENT_ID           agent for runzero_agent tests
//   RUNZERO_TEST_ASSET_ID           asset for runzero_asset tests (comments/criticality are changed)
//   RUNZERO_TEST_SERVICE_ID         service for the runzero_service data source
//   RUNZERO_TEST_TASK_ID            task for the runzero_task data source
//   RUNZERO_TEST_HOSTED_ZONE_ID     hosted zone for the runzero_hosted_zone data source
//   RUNZERO_TEST_WIRELESS_LAN_ID    wireless LAN for the runzero_wireless_lan data source
//
// Adopted objects (agents, explorers, assets, the org API key) are released with a
// `removed` block at the end of each test, so they are never deleted.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
	"github.com/petrutb/terraform-provider-runzero/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"runzero": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// `removed` blocks need Terraform 1.7+.
var requireRemovedBlock = []tfversion.TerraformVersionCheck{
	tfversion.SkipBelow(tfversion.Version1_7_0),
}

func hasOAuth() bool {
	return os.Getenv("RUNZERO_CLIENT_ID") != "" && os.Getenv("RUNZERO_CLIENT_SECRET") != ""
}

// testAccPreCheckOrg ensures organization-scoped credentials are configured.
func testAccPreCheckOrg(t *testing.T) {
	t.Helper()
	if os.Getenv("RUNZERO_ORG_TOKEN") == "" && !hasOAuth() {
		t.Fatal("RUNZERO_ORG_TOKEN (or RUNZERO_CLIENT_ID/RUNZERO_CLIENT_SECRET) must be set for acceptance tests")
	}
}

// testAccPreCheckAccount ensures account-scoped credentials are configured.
func testAccPreCheckAccount(t *testing.T) {
	t.Helper()
	if os.Getenv("RUNZERO_ACCOUNT_TOKEN") == "" && !hasOAuth() {
		t.Fatal("RUNZERO_ACCOUNT_TOKEN (or RUNZERO_CLIENT_ID/RUNZERO_CLIENT_SECRET) must be set for acceptance tests")
	}
}

// requireEnv returns the value of an optional fixture variable or skips the test.
func requireEnv(t *testing.T, name string) string {
	t.Helper()
	v := os.Getenv(name)
	if v == "" {
		t.Skipf("%s is not set", name)
	}
	return v
}

func randomName() string {
	return acctest.RandomWithPrefix("tf-acc")
}

func randomEmail() string {
	domain := os.Getenv("RUNZERO_TEST_EMAIL_DOMAIN")
	if domain == "" {
		domain = "example.com"
	}
	return fmt.Sprintf("%s@%s", acctest.RandomWithPrefix("tf-acc"), domain)
}

// removedConfig stops managing a resource without destroying the real object.
func removedConfig(address string) string {
	return fmt.Sprintf(`
removed {
  from = %s

  lifecycle {
    destroy = false
  }
}
`, address)
}

// testClient builds an API client from the same environment variables as the provider.
func testClient() (*client.Client, error) {
	apiURL := os.Getenv("RUNZERO_API_URL")
	if apiURL == "" {
		apiURL = "https://console.runzero.com/api/v1.0"
	}
	return client.NewClient(client.Config{
		BaseURL:      apiURL,
		OrgToken:     os.Getenv("RUNZERO_ORG_TOKEN"),
		AccountToken: os.Getenv("RUNZERO_ACCOUNT_TOKEN"),
		ClientID:     os.Getenv("RUNZERO_CLIENT_ID"),
		ClientSecret: os.Getenv("RUNZERO_CLIENT_SECRET"),
	})
}

type getFunc func(ctx context.Context, c *client.Client, rs *terraform.ResourceState) error

// byID adapts a client getter such as (*client.Client).GetGroup to a getFunc.
func byID[T any](get func(*client.Client, context.Context, string) (T, error)) getFunc {
	return func(ctx context.Context, c *client.Client, rs *terraform.ResourceState) error {
		_, err := get(c, ctx, rs.Primary.ID)
		return err
	}
}

// checkDestroy verifies that every resource of the given type returns 404 after destroy.
func checkDestroy(resourceType string, get getFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c, err := testClient()
		if err != nil {
			return err
		}
		ctx := context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			err := get(ctx, c, rs)
			if err == nil {
				return fmt.Errorf("%s %s still exists", resourceType, rs.Primary.ID)
			}
			var apiErr *client.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
				continue
			}
			return fmt.Errorf("checking %s %s: %w", resourceType, rs.Primary.ID, err)
		}
		return nil
	}
}

// captureAttr stores an attribute value for comparison in a later step.
func captureAttr(name, key string, dst *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%s not found in state", name)
		}
		*dst = rs.Primary.Attributes[key]
		return nil
	}
}

// checkAttrChanged fails if an attribute still has the value captured earlier.
func checkAttrChanged(name, key string, previous *string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, key, func(v string) error {
		if v == *previous {
			return fmt.Errorf("%s.%s did not change", name, key)
		}
		return nil
	})
}
