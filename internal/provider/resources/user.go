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
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
)

type UserResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Email          types.String `tfsdk:"email"`
	FirstName      types.String `tfsdk:"first_name"`
	LastName       types.String `tfsdk:"last_name"`
	ClientAdmin    types.Bool   `tfsdk:"client_admin"`
	OrgDefaultRole types.String `tfsdk:"org_default_role"`
	OrgRoles       types.Map    `tfsdk:"org_roles"`
	SSOOnly        types.Bool   `tfsdk:"sso_only"`
	MFAEnabled     types.Bool   `tfsdk:"mfa_enabled"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the user.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The email address of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"first_name": schema.StringAttribute{
				Description: "The first name of the user.",
				Optional:    true,
				Computed:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "The last name of the user.",
				Optional:    true,
				Computed:    true,
			},
			"client_admin": schema.BoolAttribute{
				Description: "Whether the user is a client administrator.",
				Optional:    true,
				Computed:    true,
			},
			"org_default_role": schema.StringAttribute{
				Description: "The default organization role for the user.",
				Optional:    true,
				Computed:    true,
			},
			"org_roles": schema.MapAttribute{
				Description: "Map of organization ID to role for the user.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					planmodifiers.ServerEnrichedMap(),
				},
			},
			"sso_only": schema.BoolAttribute{
				Description: "Whether the user can only authenticate via SSO.",
				Computed:    true,
			},
			"mfa_enabled": schema.BoolAttribute{
				Description: "Whether MFA is enabled for the user.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the user was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the user was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgRoles := extractMapString(ctx, plan.OrgRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.CreateUser(ctx, client.CreateUserRequest{
		Email:          plan.Email.ValueString(),
		FirstName:      plan.FirstName.ValueString(),
		LastName:       plan.LastName.ValueString(),
		ClientAdmin:    plan.ClientAdmin.ValueBool(),
		OrgDefaultRole: plan.OrgDefaultRole.ValueString(),
		OrgRoles:       orgRoles,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	configOrgRoles := plan.OrgRoles
	mapUserToState(user, &plan)
	plan.OrgRoles = configOrgRoles
	tflog.Trace(ctx, "created user", map[string]interface{}{"id": user.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	priorOrgRoles := state.OrgRoles
	mapUserToState(user, &state)
	if !priorOrgRoles.IsNull() && !priorOrgRoles.IsUnknown() {
		state.OrgRoles = filterMapToKeys(ctx, state.OrgRoles, priorOrgRoles, &resp.Diagnostics)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgRoles := extractMapString(ctx, plan.OrgRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	clientAdmin := plan.ClientAdmin.ValueBool()
	user, err := r.client.UpdateUser(ctx, state.ID.ValueString(), client.UpdateUserRequest{
		FirstName:      plan.FirstName.ValueString(),
		LastName:       plan.LastName.ValueString(),
		ClientAdmin:    &clientAdmin,
		OrgDefaultRole: plan.OrgDefaultRole.ValueString(),
		OrgRoles:       orgRoles,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	configOrgRoles := plan.OrgRoles
	mapUserToState(user, &plan)
	plan.OrgRoles = configOrgRoles
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapUserToState(user *client.User, state *UserResourceModel) {
	state.ID = types.StringValue(user.ID)
	state.Email = types.StringValue(user.Email)
	state.FirstName = types.StringValue(user.FirstName)
	state.LastName = types.StringValue(user.LastName)
	state.ClientAdmin = types.BoolValue(user.ClientAdmin)
	state.OrgDefaultRole = types.StringValue(user.OrgDefaultRole)
	state.SSOOnly = types.BoolValue(user.SSOOnly)
	state.MFAEnabled = types.BoolValue(user.MFAEnabled)
	state.CreatedAt = types.Int64Value(user.CreatedAt)
	state.UpdatedAt = types.Int64Value(user.UpdatedAt)

	if user.OrgRoles != nil {
		elements := make(map[string]attr.Value, len(user.OrgRoles))
		for k, v := range user.OrgRoles {
			elements[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		state.OrgRoles = types.MapValueMust(types.StringType, elements)
	} else {
		state.OrgRoles = types.MapNull(types.StringType)
	}
}
