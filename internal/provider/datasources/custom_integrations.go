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
	_ datasource.DataSource              = &CustomIntegrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &CustomIntegrationsDataSource{}
)

type CustomIntegrationsDataSource struct {
	client *client.Client
}

type CustomIntegrationsDataSourceModel struct {
	CustomIntegrations []CustomIntegrationDSModel `tfsdk:"custom_integrations"`
}

type CustomIntegrationDSModel struct {
	ID          types.String `tfsdk:"id"`
	ClientID    types.String `tfsdk:"client_id"`
	CreatedByID types.String `tfsdk:"created_by_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	DestroyedAt types.String `tfsdk:"destroyed_at"`
	Name        types.String `tfsdk:"name"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}

func NewCustomIntegrationsDataSource() datasource.DataSource {
	return &CustomIntegrationsDataSource{}
}

func (d *CustomIntegrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_integrations"
}

func (d *CustomIntegrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero custom integrations.",
		Attributes: map[string]schema.Attribute{
			"custom_integrations": schema.ListNestedAttribute{
				Description: "List of custom integrations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the custom integration.",
							Computed:    true,
						},
						"client_id": schema.StringAttribute{
							Description: "The client UUID.",
							Computed:    true,
						},
						"created_by_id": schema.StringAttribute{
							Description: "The UUID of the user who created the custom integration.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Timestamp when the custom integration was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Timestamp when the custom integration was last updated.",
							Computed:    true,
						},
						"destroyed_at": schema.StringAttribute{
							Description: "Timestamp when the custom integration was destroyed.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the custom integration.",
							Computed:    true,
						},
						"icon": schema.StringAttribute{
							Description: "The icon of the custom integration.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the custom integration.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CustomIntegrationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CustomIntegrationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrations, err := d.client.ListCustomIntegrations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading custom integrations", err.Error())
		return
	}

	state.CustomIntegrations = make([]CustomIntegrationDSModel, len(integrations))
	for i, k := range integrations {
		state.CustomIntegrations[i] = CustomIntegrationDSModel{
			ID:          types.StringValue(k.ID),
			ClientID:    types.StringValue(k.ClientID),
			CreatedByID: types.StringValue(k.CreatedByID),
			CreatedAt:   types.StringValue(k.CreatedAt),
			UpdatedAt:   types.StringValue(k.UpdatedAt),
			DestroyedAt: types.StringValue(k.DestroyedAt),
			Name:        types.StringValue(k.Name),
			Icon:        types.StringValue(k.Icon),
			Description: types.StringValue(k.Description),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
