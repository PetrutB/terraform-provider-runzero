package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource                = &UserInviteResource{}
	_ resource.ResourceWithImportState = &UserInviteResource{}
	_ resource.ResourceWithConfigure   = &UserInviteResource{}
)

type UserInviteResource struct {
	client *client.Client
}

type UserInviteResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ClientID       types.String `tfsdk:"client_id"`
	Name           types.String `tfsdk:"name"`
	FirstName      types.String `tfsdk:"first_name"`
	LastName       types.String `tfsdk:"last_name"`
	Email          types.String `tfsdk:"email"`
	ClientAdmin    types.Bool   `tfsdk:"client_admin"`
	OrgDefaultRole types.String `tfsdk:"org_default_role"`
	OrgRoles       types.Map    `tfsdk:"org_roles"`
	Subject        types.String `tfsdk:"subject"`
	Message        types.String `tfsdk:"message"`
	SSOOnly        types.Bool   `tfsdk:"sso_only"`
	MFAEnabled     types.Bool   `tfsdk:"mfa_enabled"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewUserInviteResource() resource.Resource {
	return &UserInviteResource{}
}

func (r *UserInviteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_invite"
}

func (r *UserInviteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Invites a user to RunZero via email. This is a create-only resource that sends an invitation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the invited user.",
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
				Description: "The display name of the user.",
				Computed:    true,
			},
			"first_name": schema.StringAttribute{
				Description: "The first name of the user to invite.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_name": schema.StringAttribute{
				Description: "The last name of the user to invite.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The email address of the user to invite.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_admin": schema.BoolAttribute{
				Description: "Whether the user should be a client administrator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"org_default_role": schema.StringAttribute{
				Description: "The default organization role for the user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_roles": schema.MapAttribute{
				Description: "Map of organization UUIDs to role names.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"subject": schema.StringAttribute{
				Description: "Custom subject line for the invite email. Write-only; not returned by the API.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"message": schema.StringAttribute{
				Description: "Custom message body for the invite email. Write-only; not returned by the API.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

func (r *UserInviteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserInviteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserInviteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgRoles := extractMapString(ctx, plan.OrgRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.InviteUser(ctx, client.UserInviteRequest{
		FirstName:      plan.FirstName.ValueString(),
		LastName:       plan.LastName.ValueString(),
		Email:          plan.Email.ValueString(),
		ClientAdmin:    plan.ClientAdmin.ValueBool(),
		OrgDefaultRole: plan.OrgDefaultRole.ValueString(),
		OrgRoles:       orgRoles,
		Subject:        plan.Subject.ValueString(),
		Message:        plan.Message.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error inviting user", err.Error())
		return
	}

	mapInviteUserToState(user, &plan)
	tflog.Trace(ctx, "invited user", map[string]interface{}{"id": user.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserInviteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserInviteResourceModel
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
		resp.Diagnostics.AddError("Error reading invited user", err.Error())
		return
	}

	mapInviteUserToState(user, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *UserInviteResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"User invitations cannot be updated. Change attributes to trigger a replacement.",
	)
}

func (r *UserInviteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserInviteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting invited user", err.Error())
		return
	}
}

func (r *UserInviteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapInviteUserToState(user *client.User, state *UserInviteResourceModel) {
	state.ID = types.StringValue(user.ID)
	state.ClientID = types.StringValue(user.ClientID)
	state.Name = types.StringValue(user.Name)
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
			elements[k] = types.StringValue(v)
		}
		state.OrgRoles = types.MapValueMust(types.StringType, elements)
	} else {
		state.OrgRoles = types.MapNull(types.StringType)
	}
}
