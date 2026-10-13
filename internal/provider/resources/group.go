package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
	"github.com/petrutb/terraform-provider-runzero/internal/provider/planmodifiers"
)

var (
	_ resource.Resource                = &GroupResource{}
	_ resource.ResourceWithImportState = &GroupResource{}
	_ resource.ResourceWithConfigure   = &GroupResource{}
)

type GroupResource struct {
	client *client.Client
}

type GroupResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	RoleSummary    types.String `tfsdk:"role_summary"`
	UserCount      types.Int64  `tfsdk:"user_count"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	ExpiresAt      types.Int64  `tfsdk:"expires_at"`
	OrgDefaultRole types.String `tfsdk:"org_default_role"`
	OrgRoles       types.Map    `tfsdk:"org_roles"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the group.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the group.",
				Optional:    true,
				Computed:    true,
			},
			"role_summary": schema.StringAttribute{
				Description: "Summary of roles assigned to the group.",
				Computed:    true,
			},
			"user_count": schema.Int64Attribute{
				Description: "Number of users in the group.",
				Computed:    true,
			},
			"created_by_email": schema.StringAttribute{
				Description: "Email of the user who created the group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.Int64Attribute{
				Description: "Unix timestamp when the group expires.",
				Optional:    true,
				Computed:    true,
			},
			"org_default_role": schema.StringAttribute{
				Description: "The default organization role for the group.",
				Optional:    true,
				Computed:    true,
			},
			"org_roles": schema.MapAttribute{
				Description: "Map of organization ID to role for the group.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					planmodifiers.ServerEnrichedMap(),
				},
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the group was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the group was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *GroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgRoles := extractMapString(ctx, plan.OrgRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.CreateGroup(ctx, client.CreateGroupRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		ExpiresAt:      plan.ExpiresAt.ValueInt64(),
		OrgDefaultRole: plan.OrgDefaultRole.ValueString(),
		OrgRoles:       orgRoles,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating group", err.Error())
		return
	}

	configOrgRoles := plan.OrgRoles
	mapGroupToState(ctx, group, &plan)
	plan.OrgRoles = configOrgRoles
	tflog.Trace(ctx, "created group", map[string]interface{}{"id": group.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetGroup(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading group", err.Error())
		return
	}

	priorOrgRoles := state.OrgRoles
	mapGroupToState(ctx, group, &state)
	if !priorOrgRoles.IsNull() && !priorOrgRoles.IsUnknown() {
		state.OrgRoles = filterMapToKeys(ctx, state.OrgRoles, priorOrgRoles, &resp.Diagnostics)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgRoles := extractMapString(ctx, plan.OrgRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.UpdateGroup(ctx, client.UpdateGroupRequest{
		ID:             state.ID.ValueString(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		ExpiresAt:      plan.ExpiresAt.ValueInt64(),
		OrgDefaultRole: plan.OrgDefaultRole.ValueString(),
		OrgRoles:       orgRoles,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating group", err.Error())
		return
	}

	configOrgRoles := plan.OrgRoles
	mapGroupToState(ctx, group, &plan)
	plan.OrgRoles = configOrgRoles
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroup(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", err.Error())
		return
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapGroupToState(ctx context.Context, group *client.Group, state *GroupResourceModel) {
	_ = ctx
	state.ID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)
	state.RoleSummary = types.StringValue(group.RoleSummary)
	state.UserCount = types.Int64Value(group.UserCount)
	state.CreatedByEmail = types.StringValue(group.CreatedByEmail)
	state.ExpiresAt = types.Int64Value(group.ExpiresAt)
	state.OrgDefaultRole = types.StringValue(group.OrgDefaultRole)
	state.CreatedAt = types.Int64Value(group.CreatedAt)
	state.UpdatedAt = types.Int64Value(group.UpdatedAt)

	if group.OrgRoles != nil {
		elements := make(map[string]attr.Value, len(group.OrgRoles))
		for k, v := range group.OrgRoles {
			elements[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		state.OrgRoles = types.MapValueMust(types.StringType, elements)
	} else {
		state.OrgRoles = types.MapNull(types.StringType)
	}
}
