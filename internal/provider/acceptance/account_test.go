package acceptance

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

func TestAccOrganization(t *testing.T) {
	name := randomName()
	const res = "runzero_organization.test"
	config := func(description string) string {
		return fmt.Sprintf(`
resource "runzero_organization" "test" {
  name        = %q
  description = %q
}

data "runzero_organization" "test" {
  id = runzero_organization.test.id
}

data "runzero_organizations" "all" {
  depends_on = [runzero_organization.test]
}
`, name, description)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_organization", byID((*client.Client).GetOrganization)),
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "name", name),
					resource.TestCheckResourceAttr(res, "description", "Terraform acceptance test"),
					resource.TestCheckResourceAttrPair("data.runzero_organization.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_organizations.all", "organizations.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config("Updated by Terraform acceptance test"),
				Check:  resource.TestCheckResourceAttr(res, "description", "Updated by Terraform acceptance test"),
			},
		},
	})
}

func TestAccUser(t *testing.T) {
	email := randomEmail()
	const res = "runzero_user.test"
	config := func(lastName string) string {
		return fmt.Sprintf(`
resource "runzero_user" "test" {
  email            = %q
  first_name       = "Terraform"
  last_name        = %q
  org_default_role = "viewer"
}

data "runzero_user" "test" {
  id = runzero_user.test.id
}

data "runzero_users" "all" {
  depends_on = [runzero_user.test]
}
`, email, lastName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_user", byID((*client.Client).GetUser)),
		Steps: []resource.TestStep{
			{
				Config: config("Acceptance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "email", email),
					resource.TestCheckResourceAttr(res, "org_default_role", "viewer"),
					resource.TestCheckResourceAttrPair("data.runzero_user.test", "email", res, "email"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_users.all", "users.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config("Updated"),
				Check:  resource.TestCheckResourceAttr(res, "last_name", "Updated"),
			},
		},
	})
}

func TestAccUserInvite(t *testing.T) {
	email := requireEnv(t, "RUNZERO_TEST_INVITE_EMAIL")
	const res = "runzero_user_invite.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_user_invite", byID((*client.Client).GetUser)),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "runzero_user_invite" "test" {
  email            = %q
  first_name       = "Terraform"
  last_name        = "Invite"
  org_default_role = "viewer"
  subject          = "runZero Terraform acceptance test"
  message          = "This invitation was sent by an automated test and can be ignored."
}
`, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "email", email),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subject", "message"},
			},
		},
	})
}

func TestAccGroup(t *testing.T) {
	name := randomName()
	const res = "runzero_group.test"
	config := func(description string) string {
		return fmt.Sprintf(`
resource "runzero_group" "test" {
  name             = %q
  description      = %q
  org_default_role = "viewer"
}

data "runzero_group" "test" {
  id = runzero_group.test.id
}

data "runzero_groups" "all" {
  depends_on = [runzero_group.test]
}
`, name, description)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_group", byID((*client.Client).GetGroup)),
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "name", name),
					resource.TestCheckResourceAttr(res, "org_default_role", "viewer"),
					resource.TestCheckResourceAttrPair("data.runzero_group.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_groups.all", "groups.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config("Updated by Terraform acceptance test"),
				Check:  resource.TestCheckResourceAttr(res, "description", "Updated by Terraform acceptance test"),
			},
		},
	})
}

func TestAccSSOGroupMapping(t *testing.T) {
	name := randomName()
	const res = "runzero_sso_group_mapping.test"
	config := func(value string) string {
		return fmt.Sprintf(`
resource "runzero_group" "test" {
  name             = %q
  org_default_role = "viewer"
}

resource "runzero_sso_group_mapping" "test" {
  group_id      = runzero_group.test.id
  sso_attribute = "groups"
  sso_value     = %q
  description   = "Terraform acceptance test"
}

data "runzero_sso_group_mapping" "test" {
  id = runzero_sso_group_mapping.test.id
}

data "runzero_sso_group_mappings" "all" {
  depends_on = [runzero_sso_group_mapping.test]
}
`, name, value)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_sso_group_mapping", byID((*client.Client).GetSSOGroupMapping)),
		Steps: []resource.TestStep{
			{
				Config: config(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttrPair(res, "group_id", "runzero_group.test", "id"),
					resource.TestCheckResourceAttr(res, "sso_value", name),
					resource.TestCheckResourceAttrPair("data.runzero_sso_group_mapping.test", "sso_value", res, "sso_value"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_sso_group_mappings.all", "sso_group_mappings.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config(name + "-updated"),
				Check:  resource.TestCheckResourceAttr(res, "sso_value", name+"-updated"),
			},
		},
	})
}

func TestAccAPIKey(t *testing.T) {
	const res = "runzero_api_key.test"
	var token string
	config := func(comment, rotate string) string {
		return fmt.Sprintf(`
resource "runzero_api_key" "test" {
  comment = %q
  rotate  = %q
}

data "runzero_api_key" "test" {
  id = runzero_api_key.test.id
}

data "runzero_api_keys" "all" {
  depends_on = [runzero_api_key.test]
}
`, comment, rotate)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_api_key", byID((*client.Client).GetAPIKey)),
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test", "1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttrSet(res, "token"),
					resource.TestCheckResourceAttr(res, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttrPair("data.runzero_api_key.test", "comment", res, "comment"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_api_keys.all", "api_keys.*.id", res, "id"),
					captureAttr(res, "token", &token),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotate"},
			},
			{
				Config: config("Updated by Terraform acceptance test", "1"),
				Check:  resource.TestCheckResourceAttr(res, "comment", "Updated by Terraform acceptance test"),
			},
			{
				Config: config("Updated by Terraform acceptance test", "2"),
				Check:  checkAttrChanged(res, "token", &token),
			},
		},
	})
}

func TestAccExportToken(t *testing.T) {
	name := randomName()
	const res = "runzero_export_token.test"
	var token string
	config := func(comment, rotate string) string {
		return fmt.Sprintf(`
resource "runzero_organization" "test" {
  name = %q
}

resource "runzero_export_token" "test" {
  org_id  = runzero_organization.test.id
  comment = %q
  rotate  = %q
}

data "runzero_export_token" "test" {
  org_id = runzero_organization.test.id
  id     = runzero_export_token.test.id
}

data "runzero_export_tokens" "all" {
  org_id     = runzero_organization.test.id
  depends_on = [runzero_export_token.test]
}
`, name, comment, rotate)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: checkDestroy("runzero_export_token", func(ctx context.Context, c *client.Client, rs *terraform.ResourceState) error {
			_, err := c.GetExportToken(ctx, rs.Primary.Attributes["org_id"], rs.Primary.ID)
			return err
		}),
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test", "1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttrSet(res, "token"),
					resource.TestCheckResourceAttrPair(res, "org_id", "runzero_organization.test", "id"),
					resource.TestCheckResourceAttrPair("data.runzero_export_token.test", "comment", res, "comment"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_export_tokens.all", "export_tokens.*.id", res, "id"),
					captureAttr(res, "token", &token),
				),
			},
			{
				Config: config("Terraform acceptance test", "2"),
				Check:  checkAttrChanged(res, "token", &token),
			},
		},
	})
}

func TestAccCredential(t *testing.T) {
	name := randomName()
	const res = "runzero_credential.test"
	config := func(credName string) string {
		return fmt.Sprintf(`
resource "runzero_credential" "test" {
  name   = %q
  type   = "snmpv2_community"
  global = true
  secret = {
    community = "tf-acc-community"
  }
}

data "runzero_credential" "test" {
  id = runzero_credential.test.id
}

data "runzero_credentials" "all" {
  depends_on = [runzero_credential.test]
}
`, credName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_credential", byID((*client.Client).GetCredential)),
		Steps: []resource.TestStep{
			{
				Config: config(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "type", "snmpv2_community"),
					resource.TestCheckResourceAttr(res, "global", "true"),
					resource.TestCheckResourceAttrPair("data.runzero_credential.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_credentials.all", "credentials.*.id", res, "id"),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
			{
				Config: config(name + "-updated"),
				Check:  resource.TestCheckResourceAttr(res, "name", name+"-updated"),
			},
		},
	})
}

func TestAccCustomIntegration(t *testing.T) {
	name := randomName()
	const res = "runzero_custom_integration.test"
	config := func(description string) string {
		return fmt.Sprintf(`
resource "runzero_custom_integration" "test" {
  name        = %q
  description = %q
}

data "runzero_custom_integration" "test" {
  id = runzero_custom_integration.test.id
}

data "runzero_custom_integrations" "all" {
  depends_on = [runzero_custom_integration.test]
}
`, name, description)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_custom_integration", byID((*client.Client).GetCustomIntegration)),
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "name", name),
					resource.TestCheckResourceAttrPair("data.runzero_custom_integration.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_custom_integrations.all", "custom_integrations.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config("Updated by Terraform acceptance test"),
				Check:  resource.TestCheckResourceAttr(res, "description", "Updated by Terraform acceptance test"),
			},
		},
	})
}

func TestAccAssetOwnershipType(t *testing.T) {
	name := randomName()
	const res = "runzero_asset_ownership_type.test"
	config := func(typeName string) string {
		return fmt.Sprintf(`
resource "runzero_asset_ownership_type" "test" {
  name  = %q
  order = 99
}

data "runzero_asset_ownership_type" "test" {
  id = runzero_asset_ownership_type.test.id
}

data "runzero_asset_ownership_types" "all" {
  depends_on = [runzero_asset_ownership_type.test]
}
`, typeName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_asset_ownership_type", byID((*client.Client).GetAssetOwnershipType)),
		Steps: []resource.TestStep{
			{
				Config: config(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "name", name),
					resource.TestCheckResourceAttr(res, "order", "99"),
					resource.TestCheckResourceAttrPair("data.runzero_asset_ownership_type.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_asset_ownership_types.all", "asset_ownership_types.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config(name + "-updated"),
				Check:  resource.TestCheckResourceAttr(res, "name", name+"-updated"),
			},
		},
	})
}

func TestAccScanTemplate(t *testing.T) {
	name := randomName()
	const res = "runzero_scan_template.test"
	config := func(description string) string {
		return fmt.Sprintf(`
resource "runzero_organization" "test" {
  name = %[1]q
}

resource "runzero_scan_template" "test" {
  name            = %[1]q
  description     = %[2]q
  organization_id = runzero_organization.test.id
  global          = false
  acl             = jsonencode({ (runzero_organization.test.id) = "user" })
}

data "runzero_scan_template" "test" {
  id = runzero_scan_template.test.id
}

data "runzero_scan_templates" "all" {
  depends_on = [runzero_scan_template.test]
}
`, name, description)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckAccount(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy("runzero_scan_template", byID((*client.Client).GetScanTemplate)),
		Steps: []resource.TestStep{
			{
				Config: config("Terraform acceptance test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(res, "id"),
					resource.TestCheckResourceAttr(res, "name", name),
					resource.TestCheckResourceAttr(res, "global", "false"),
					resource.TestCheckResourceAttrPair("data.runzero_scan_template.test", "name", res, "name"),
					resource.TestCheckTypeSetElemAttrPair("data.runzero_scan_templates.all", "scan_templates.*.id", res, "id"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: config("Updated by Terraform acceptance test"),
				Check:  resource.TestCheckResourceAttr(res, "description", "Updated by Terraform acceptance test"),
			},
		},
	})
}
