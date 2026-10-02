package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Agents, explorers, and assets cannot be created through the API. These tests adopt an
// existing object by ID, verify import, and then release it from state with a `removed`
// block so the real object is not deleted.

var volatileAgentAttrs = []string{"connected", "last_checkin", "updated_at", "external_ip", "internal_ip", "version"}

func TestAccAgent(t *testing.T) {
	agentID := requireEnv(t, "RUNZERO_TEST_AGENT_ID")
	const res = "runzero_agent.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks:   requireRemovedBlock,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "runzero_agent" "test" {
  id = %q
}

data "runzero_agent" "test" {
  id = runzero_agent.test.id
}

data "runzero_agents" "all" {}
`, agentID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", agentID),
					resource.TestCheckResourceAttrSet(res, "name"),
					resource.TestCheckResourceAttrPair("data.runzero_agent.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttr("data.runzero_agents.all", "agents.*.id", agentID),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateId:           agentID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: volatileAgentAttrs,
			},
			{
				Config: removedConfig(res),
			},
		},
	})
}

func TestAccExplorer(t *testing.T) {
	explorerID := requireEnv(t, "RUNZERO_TEST_EXPLORER_ID")
	const res = "runzero_explorer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks:   requireRemovedBlock,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "runzero_explorer" "test" {
  id = %q
}

data "runzero_explorer" "test" {
  id = runzero_explorer.test.id
}

data "runzero_explorers" "all" {}
`, explorerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", explorerID),
					resource.TestCheckResourceAttrSet(res, "name"),
					resource.TestCheckResourceAttrPair("data.runzero_explorer.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttr("data.runzero_explorers.all", "explorers.*.id", explorerID),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateId:           explorerID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: volatileAgentAttrs,
			},
			{
				Config: removedConfig(res),
			},
		},
	})
}

func TestAccAsset(t *testing.T) {
	assetID := requireEnv(t, "RUNZERO_TEST_ASSET_ID")
	const res = "runzero_asset.test"
	config := func(comments string) string {
		return fmt.Sprintf(`
resource "runzero_asset" "test" {
  id       = %q
  comments = %q
}

data "runzero_asset" "test" {
  id = runzero_asset.test.id
}

data "runzero_assets" "search" {
  search     = "id:%s"
  depends_on = [runzero_asset.test]
}
`, assetID, comments, assetID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOrg(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks:   requireRemovedBlock,
		Steps: []resource.TestStep{
			{
				Config: config("Managed by Terraform acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", assetID),
					resource.TestCheckResourceAttr(res, "comments", "Managed by Terraform acceptance test"),
					resource.TestCheckResourceAttrPair("data.runzero_asset.test", "comments", res, "comments"),
					resource.TestCheckTypeSetElemAttr("data.runzero_assets.search", "assets.*.id", assetID),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateId:           assetID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alive", "last_seen", "updated_at", "risk", "risk_rank"},
			},
			{
				Config: config("Updated by Terraform acceptance test"),
				Check:  resource.TestCheckResourceAttr(res, "comments", "Updated by Terraform acceptance test"),
			},
			{
				Config: removedConfig(res),
			},
		},
	})
}
