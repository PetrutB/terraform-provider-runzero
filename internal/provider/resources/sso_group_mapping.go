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
	_ resource.Resource                = &SSOGroupMappingResource{}
	_ resource.ResourceWithImportState = &SSOGroupMappingResource{}
	_ resource.ResourceWithConfigure   = &SSOGroupMappingResource{}
)

type SSOGroupMappingResource struct {
	client *client.Client
}

type SSOGroupMappingResourceModel struct {
	ID             types.String `tfsdk:"id"`
	GroupID        types.String `tfsdk:"group_id"`
	GroupName      types.String `tfsdk:"group_name"`
	SSOAttribute   types.String `tfsdk:"sso_attribute"`
	SSOValue       types.String `tfsdk:"sso_value"`
	Description    types.String `tfsdk:"description"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewSSOGroupMappingResource() resource.Resource {
	return &SSOGroupMappingResource{}
}

func (r *SSOGroupMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_group_mapping"
}

func (r *SSOGroupMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero SSO group mapping.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the SSO group mapping.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "The RunZero group ID to map to.",
				Required:    true,
			},
			"group_name": schema.StringAttribute{
				Description: "The name of the mapped group (read-only).",
				Computed:    true,
			},
			"sso_attribute": schema.StringAttribute{
				Description: "The SSO attribute name for this mapping.",
				Required:    true,
			},
			"sso_value": schema.StringAttribute{
				Description: "The SSO attribute value for this mapping.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the SSO group mapping.",
				Optional:    true,
				Computed:    true,
			},
			"created_by_email": schema.StringAttribute{
				Description: "Email of the user who created the mapping.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the SSO group mapping was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the SSO group mapping was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *SSOGroupMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSOGroupMappingResource) mapToModel(mapping *client.SSOGroupMapping, model *SSOGroupMappingResourceModel) {
	model.ID = types.StringValue(mapping.ID)
	model.GroupID = types.StringValue(mapping.GroupID)
	model.GroupName = types.StringValue(mapping.GroupName)
	model.SSOAttribute = types.StringValue(mapping.SSOAttribute)
	model.SSOValue = types.StringValue(mapping.SSOValue)
	model.Description = types.StringValue(mapping.Description)
	model.CreatedByEmail = types.StringValue(mapping.CreatedByEmail)
	model.CreatedAt = types.Int64Value(mapping.CreatedAt)
	model.UpdatedAt = types.Int64Value(mapping.UpdatedAt)
}

func (r *SSOGroupMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSOGroupMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.CreateSSOGroupMapping(ctx, client.CreateSSOGroupMappingRequest{
		GroupID:      plan.GroupID.ValueString(),
		SSOAttribute: plan.SSOAttribute.ValueString(),
		SSOValue:     plan.SSOValue.ValueString(),
		Description:  plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSO group mapping", err.Error())
		return
	}

	r.mapToModel(mapping, &plan)
	tflog.Trace(ctx, "created SSO group mapping", map[string]interface{}{"id": mapping.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SSOGroupMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSOGroupMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.GetSSOGroupMapping(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SSO group mapping", err.Error())
		return
	}

	r.mapToModel(mapping, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SSOGroupMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SSOGroupMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.UpdateSSOGroupMapping(ctx, client.UpdateSSOGroupMappingRequest{
		ID:           plan.ID.ValueString(),
		GroupID:      plan.GroupID.ValueString(),
		SSOAttribute: plan.SSOAttribute.ValueString(),
		SSOValue:     plan.SSOValue.ValueString(),
		Description:  plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating SSO group mapping", err.Error())
		return
	}

	r.mapToModel(mapping, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SSOGroupMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSOGroupMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSSOGroupMapping(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting SSO group mapping", err.Error())
		return
	}
}

func (r *SSOGroupMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
