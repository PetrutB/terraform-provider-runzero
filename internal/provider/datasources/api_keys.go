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
	_ datasource.DataSource              = &APIKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &APIKeysDataSource{}
)

type APIKeysDataSource struct {
	client *client.Client
}

type APIKeysDataSourceModel struct {
	APIKeys []APIKeyDSModel `tfsdk:"api_keys"`
}

type APIKeyDSModel struct {
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

func NewAPIKeysDataSource() datasource.DataSource {
	return &APIKeysDataSource{}
}

func (d *APIKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_keys"
}

func (d *APIKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero API keys.",
		Attributes: map[string]schema.Attribute{
			"api_keys": schema.ListNestedAttribute{
				Description: "List of API keys.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the API key.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *APIKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *APIKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state APIKeysDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKeys, err := d.client.ListAPIKeys(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading API keys", err.Error())
		return
	}

	state.APIKeys = make([]APIKeyDSModel, len(apiKeys))
	for i, k := range apiKeys {
		state.APIKeys[i] = APIKeyDSModel{
			ID:             types.StringValue(k.ID),
			ClientID:       types.StringValue(k.ClientID),
			OrganizationID: types.StringValue(k.OrganizationID),
			CreatedAt:      types.Int64Value(k.CreatedAt),
			CreatedBy:      types.StringValue(k.CreatedBy),
			Comment:        types.StringValue(k.Comment),
			LastUsedAt:     types.Int64Value(k.LastUsedAt),
			LastUsedIP:     types.StringValue(k.LastUsedIP),
			LastUsedUA:     types.StringValue(k.LastUsedUA),
			Counter:        types.Int64Value(k.Counter),
			UsageToday:     types.Int64Value(k.UsageToday),
			UsageLimit:     types.Int64Value(k.UsageLimit),
			Token:          types.StringValue(k.Token),
			Inactive:       types.BoolValue(k.Inactive),
			Type:           types.StringValue(k.Type),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
