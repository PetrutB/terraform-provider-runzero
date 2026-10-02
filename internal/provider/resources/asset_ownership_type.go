package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource                = &AssetOwnershipTypeResource{}
	_ resource.ResourceWithImportState = &AssetOwnershipTypeResource{}
	_ resource.ResourceWithConfigure   = &AssetOwnershipTypeResource{}
)

type AssetOwnershipTypeResource struct {
	client *client.Client
}

type AssetOwnershipTypeResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Reference types.Int64  `tfsdk:"reference"`
	Order     types.Int64  `tfsdk:"order"`
	Hidden    types.Bool   `tfsdk:"hidden"`
}

func NewAssetOwnershipTypeResource() resource.Resource {
	return &AssetOwnershipTypeResource{}
}

func (r *AssetOwnershipTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_ownership_type"
}

func (r *AssetOwnershipTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero asset ownership type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the asset ownership type.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the asset ownership type.",
				Required:    true,
			},
			"reference": schema.Int64Attribute{
				Description: "Reference number for the ownership type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"order": schema.Int64Attribute{
				Description: "Display order for the ownership type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"hidden": schema.BoolAttribute{
				Description: "Whether the ownership type is hidden.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AssetOwnershipTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AssetOwnershipTypeResource) mapToModel(ot *client.AssetOwnershipType, model *AssetOwnershipTypeResourceModel) {
	model.ID = types.StringValue(ot.ID)
	model.Name = types.StringValue(ot.Name)
	model.Reference = types.Int64Value(ot.Reference)
	model.Order = types.Int64Value(ot.Order)
	model.Hidden = types.BoolValue(ot.Hidden)
}

func (r *AssetOwnershipTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AssetOwnershipTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipType, err := r.client.CreateAssetOwnershipType(ctx, client.CreateAssetOwnershipTypeRequest{
		Name:      plan.Name.ValueString(),
		Reference: plan.Reference.ValueInt64(),
		Order:     plan.Order.ValueInt64(),
		Hidden:    plan.Hidden.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating asset ownership type", err.Error())
		return
	}

	r.mapToModel(ownershipType, &plan)
	tflog.Trace(ctx, "created asset ownership type", map[string]interface{}{"id": ownershipType.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AssetOwnershipTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssetOwnershipTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipType, err := r.client.GetAssetOwnershipType(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading asset ownership type", err.Error())
		return
	}

	r.mapToModel(ownershipType, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AssetOwnershipTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AssetOwnershipTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AssetOwnershipTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipType, err := r.client.UpdateAssetOwnershipType(ctx, state.ID.ValueString(), client.UpdateAssetOwnershipTypeRequest{
		Name:      plan.Name.ValueString(),
		Reference: plan.Reference.ValueInt64(),
		Order:     plan.Order.ValueInt64(),
		Hidden:    plan.Hidden.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating asset ownership type", err.Error())
		return
	}

	r.mapToModel(ownershipType, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AssetOwnershipTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AssetOwnershipTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAssetOwnershipType(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting asset ownership type", err.Error())
		return
	}
}

func (r *AssetOwnershipTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
