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
	_ datasource.DataSource              = &APIKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &APIKeyDataSource{}
)

type APIKeyDataSource struct {
	client *client.Client
}

type APIKeyDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	ClientID       types.String `tfsdk:"client_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Comment        types.String `tfsdk:"comment"`
	LastUsedAt     types.Int64  `tfsdk:"last_used_at"`
	LastUsedIP     types.String `tfsdk:"last_used_ip"`
	LastUsedUA     types.String `tfsdk:"last_used_ua"`
	Counter        types.Int64  `tfsdk:"counter"`
	UsageToday     types.Int64  `tfsdk:"usage_today"`
	UsageLimit     types.Int64  `tfsdk:"usage_limit"`
	Token          types.String `tfsdk:"token"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	Type           types.String `tfsdk:"type"`
}

func NewAPIKeyDataSource() datasource.DataSource {
	return &APIKeyDataSource{}
}

func (d *APIKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (d *APIKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero API key by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the API key.",
				Required:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client UUID.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization UUID the key is scoped to.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the API key was created.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Email of the user who created the key.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the API key.",
				Computed:    true,
			},
			"last_used_at": schema.Int64Attribute{
				Description: "Unix timestamp when the key was last used.",
				Computed:    true,
			},
			"last_used_ip": schema.StringAttribute{
				Description: "IP address from which the key was last used.",
				Computed:    true,
			},
			"last_used_ua": schema.StringAttribute{
				Description: "User agent string from the last use.",
				Computed:    true,
			},
			"counter": schema.Int64Attribute{
				Description: "Total usage counter.",
				Computed:    true,
			},
			"usage_today": schema.Int64Attribute{
				Description: "Usage count for today.",
				Computed:    true,
			},
			"usage_limit": schema.Int64Attribute{
				Description: "Maximum allowed usage.",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The API key token value.",
				Computed:    true,
				Sensitive:   true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the API key is inactive.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the API key.",
				Computed:    true,
			},
		},
	}
}

func (d *APIKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *APIKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state APIKeyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := d.client.GetAPIKey(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading API key", err.Error())
		return
	}

	state.ID = types.StringValue(apiKey.ID)
	state.ClientID = types.StringValue(apiKey.ClientID)
	state.OrganizationID = types.StringValue(apiKey.OrganizationID)
	state.CreatedAt = types.Int64Value(apiKey.CreatedAt)
	state.CreatedBy = types.StringValue(apiKey.CreatedBy)
	state.Comment = types.StringValue(apiKey.Comment)
	state.LastUsedAt = types.Int64Value(apiKey.LastUsedAt)
	state.LastUsedIP = types.StringValue(apiKey.LastUsedIP)
	state.LastUsedUA = types.StringValue(apiKey.LastUsedUA)
	state.Counter = types.Int64Value(apiKey.Counter)
	state.UsageToday = types.Int64Value(apiKey.UsageToday)
	state.UsageLimit = types.Int64Value(apiKey.UsageLimit)
	state.Token = types.StringValue(apiKey.Token)
	state.Inactive = types.BoolValue(apiKey.Inactive)
	state.Type = types.StringValue(apiKey.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
