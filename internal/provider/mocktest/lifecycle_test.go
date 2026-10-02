package mocktest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/petrutb/terraform-provider-runzero/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"runzero": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// clearProviderEnv keeps developer credentials from leaking into unit tests.
func clearProviderEnv(t *testing.T) {
	t.Helper()
	for _, v := range []string{
		"RUNZERO_API_URL", "RUNZERO_ORG_TOKEN", "RUNZERO_ACCOUNT_TOKEN",
		"RUNZERO_CLIENT_ID", "RUNZERO_CLIENT_SECRET",
	} {
		t.Setenv(v, "")
	}
}

func checkFakeObject(api *fakeAPI, collection, resourceName, field, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in state", resourceName)
		}
		obj, ok := api.get(collection + "/" + rs.Primary.ID)
		if !ok {
			return fmt.Errorf("%s %s does not exist in the API", collection, rs.Primary.ID)
		}
		if got := fmt.Sprint(obj[field]); got != want {
			return fmt.Errorf("API %s: want %q, got %q", field, want, got)
		}
		return nil
	}
}

func checkFakeEmpty(api *fakeAPI, collection string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if n := api.count(collection); n != 0 {
			return fmt.Errorf("%d objects left in %s after destroy", n, collection)
		}
		return nil
	}
}

func TestUnitProvider_MissingCredentials(t *testing.T) {
	clearProviderEnv(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "runzero" {
  api_url = "http://127.0.0.1:1"
}

data "runzero_sites" "all" {}
`,
				ExpectError: regexp.MustCompile(`Missing RunZero Credentials`),
			},
		},
	})
}

func TestUnitProvider_APIErrorSurfaced(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: api.providerConfig() + `
data "runzero_site" "missing" {
  id = "does-not-exist"
}
`,
				ExpectError: regexp.MustCompile(`(?s)Error reading site.*404`),
			},
		},
	})
}

func TestUnitSiteResource(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)
	const name = "runzero_site.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkFakeEmpty(api, "/org/sites"),
		Steps: []resource.TestStep{
			{
				Config: api.providerConfig() + `
resource "runzero_site" "test" {
  name        = "unit-site"
  description = "Unit test site"
  scope       = "192.168.1.0/24"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttr(name, "name", "unit-site"),
					resource.TestCheckResourceAttr(name, "description", "Unit test site"),
					resource.TestCheckResourceAttr(name, "scope", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(name, "excludes", ""),
					resource.TestCheckResourceAttrSet(name, "created_at"),
					checkFakeObject(api, "/org/sites", name, "scope", "192.168.1.0/24"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: api.providerConfig() + `
resource "runzero_site" "test" {
  name        = "unit-site-renamed"
  description = "Unit test site"
  scope       = "10.0.0.0/8"
  excludes    = "10.0.0.1"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(name, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "unit-site-renamed"),
					resource.TestCheckResourceAttr(name, "scope", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(name, "excludes", "10.0.0.1"),
					checkFakeObject(api, "/org/sites", name, "name", "unit-site-renamed"),
				),
			},
		},
	})
}

func TestUnitSiteResource_RecreatedWhenDeletedOutsideTerraform(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)
	const name = "runzero_site.test"
	config := api.providerConfig() + `
resource "runzero_site" "test" {
  name = "drift-site"
}
`
	var firstID string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					firstID = s.RootModule().Resources[name].Primary.ID
					return nil
				},
			},
			{
				PreConfig: func() { api.remove("/org/sites/" + firstID) },
				Config:    config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(name, plancheck.ResourceActionCreate),
					},
				},
				Check: func(s *terraform.State) error {
					if s.RootModule().Resources[name].Primary.ID == firstID {
						return fmt.Errorf("expected a new site ID after recreation")
					}
					return nil
				},
			},
		},
	})
}

func TestUnitSiteDataSources(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: api.providerConfig() + `
resource "runzero_site" "a" {
  name  = "site-a"
  scope = "10.1.0.0/16"
}

resource "runzero_site" "b" {
  name = "site-b"
}

data "runzero_site" "a" {
  id = runzero_site.a.id
}

data "runzero_sites" "all" {
  depends_on = [runzero_site.a, runzero_site.b]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.runzero_site.a", "id", "runzero_site.a", "id"),
					resource.TestCheckResourceAttr("data.runzero_site.a", "name", "site-a"),
					resource.TestCheckResourceAttr("data.runzero_site.a", "scope", "10.1.0.0/16"),
					resource.TestCheckResourceAttr("data.runzero_sites.all", "sites.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("data.runzero_sites.all", "sites.*", map[string]string{"name": "site-b"}),
				),
			},
		},
	})
}

func TestUnitGroupResource(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)
	const name = "runzero_group.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkFakeEmpty(api, "/account/groups"),
		Steps: []resource.TestStep{
			{
				Config: api.providerConfig() + `
resource "runzero_group" "test" {
  name             = "unit-group"
  description      = "Unit test group"
  org_default_role = "viewer"
  org_roles = {
    "org-1" = "admin"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttr(name, "name", "unit-group"),
					resource.TestCheckResourceAttr(name, "org_default_role", "viewer"),
					resource.TestCheckResourceAttr(name, "org_roles.org-1", "admin"),
					checkFakeObject(api, "/account/groups", name, "org_default_role", "viewer"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: api.providerConfig() + `
resource "runzero_group" "test" {
  name             = "unit-group-renamed"
  description      = "Unit test group"
  org_default_role = "user"
  org_roles = {
    "org-1" = "admin"
  }
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(name, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "unit-group-renamed"),
					resource.TestCheckResourceAttr(name, "org_default_role", "user"),
					checkFakeObject(api, "/account/groups", name, "name", "unit-group-renamed"),
				),
			},
		},
	})
}

