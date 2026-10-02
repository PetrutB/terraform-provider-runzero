package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
	_ resource.ResourceWithConfigure   = &OrganizationResource{}
)

type OrganizationResource struct {
	client *client.Client
}

type OrganizationResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	ClientID                        types.String `tfsdk:"client_id"`
	Inactive                        types.Bool   `tfsdk:"inactive"`
	Project                         types.Bool   `tfsdk:"project"`
	ParentID                        types.String `tfsdk:"parent_id"`
	AssetCount                      types.Int64  `tfsdk:"asset_count"`
	ServiceCount                    types.Int64  `tfsdk:"service_count"`
	ExpirationAssetsStale           types.String `tfsdk:"expiration_assets_stale"`
	ExpirationAssetsOffline         types.String `tfsdk:"expiration_assets_offline"`
	ExpirationIntegrationAttributes types.String `tfsdk:"expiration_integration_attributes"`
	ExpirationScans                 types.String `tfsdk:"expiration_scans"`
	ExpirationVulnerabilities       types.String `tfsdk:"expiration_vulnerabilities"`
	KeepLatestIntegrationAttributes types.String `tfsdk:"keep_latest_integration_attributes"`
	CreatedAt                       types.Int64  `tfsdk:"created_at"`
	UpdatedAt                       types.Int64  `tfsdk:"updated_at"`
}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

func (r *OrganizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the organization.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the organization.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the organization.",
				Optional:    true,
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID of the organization.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the organization is inactive.",
				Computed:    true,
			},
			"project": schema.BoolAttribute{
				Description: "Whether this is a project organization.",
				Computed:    true,
			},
			"parent_id": schema.StringAttribute{
				Description: "The parent organization ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asset_count": schema.Int64Attribute{
				Description: "Number of assets in the organization.",
				Computed:    true,
			},
			"service_count": schema.Int64Attribute{
				Description: "Number of services in the organization.",
				Computed:    true,
			},
			"expiration_assets_stale": schema.StringAttribute{
				Description: "Duration before assets are marked stale.",
				Optional:    true,
				Computed:    true,
			},
			"expiration_assets_offline": schema.StringAttribute{
				Description: "Duration before assets are marked offline.",
				Optional:    true,
				Computed:    true,
			},
			"expiration_integration_attributes": schema.StringAttribute{
				Description: "Duration before integration attributes expire.",
				Optional:    true,
				Computed:    true,
			},
			"expiration_scans": schema.StringAttribute{
				Description: "Duration before scans expire.",
				Optional:    true,
				Computed:    true,
			},
			"expiration_vulnerabilities": schema.StringAttribute{
				Description: "Duration before vulnerabilities expire.",
				Optional:    true,
				Computed:    true,
			},
			"keep_latest_integration_attributes": schema.StringAttribute{
				Description: "Whether to keep the latest integration attributes (as string boolean).",
				Optional:    true,
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the organization was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the organization was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *OrganizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.CreateOrganization(ctx, client.CreateOrganizationRequest{
		Name:                            plan.Name.ValueString(),
		Description:                     plan.Description.ValueString(),
		ParentID:                        plan.ParentID.ValueString(),
		ExpirationAssetsStale:           plan.ExpirationAssetsStale.ValueString(),
		ExpirationAssetsOffline:         plan.ExpirationAssetsOffline.ValueString(),
		ExpirationIntegrationAttributes: plan.ExpirationIntegrationAttributes.ValueString(),
		ExpirationScans:                 plan.ExpirationScans.ValueString(),
		ExpirationVulnerabilities:       plan.ExpirationVulnerabilities.ValueString(),
		KeepLatestIntegrationAttributes: plan.KeepLatestIntegrationAttributes.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating organization", err.Error())
		return
	}

	mapOrgToState(org, &plan)
	tflog.Trace(ctx, "created organization", map[string]interface{}{"id": org.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.GetOrganization(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading organization", err.Error())
		return
	}

	mapOrgToState(org, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.UpdateOrganization(ctx, state.ID.ValueString(), client.UpdateOrganizationRequest{
		Name:                            plan.Name.ValueString(),
		Description:                     plan.Description.ValueString(),
		ExpirationAssetsStale:           plan.ExpirationAssetsStale.ValueString(),
		ExpirationAssetsOffline:         plan.ExpirationAssetsOffline.ValueString(),
		ExpirationIntegrationAttributes: plan.ExpirationIntegrationAttributes.ValueString(),
		ExpirationScans:                 plan.ExpirationScans.ValueString(),
		ExpirationVulnerabilities:       plan.ExpirationVulnerabilities.ValueString(),
		KeepLatestIntegrationAttributes: plan.KeepLatestIntegrationAttributes.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating organization", err.Error())
		return
	}

	mapOrgToState(org, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrganization(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting organization", err.Error())
		return
	}
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapOrgToState(org *client.Organization, state *OrganizationResourceModel) {
	state.ID = types.StringValue(org.ID)
	state.Name = types.StringValue(org.Name)
	state.Description = types.StringValue(org.Description)
	state.ClientID = types.StringValue(org.ClientID)
	state.Inactive = types.BoolValue(org.Inactive)
	state.Project = types.BoolValue(org.Project)
	state.ParentID = types.StringValue(org.ParentID)
	state.AssetCount = types.Int64Value(org.AssetCount)
	state.ServiceCount = types.Int64Value(org.ServiceCount)
	state.ExpirationAssetsStale = types.StringValue(strconv.FormatInt(org.ExpirationAssetsStale, 10))
	state.ExpirationAssetsOffline = types.StringValue(strconv.FormatInt(org.ExpirationAssetsOffline, 10))
	state.ExpirationIntegrationAttributes = types.StringValue(strconv.FormatInt(org.ExpirationIntegrationAttributes, 10))
	state.ExpirationScans = types.StringValue(strconv.FormatInt(org.ExpirationScans, 10))
	state.ExpirationVulnerabilities = types.StringValue(strconv.FormatInt(org.ExpirationVulnerabilities, 10))
	state.KeepLatestIntegrationAttributes = types.StringValue(strconv.FormatBool(org.KeepLatestIntegrationAttributes))
	state.CreatedAt = types.Int64Value(org.CreatedAt)
	state.UpdatedAt = types.Int64Value(org.UpdatedAt)
}
