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
	_ datasource.DataSource              = &CustomIntegrationDataSource{}
	_ datasource.DataSourceWithConfigure = &CustomIntegrationDataSource{}
)

type CustomIntegrationDataSource struct {
	client *client.Client
}

type CustomIntegrationDataSourceModel struct {
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

func NewCustomIntegrationDataSource() datasource.DataSource {
	return &CustomIntegrationDataSource{}
}

func (d *CustomIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_integration"
}

func (d *CustomIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero custom integration by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the custom integration.",
				Required:    true,
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
	}
}

func (d *CustomIntegrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CustomIntegrationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := d.client.GetCustomIntegration(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading custom integration", err.Error())
		return
	}

	state.ID = types.StringValue(integration.ID)
	state.ClientID = types.StringValue(integration.ClientID)
	state.CreatedByID = types.StringValue(integration.CreatedByID)
	state.CreatedAt = types.StringValue(integration.CreatedAt)
	state.UpdatedAt = types.StringValue(integration.UpdatedAt)
	state.DestroyedAt = types.StringValue(integration.DestroyedAt)
	state.Name = types.StringValue(integration.Name)
	state.Icon = types.StringValue(integration.Icon)
	state.Description = types.StringValue(integration.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
