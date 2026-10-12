package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

type UserDataSource struct {
	client *client.Client
}

type UserDataSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	ClientID              types.String `tfsdk:"client_id"`
	Name                  types.String `tfsdk:"name"`
	Email                 types.String `tfsdk:"email"`
	FirstName             types.String `tfsdk:"first_name"`
	LastName              types.String `tfsdk:"last_name"`
	ClientAdmin           types.Bool   `tfsdk:"client_admin"`
	OrgDefaultRole        types.String `tfsdk:"org_default_role"`
	OrgRoles              types.Map    `tfsdk:"org_roles"`
	SSOOnly               types.Bool   `tfsdk:"sso_only"`
	MFAEnabled            types.Bool   `tfsdk:"mfa_enabled"`
	LoginFailures         types.Int64  `tfsdk:"login_failures"`
	LastLoginIP           types.String `tfsdk:"last_login_ip"`
	LastLoginAt           types.Int64  `tfsdk:"last_login_at"`
	LastLoginUA           types.String `tfsdk:"last_login_ua"`
	LastActivityAt        types.Int64  `tfsdk:"last_activity_at"`
	Actions               types.Int64  `tfsdk:"actions"`
	LastActionAt          types.Int64  `tfsdk:"last_action_at"`
	InviteTokenExpiration types.Int64  `tfsdk:"invite_token_expiration"`
	ResetTokenExpiration  types.Int64  `tfsdk:"reset_token_expiration"`
	PasswordEnabledAt     types.Int64  `tfsdk:"password_enabled_at"`
	TermsAgreedAsOf       types.Int64  `tfsdk:"terms_agreed_as_of"`
	CreatedAt             types.Int64  `tfsdk:"created_at"`
	UpdatedAt             types.Int64  `tfsdk:"updated_at"`
}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := userDSSchemaAttributes()
	attrs["id"] = schema.StringAttribute{
		Description: "The UUID of the user.",
		Required:    true,
	}
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero user by ID.",
		Attributes:  attrs,
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	mapUserToDS(ctx, user, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func userDSSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The UUID of the user.",
			Computed:    true,
		},
		"client_id": schema.StringAttribute{
			Description: "The client UUID.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The display name of the user.",
			Computed:    true,
		},
		"email": schema.StringAttribute{
			Description: "The email address of the user.",
			Computed:    true,
		},
		"first_name": schema.StringAttribute{
			Description: "The first name of the user.",
			Computed:    true,
		},
		"last_name": schema.StringAttribute{
			Description: "The last name of the user.",
			Computed:    true,
		},
		"client_admin": schema.BoolAttribute{
			Description: "Whether the user is a client administrator.",
			Computed:    true,
		},
		"org_default_role": schema.StringAttribute{
			Description: "The default organization role for the user.",
			Computed:    true,
		},
		"org_roles": schema.MapAttribute{
			Description: "Map of organization UUIDs to role names.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"sso_only": schema.BoolAttribute{
			Description: "Whether the user can only authenticate via SSO.",
			Computed:    true,
		},
		"mfa_enabled": schema.BoolAttribute{
			Description: "Whether multi-factor authentication is enabled.",
			Computed:    true,
		},
		"login_failures": schema.Int64Attribute{
			Description: "Number of consecutive login failures.",
			Computed:    true,
		},
		"last_login_ip": schema.StringAttribute{
			Description: "IP address of the last login.",
			Computed:    true,
		},
		"last_login_at": schema.Int64Attribute{
			Description: "Unix timestamp of the last login.",
			Computed:    true,
		},
		"last_login_ua": schema.StringAttribute{
			Description: "User agent string of the last login.",
			Computed:    true,
		},
		"last_activity_at": schema.Int64Attribute{
			Description: "Unix timestamp of the last activity.",
			Computed:    true,
		},
		"actions": schema.Int64Attribute{
			Description: "Total number of actions performed.",
			Computed:    true,
		},
		"last_action_at": schema.Int64Attribute{
			Description: "Unix timestamp of the last action.",
			Computed:    true,
		},
		"invite_token_expiration": schema.Int64Attribute{
			Description: "Unix timestamp when the invite token expires.",
			Computed:    true,
		},
		"reset_token_expiration": schema.Int64Attribute{
			Description: "Unix timestamp when the password reset token expires.",
			Computed:    true,
		},
		"password_enabled_at": schema.Int64Attribute{
			Description: "Unix timestamp when password authentication was enabled.",
			Computed:    true,
		},
		"terms_agreed_as_of": schema.Int64Attribute{
			Description: "Unix timestamp when the user agreed to the terms.",
			Computed:    true,
		},
		"created_at": schema.Int64Attribute{
			Description: "Unix timestamp when the user was created.",
			Computed:    true,
		},
		"updated_at": schema.Int64Attribute{
			Description: "Unix timestamp when the user was last updated.",
			Computed:    true,
		},
	}
}

func mapUserToDS(ctx context.Context, u *client.User, m *UserDataSourceModel, diags *diag.Diagnostics) {
	m.ID = types.StringValue(u.ID)
	m.ClientID = types.StringValue(u.ClientID)
	m.Name = types.StringValue(u.Name)
	m.Email = types.StringValue(u.Email)
	m.FirstName = types.StringValue(u.FirstName)
	m.LastName = types.StringValue(u.LastName)
	m.ClientAdmin = types.BoolValue(u.ClientAdmin)
	m.OrgDefaultRole = types.StringValue(u.OrgDefaultRole)
	m.SSOOnly = types.BoolValue(u.SSOOnly)
	m.MFAEnabled = types.BoolValue(u.MFAEnabled)
	m.LoginFailures = types.Int64Value(u.LoginFailures)
	m.LastLoginIP = types.StringValue(u.LastLoginIP)
	m.LastLoginAt = types.Int64Value(u.LastLoginAt)
	m.LastLoginUA = types.StringValue(u.LastLoginUA)
	m.LastActivityAt = types.Int64Value(u.LastActivityAt)
	m.Actions = types.Int64Value(u.Actions)
	m.LastActionAt = types.Int64Value(u.LastActionAt)
	m.InviteTokenExpiration = types.Int64Value(u.InviteTokenExpiration)
	m.ResetTokenExpiration = types.Int64Value(u.ResetTokenExpiration)
	m.PasswordEnabledAt = types.Int64Value(u.PasswordEnabledAt)
	m.TermsAgreedAsOf = types.Int64Value(u.TermsAgreedAsOf)
	m.CreatedAt = types.Int64Value(u.CreatedAt)
	m.UpdatedAt = types.Int64Value(u.UpdatedAt)

	orgRolesStr := make(map[string]string, len(u.OrgRoles))
	for k, v := range u.OrgRoles {
		orgRolesStr[k] = fmt.Sprintf("%v", v)
	}
	orgRoles, d := types.MapValueFrom(ctx, types.StringType, orgRolesStr)
	diags.Append(d...)
	m.OrgRoles = orgRoles
}