func TestUnitScanResource(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)
	const name = "runzero_scan.test"
	scanConfig := func(targets string) string {
		return api.providerConfig() + fmt.Sprintf(`
resource "runzero_site" "test" {
  name = "scan-site"
}

resource "runzero_scan" "test" {
  site_id        = runzero_site.test.id
  name           = "unit-scan"
  targets        = %q
  scan_frequency = "weekly"
  subnet_ping    = "true"
}
`, targets)
	}
	var firstID string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if !api.received("POST /org/tasks/" + firstID + "/stop") {
				return fmt.Errorf("scan task %s was not stopped on replace or destroy", firstID)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: scanConfig("10.0.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttr(name, "status", "new"),
					resource.TestCheckResourceAttr(name, "targets", "10.0.0.0/24"),
					func(s *terraform.State) error {
						firstID = s.RootModule().Resources[name].Primary.ID
						return nil
					},
				),
			},
			{
				Config: scanConfig("10.0.1.0/24"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(name, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.TestCheckResourceAttr(name, "targets", "10.0.1.0/24"),
			},
		},
	})
}

func TestUnitOrgAPIKeyResource(t *testing.T) {
	clearProviderEnv(t)
	api := newFakeAPI(t)
	api.setOrgKey(map[string]any{
		"id":              "key-1",
		"organization_id": "org-1",
		"token":           "OT_original",
		"type":            "org",
		"counter":         5,
	})
	const name = "runzero_org_api_key.test"
	keyConfig := func(rotate string) string {
		return api.providerConfig() + fmt.Sprintf(`
resource "runzero_org_api_key" "test" {
  rotate = %q
}
`, rotate)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			if !api.received("DELETE /org/key") {
				return fmt.Errorf("org key was not deleted on destroy")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: keyConfig("v1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", "key-1"),
					resource.TestCheckResourceAttr(name, "token", "OT_original"),
					resource.TestCheckResourceAttr(name, "counter", "5"),
				),
			},
			{
				Config: keyConfig("v2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(name, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", "key-1"),
					resource.TestMatchResourceAttr(name, "token", regexp.MustCompile(`^OT_rotated-`)),
					resource.TestCheckResourceAttr(name, "counter", "0"),
				),
			},
		},
	})
}
