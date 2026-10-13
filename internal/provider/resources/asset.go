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
	_ resource.Resource                = &AssetResource{}
	_ resource.ResourceWithImportState = &AssetResource{}
	_ resource.ResourceWithConfigure   = &AssetResource{}
)

type AssetResource struct {
	client *client.Client
}

type AssetResourceModel struct {
	ID          types.String  `tfsdk:"id"`
	Type        types.String  `tfsdk:"type"`
	OS          types.String  `tfsdk:"os"`
	OSVersion   types.String  `tfsdk:"os_version"`
	HW          types.String  `tfsdk:"hw"`
	SiteID      types.String  `tfsdk:"site_id"`
	Comments    types.String  `tfsdk:"comments"`
	Criticality types.String  `tfsdk:"criticality"`
	Risk        types.String  `tfsdk:"risk"`
	RiskRank    types.Float64 `tfsdk:"risk_rank"`
	Alive       types.Bool    `tfsdk:"alive"`
	FirstSeen   types.Int64   `tfsdk:"first_seen"`
	LastSeen    types.Int64   `tfsdk:"last_seen"`
	CreatedAt   types.Int64   `tfsdk:"created_at"`
	UpdatedAt   types.Int64   `tfsdk:"updated_at"`
}

func NewAssetResource() resource.Resource {
	return &AssetResource{}
}

func (r *AssetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (r *AssetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero asset. Assets are discovered via scans and adopted into Terraform by ID. Supports updating comments and criticality.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the asset.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the asset (e.g. Server, Workstation).",
				Computed:    true,
			},
			"os": schema.StringAttribute{
				Description: "The operating system of the asset.",
				Computed:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "The OS version of the asset.",
				Computed:    true,
			},
			"hw": schema.StringAttribute{
				Description: "The hardware description of the asset.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID the asset belongs to.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"comments": schema.StringAttribute{
				Description: "Comments for the asset.",
				Optional:    true,
				Computed:    true,
			},
			"criticality": schema.StringAttribute{
				Description: "The criticality level of the asset (e.g. low, medium, high, critical).",
				Optional:    true,
				Computed:    true,
			},
			"risk": schema.StringAttribute{
				Description: "The risk level of the asset.",
				Computed:    true,
			},
			"risk_rank": schema.Float64Attribute{
				Description: "The numeric risk rank (-1=none, 0=info, 4=critical).",
				Computed:    true,
			},
			"alive": schema.BoolAttribute{
				Description: "Whether the asset is alive.",
				Computed:    true,
			},
			"first_seen": schema.Int64Attribute{
				Description: "Unix timestamp when the asset was first seen.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_seen": schema.Int64Attribute{
				Description: "Unix timestamp when the asset was last seen.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the asset was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the asset was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *AssetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AssetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := r.client.GetAsset(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error adopting asset", fmt.Sprintf("Asset %s not found: %s", plan.ID.ValueString(), err.Error()))
		return
	}

	if !plan.Comments.IsNull() && !plan.Comments.IsUnknown() && plan.Comments.ValueString() != asset.Comments {
		asset, err = r.client.UpdateAssetComments(ctx, asset.ID, client.UpdateAssetCommentsRequest{
			Comments: plan.Comments.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating asset comments", err.Error())
			return
		}
	}

	if !plan.Criticality.IsNull() && !plan.Criticality.IsUnknown() {
		asset, err = r.client.UpdateAssetCriticality(ctx, asset.ID, client.UpdateAssetCriticalityRequest{
			Criticality: plan.Criticality.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating asset criticality", err.Error())
			return
		}
	}

	r.setStateFromAsset(&plan, asset)

	tflog.Trace(ctx, "adopted asset", map[string]interface{}{"id": asset.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := r.client.GetAsset(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading asset", err.Error())
		return
	}

	r.setStateFromAsset(&state, asset)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AssetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AssetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	if !plan.Comments.IsNull() && plan.Comments.ValueString() != state.Comments.ValueString() {
		_, err := r.client.UpdateAssetComments(ctx, id, client.UpdateAssetCommentsRequest{
			Comments: plan.Comments.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating asset comments", err.Error())
			return
		}
	}

	if !plan.Criticality.IsNull() && plan.Criticality.ValueString() != state.Criticality.ValueString() {
		_, err := r.client.UpdateAssetCriticality(ctx, id, client.UpdateAssetCriticalityRequest{
			Criticality: plan.Criticality.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating asset criticality", err.Error())
			return
		}
	}

	asset, err := r.client.GetAsset(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading asset after update", err.Error())
		return
	}

	r.setStateFromAsset(&plan, asset)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AssetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAsset(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting asset", err.Error())
		return
	}
}

func (r *AssetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AssetResource) setStateFromAsset(model *AssetResourceModel, asset *client.Asset) {
	model.ID = types.StringValue(asset.ID)
	model.Type = types.StringValue(asset.Type)
	model.OS = types.StringValue(asset.OS)
	model.OSVersion = types.StringValue(asset.OSVersion)
	model.HW = types.StringValue(asset.HW)
	model.SiteID = types.StringValue(asset.SiteID)
	model.Comments = types.StringValue(asset.Comments)
	model.Criticality = types.StringValue(asset.Criticality)
	model.Risk = types.StringValue(asset.Risk)
	model.RiskRank = types.Float64Value(asset.RiskRank)
	model.Alive = types.BoolValue(asset.Alive)
	model.FirstSeen = types.Int64Value(asset.FirstSeen)
	model.LastSeen = types.Int64Value(asset.LastSeen)
	model.CreatedAt = types.Int64Value(asset.CreatedAt)
	model.UpdatedAt = types.Int64Value(asset.UpdatedAt)
}
