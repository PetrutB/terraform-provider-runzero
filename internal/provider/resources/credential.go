package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource                = &CredentialResource{}
	_ resource.ResourceWithImportState = &CredentialResource{}
	_ resource.ResourceWithConfigure   = &CredentialResource{}
)

type CredentialResource struct {
	client *client.Client
}

type CredentialResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ClientID       types.String `tfsdk:"client_id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Secret         types.Map    `tfsdk:"secret"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedByID    types.String `tfsdk:"created_by_id"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	ACL            types.Map    `tfsdk:"acl"`
	Global         types.Bool   `tfsdk:"global"`
	CIDRs          types.List   `tfsdk:"cidrs"`
	LastUsedAt     types.Int64  `tfsdk:"last_used_at"`
	LastUsedByID   types.String `tfsdk:"last_used_by_id"`
}

func NewCredentialResource() resource.Resource {
	return &CredentialResource{}
}

func (r *CredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (r *CredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero credential.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the credential.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the credential.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the credential.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"aws_access_secret",
						"azure_client_secret",
						"azure_username_password",
						"censys_search_api_key",
						"crowdstrike_api_key",
						"custom_integration_credential",
						"gcp_service_account_key",
						"google_workspace_client_secret",
						"insightvm_username_password",
						"ldap_username_password",
						"mecm_db_connection_string",
						"meraki_api_key",
						"miradore_api_key_v1",
						"nessus_access_secret",
						"netbox_api_key",
						"palo_alto_api_key",
						"qualys_username_password",
						"sentinelone_api_key",
						"shodan_api_key",
						"snmpv2_community",
						"snmpv3_credential",
						"tanium_api_token",
						"tenable_access_secret",
						"tenable_security_center_access_secret",
						"vmware_username_password",
						"wiz_client_secret",
						"x509_keypair",
					),
				},
			},
			"secret": schema.MapAttribute{
				Description: "Type-specific secret material. Keys depend on the credential type (e.g. AWS: access_key/secret_access_key; SNMPv2: community; LDAP: url/username/password). Write-only; never returned by the API.",
				Optional:    true,
				Sensitive:   true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the credential was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_by_id": schema.StringAttribute{
				Description: "UUID of the user who created the credential.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by_email": schema.StringAttribute{
				Description: "Email of the user who created the credential.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"acl": schema.MapAttribute{
				Description: "Access control list mapping organization UUIDs to permission levels.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"global": schema.BoolAttribute{
				Description: "Whether the credential is globally accessible.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"cidrs": schema.ListAttribute{
				Description: "List of allowed CIDR ranges.",
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"last_used_at": schema.Int64Attribute{
				Description: "Unix timestamp when the credential was last used.",
				Computed:    true,
			},
			"last_used_by_id": schema.StringAttribute{
				Description: "UUID of the user who last used the credential.",
				Computed:    true,
			},
		},
	}
}

func (r *CredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *CredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aclMap := make(map[string]string)
	if !plan.ACL.IsNull() && !plan.ACL.IsUnknown() {
		resp.Diagnostics.Append(plan.ACL.ElementsAs(ctx, &aclMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	secretMap := make(map[string]string)
	if !plan.Secret.IsNull() && !plan.Secret.IsUnknown() {
		resp.Diagnostics.Append(plan.Secret.ElementsAs(ctx, &secretMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	credential, err := r.client.CreateCredential(ctx, client.CreateCredentialRequest{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		ACL:    aclMap,
		Global: plan.Global.ValueBool(),
		Secret: secretMap,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating credential", err.Error())
		return
	}

	mapCredentialToState(ctx, credential, &plan, &resp.Diagnostics)

	tflog.Trace(ctx, "created credential", map[string]interface{}{"id": credential.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credential, err := r.client.GetCredential(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading credential", err.Error())
		return
	}

	mapCredentialToState(ctx, credential, &state, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Credentials cannot be updated. Change attributes to trigger a replacement.",
	)
}

func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCredential(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting credential", err.Error())
		return
	}
}

func (r *CredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapCredentialToState(ctx context.Context, cred *client.Credential, state *CredentialResourceModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(cred.ID)
	state.ClientID = types.StringValue(cred.ClientID)
	state.Name = types.StringValue(cred.Name)
	state.Type = types.StringValue(cred.Type)
	state.CreatedAt = types.Int64Value(cred.CreatedAt)
	state.CreatedByID = types.StringValue(cred.CreatedByID)
	state.CreatedByEmail = types.StringValue(cred.CreatedByEmail)
	state.Global = types.BoolValue(cred.Global)
	state.LastUsedAt = types.Int64Value(cred.LastUsedAt)
	state.LastUsedByID = types.StringValue(cred.LastUsedByID)

	aclMap, d := types.MapValueFrom(ctx, types.StringType, cred.ACL)
	diags.Append(d...)
	state.ACL = aclMap

	cidrList, d := types.ListValueFrom(ctx, types.StringType, cred.CIDRs)
	diags.Append(d...)
	state.CIDRs = cidrList
}
