package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &OrgAPIKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &OrgAPIKeyDataSource{}
)

type OrgAPIKeyDataSource struct {
	client *client.Client
}

type OrgAPIKeyDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	ClientID       types.String `tfsdk:"client_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Comment        types.String `tfsdk:"comment"`
	Token          types.String `tfsdk:"token"`
	Type           types.String `tfsdk:"type"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	Counter        types.Int64  `tfsdk:"counter"`
	UsageToday     types.Int64  `tfsdk:"usage_today"`
	UsageLimit     types.Int64  `tfsdk:"usage_limit"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	LastUsedAt     types.Int64  `tfsdk:"last_used_at"`
	LastUsedIP     types.String `tfsdk:"last_used_ip"`
	LastUsedUA     types.String `tfsdk:"last_used_ua"`
}

func NewOrgAPIKeyDataSource() datasource.DataSource {
	return &OrgAPIKeyDataSource{}
}

func (d *OrgAPIKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_api_key"
}

func (d *OrgAPIKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the API key for the organization associated with the configured org token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the organization API key.",
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client UUID that owns the key.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization UUID the key belongs to.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment associated with the key.",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The organization API key token.",
				Computed:    true,
				Sensitive:   true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the key.",
				Computed:    true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the key is inactive.",
				Computed:    true,
			},
			"counter": schema.Int64Attribute{
				Description: "Total usage counter for the key.",
				Computed:    true,
			},
			"usage_today": schema.Int64Attribute{
				Description: "Number of times the key was used today.",
				Computed:    true,
			},
			"usage_limit": schema.Int64Attribute{
				Description: "Daily usage limit for the key.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the key was created.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Identifier of who created the key.",
				Computed:    true,
			},
			"last_used_at": schema.Int64Attribute{
				Description: "Unix timestamp when the key was last used.",
				Computed:    true,
			},
			"last_used_ip": schema.StringAttribute{
				Description: "IP address that last used the key.",
				Computed:    true,
			},
			"last_used_ua": schema.StringAttribute{
				Description: "User agent that last used the key.",
				Computed:    true,
			},
		},
	}
}

func (d *OrgAPIKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrgAPIKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OrgAPIKeyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := d.client.GetOrgKey(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading org API key", err.Error())
		return
	}

	state.ID = types.StringValue(key.ID)
	state.ClientID = types.StringValue(key.ClientID)
	state.OrganizationID = types.StringValue(key.OrganizationID)
	state.Comment = types.StringValue(key.Comment)
	state.Token = types.StringValue(key.Token)
	state.Type = types.StringValue(key.Type)
	state.Inactive = types.BoolValue(key.Inactive)
	state.Counter = types.Int64Value(key.Counter)
	state.UsageToday = types.Int64Value(key.UsageToday)
	state.UsageLimit = types.Int64Value(key.UsageLimit)
	state.CreatedAt = types.Int64Value(key.CreatedAt)
	state.CreatedBy = types.StringValue(key.CreatedBy)
	state.LastUsedAt = types.Int64Value(key.LastUsedAt)
	state.LastUsedIP = types.StringValue(key.LastUsedIP)
	state.LastUsedUA = types.StringValue(key.LastUsedUA)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
