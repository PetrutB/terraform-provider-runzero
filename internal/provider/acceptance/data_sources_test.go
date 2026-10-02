package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Data sources for read-only objects. List data sources run without fixtures; single-item
// data sources need an existing object ID from the environment.

func TestAccDataSourceHealth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "runzero_health" "current" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.runzero_health.current", "status"),
					resource.TestCheckResourceAttrSet("data.runzero_health.current", "version"),
				),
			},
		},
	})
}

func TestAccDataSourceLicense(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "runzero_license" "current" {}`,
				Check:  resource.TestCheckResourceAttrSet("data.runzero_license.current", "id"),
			},
		},
	})
}

func TestAccDataSourceOrgLists(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "runzero_agents" "all" {}
data "runzero_explorers" "all" {}
data "runzero_assets" "all" {}
data "runzero_services" "all" {}
data "runzero_tasks" "all" {}
data "runzero_hosted_zones" "all" {}
data "runzero_wireless_lans" "all" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.runzero_agents.all", "agents.#"),
					resource.TestCheckResourceAttrSet("data.runzero_explorers.all", "explorers.#"),
					resource.TestCheckResourceAttrSet("data.runzero_assets.all", "assets.#"),
					resource.TestCheckResourceAttrSet("data.runzero_services.all", "services.#"),
					resource.TestCheckResourceAttrSet("data.runzero_tasks.all", "tasks.#"),
					resource.TestCheckResourceAttrSet("data.runzero_hosted_zones.all", "hosted_zones.#"),
					resource.TestCheckResourceAttrSet("data.runzero_wireless_lans.all", "wireless_lans.#"),
				),
			},
		},
	})
}

func TestAccDataSourceAccountLists(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "runzero_organizations" "all" {}
data "runzero_account_sites" "all" {}
data "runzero_account_tasks" "all" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.runzero_organizations.all", "organizations.#"),
					resource.TestCheckResourceAttrSet("data.runzero_account_sites.all", "sites.#"),
					resource.TestCheckResourceAttrSet("data.runzero_account_tasks.all", "tasks.#"),
				),
			},
		},
	})
}

func TestAccDataSourceService(t *testing.T) {
	id := requireEnv(t, "RUNZERO_TEST_SERVICE_ID")
	testAccSingleDataSource(t, "runzero_service", id, "service_asset_id")
}

func TestAccDataSourceTask(t *testing.T) {
	id := requireEnv(t, "RUNZERO_TEST_TASK_ID")
	testAccSingleDataSource(t, "runzero_task", id, "name")
}

func TestAccDataSourceHostedZone(t *testing.T) {
	id := requireEnv(t, "RUNZERO_TEST_HOSTED_ZONE_ID")
	testAccSingleDataSource(t, "runzero_hosted_zone", id, "name")
}

func TestAccDataSourceWirelessLAN(t *testing.T) {
	id := requireEnv(t, "RUNZERO_TEST_WIRELESS_LAN_ID")
	testAccSingleDataSource(t, "runzero_wireless_lan", id, "essid")
}

// testAccSingleDataSource reads an existing object by ID and checks that a key attribute is populated.
func testAccSingleDataSource(t *testing.T, typeName, id, attr string) {
	t.Helper()
	name := "data." + typeName + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data %q "test" {
  id = %q
}
`, typeName, id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", id),
					resource.TestCheckResourceAttrSet(name, attr),
				),
			},
		},
	})
}
