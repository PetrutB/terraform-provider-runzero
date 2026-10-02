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
	_ datasource.DataSource              = &HealthDataSource{}
	_ datasource.DataSourceWithConfigure = &HealthDataSource{}
)

type HealthDataSource struct {
	client *client.Client
}

type HealthDataSourceModel struct {
	Status  types.String `tfsdk:"status"`
	Version types.String `tfsdk:"version"`
}

func NewHealthDataSource() datasource.DataSource {
	return &HealthDataSource{}
}

func (d *HealthDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_health"
}

func (d *HealthDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the RunZero API health check status and version.",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Description: "Health check status (e.g. 'ok').",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "RunZero API version.",
				Computed:    true,
			},
		},
	}
}

func (d *HealthDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HealthDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	health, err := d.client.GetHealth(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading health check", err.Error())
		return
	}

	state := HealthDataSourceModel{
		Status:  types.StringValue(health.Status),
		Version: types.StringValue(health.Version),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
