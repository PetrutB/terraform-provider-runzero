package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

func TestAccSite(t *testing.T) {
	name := randomName()
	const res = "runzero_site.test"
	config := func(siteName, scope string) string {
		return fmt.Sprintf(`
resource "runzero_site" "test" {
  name        = %q
  description = "Terraform acceptance test"
  scope       = %q
  excludes    = "192.0.2.1"
}

data "runzero_site" "test" {
  id = runzero_site.test.id
}

data "runzero_sites" "all" {
  depends_on = [runzero_site.test]
}

data "runzero_account_sites" "all" {
  depends_on = [runzero_site.test]
}
`, siteName, scope)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOrg(t)
			testAccPreCheckAccount(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_site", byID((*client.Client).GetSite)),
		Steps: []resource.TestStep{
			{
				Config: config(name, "192.0.2.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "name", name),
					resource.TestCheckResourceAttr(res, "scope", "192.0.2.0/24"),
					resource.TestCheckResourceAttr(res, "excludes", "192.0.2.1"),
					resource.TestCheckResourceAttrPair("data.runzero_site.test", "name", res, "name"),
					resource.TestCheckResourceAttrPair("data.runzero_site.test", "scope", res, "scope"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_sites.all", "sites.*.id", res, "id"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_account_sites.all", "sites.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config(name+"-updated", "198.51.100.0/24"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "name", name+"-updated"),
					resource.TestCheckResourceAttr(res, "scope", "198.51.100.0/24"),
				),
			},
		},
	})
}

func TestAccScan(t *testing.T) {
	explorerID := requireEnv(t, "RUNZERO_TEST_EXPLORER_ID")
	name := randomName()
	const res = "runzero_scan.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_site", byID((*client.Client).GetSite)),
		Steps: []resource.TestStep{
			{
				// TEST-NET-1 (RFC 5737) is reserved for documentation, so nothing real is probed.
				Config: fmt.Sprintf(`
resource "runzero_site" "test" {
  name = %[1]q
}

resource "runzero_scan" "test" {
  site_id     = runzero_site.test.id
  name        = %[1]q
  description = "Terraform acceptance test"
  targets     = "192.0.2.0/30"
  explorer    = %[2]q
  rate        = "100"
  tcp_ports   = "22,443"
}

data "runzero_task" "test" {
  id = runzero_scan.test.id
}

data "runzero_tasks" "all" {
  depends_on = [runzero_scan.test]
}

data "runzero_account_tasks" "all" {
  depends_on = [runzero_scan.test]
}
`, name, explorerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttrSet(res, "status"),
					resource.TestCheckResourceAttrPair(res, "site_id", "runzero_site.test", "id"),
					resource.TestCheckResourceAttrPair("data.runzero_task.test", "id", res, "id"),
					resource.TestCheckResourceAttr("data.runzero_task.test", "name", name),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_tasks.all", "tasks.*.id", res, "id"),
					resource.TestCheckResourceAttrSet("data.runzero_account_tasks.all", "tasks.#"),
				),
			},
		},
	})
}

func TestAccOrgAPIKey(t *testing.T) {
	const res = "runzero_org_api_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks:   requireRemovedBlock,
		Steps: []resource.TestStep{
			{
				// No rotate value: rotating would invalidate the org token used by the tests.
				Config: `
resource "runzero_org_api_key" "test" {}

data "runzero_org_api_key" "current" {
  depends_on = [runzero_org_api_key.test]
}

data "runzero_org" "current" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttrSet(res, "token"),
					resource.TestCheckResourceAttrPair("data.runzero_org_api_key.current", "id", res, "id"),
					resource.TestCheckResourceAttrPair("data.runzero_org.current", "id", res, "organization_id"),
					resource.TestCheckResourceAttrSet("data.runzero_org.current", "name"),
				),
			},
			{
				// Release the key from state; destroying it would revoke the test credentials.
				Config: removedConfig(res),
			},
		},
	})
}
