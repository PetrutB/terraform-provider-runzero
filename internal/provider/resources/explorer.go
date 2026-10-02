package resources

import (
	"context"
	"errors"
	"fmt"

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
	_ resource.Resource                = &ExplorerResource{}
	_ resource.ResourceWithImportState = &ExplorerResource{}
	_ resource.ResourceWithConfigure   = &ExplorerResource{}
)

type ExplorerResource struct {
	client *client.Client
}

type ExplorerResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SiteID         types.String `tfsdk:"site_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ClientID       types.String `tfsdk:"client_id"`
	OS             types.String `tfsdk:"os"`
	Arch           types.String `tfsdk:"arch"`
	Version        types.String `tfsdk:"version"`
	ExternalIP     types.String `tfsdk:"external_ip"`
	InternalIP     types.String `tfsdk:"internal_ip"`
	Connected      types.Bool   `tfsdk:"connected"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	LastCheckin    types.Int64  `tfsdk:"last_checkin"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
	DeactivatedAt  types.Int64  `tfsdk:"deactivated_at"`
}

func NewExplorerResource() resource.Resource {
	return &ExplorerResource{}
}

func (r *ExplorerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_explorer"
}

func (r *ExplorerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero explorer. Explorers are deployed externally and adopted into Terraform by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the explorer.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the explorer.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID the explorer is associated with.",
				Optional:    true,
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID the explorer belongs to.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os": schema.StringAttribute{
				Description: "The operating system of the explorer host.",
				Computed:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the explorer host.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of the explorer.",
				Computed:    true,
			},
			"external_ip": schema.StringAttribute{
				Description: "The external IP address of the explorer.",
				Computed:    true,
			},
			"internal_ip": schema.StringAttribute{
				Description: "The internal IP address of the explorer.",
				Computed:    true,
			},
			"connected": schema.BoolAttribute{
				Description: "Whether the explorer is currently connected.",
				Computed:    true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the explorer is inactive.",
				Computed:    true,
			},
			"last_checkin": schema.Int64Attribute{
				Description: "Unix timestamp of the explorer's last check-in.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the explorer was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the explorer was last updated.",
				Computed:    true,
			},
			"deactivated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the explorer was deactivated.",
				Computed:    true,
			},
		},
	}
}

func (r *ExplorerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExplorerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ExplorerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	explorer, err := r.client.GetExplorer(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error adopting explorer", fmt.Sprintf("Explorer %s not found: %s", plan.ID.ValueString(), err.Error()))
		return
	}

	if !plan.SiteID.IsNull() && !plan.SiteID.IsUnknown() && plan.SiteID.ValueString() != explorer.SiteID {
		explorer, err = r.client.UpdateExplorer(ctx, explorer.ID, client.UpdateExplorerRequest{
			SiteID: plan.SiteID.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating explorer site", err.Error())
			return
		}
	}

	r.setStateFromExplorer(&plan, explorer)

	tflog.Trace(ctx, "adopted explorer", map[string]interface{}{"id": explorer.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ExplorerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ExplorerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	explorer, err := r.client.GetExplorer(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading explorer", err.Error())
		return
	}

	r.setStateFromExplorer(&state, explorer)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ExplorerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ExplorerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ExplorerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateExplorerRequest{}
	if !plan.SiteID.IsNull() && plan.SiteID.ValueString() != state.SiteID.ValueString() {
		updateReq.SiteID = plan.SiteID.ValueString()
	}

	explorer, err := r.client.UpdateExplorer(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating explorer", err.Error())
		return
	}

	r.setStateFromExplorer(&plan, explorer)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ExplorerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ExplorerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteExplorer(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting explorer", err.Error())
		return
	}
}

func (r *ExplorerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ExplorerResource) setStateFromExplorer(model *ExplorerResourceModel, explorer *client.Explorer) {
	model.ID = types.StringValue(explorer.ID)
	model.Name = types.StringValue(explorer.Name)
	model.SiteID = types.StringValue(explorer.SiteID)
	model.OrganizationID = types.StringValue(explorer.OrganizationID)
	model.ClientID = types.StringValue(explorer.ClientID)
	model.OS = types.StringValue(explorer.OS)
	model.Arch = types.StringValue(explorer.Arch)
	model.Version = types.StringValue(explorer.Version)
	model.ExternalIP = types.StringValue(explorer.ExternalIP)
	model.InternalIP = types.StringValue(explorer.InternalIP)
	model.Connected = types.BoolValue(explorer.Connected)
	model.Inactive = types.BoolValue(explorer.Inactive)
	model.LastCheckin = types.Int64Value(explorer.LastCheckin)
	model.CreatedAt = types.Int64Value(explorer.CreatedAt)
	model.UpdatedAt = types.Int64Value(explorer.UpdatedAt)
	model.DeactivatedAt = types.Int64Value(explorer.DeactivatedAt)
}
