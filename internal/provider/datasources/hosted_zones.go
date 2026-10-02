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
	_ datasource.DataSource              = &HostedZonesDataSource{}
	_ datasource.DataSourceWithConfigure = &HostedZonesDataSource{}
)

type HostedZonesDataSource struct {
	client *client.Client
}

type HostedZonesDataSourceModel struct {
	HostedZones []HostedZoneDSModel `tfsdk:"hosted_zones"`
}

type HostedZoneDSModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	ExplorersConcurrency types.Int64  `tfsdk:"explorers_concurrency"`
	ExplorersTotal       types.Int64  `tfsdk:"explorers_total"`
	TasksActive          types.Int64  `tfsdk:"tasks_active"`
	TasksWaiting         types.Int64  `tfsdk:"tasks_waiting"`
	OrganizationID       types.String `tfsdk:"organization_id"`
}

func NewHostedZonesDataSource() datasource.DataSource {
	return &HostedZonesDataSource{}
}

func (d *HostedZonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosted_zones"
}

func (d *HostedZonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero hosted zones.",
		Attributes: map[string]schema.Attribute{
			"hosted_zones": schema.ListNestedAttribute{
				Description: "List of hosted zones.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the hosted zone.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the hosted zone.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the hosted zone is enabled.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Last update timestamp (ISO 8601 date-time).",
							Computed:    true,
						},
						"explorers_concurrency": schema.Int64Attribute{
							Description: "Number of concurrent explorers.",
							Computed:    true,
						},
						"explorers_total": schema.Int64Attribute{
							Description: "Total number of explorers.",
							Computed:    true,
						},
						"tasks_active": schema.Int64Attribute{
							Description: "Number of active tasks.",
							Computed:    true,
						},
						"tasks_waiting": schema.Int64Attribute{
							Description: "Number of waiting tasks.",
							Computed:    true,
						},
						"organization_id": schema.StringAttribute{
							Description: "The organization ID.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *HostedZonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostedZonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state HostedZonesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zones, err := d.client.ListHostedZones(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading hosted zones", err.Error())
		return
	}

	state.HostedZones = make([]HostedZoneDSModel, len(zones))
	for i, z := range zones {
		state.HostedZones[i] = HostedZoneDSModel{
			ID:                   types.StringValue(z.ID),
			Name:                 types.StringValue(z.Name),
			Enabled:              types.BoolValue(z.Enabled),
			UpdatedAt:            types.StringValue(z.UpdatedAt),
			ExplorersConcurrency: types.Int64Value(z.ExplorersConcurrency),
			ExplorersTotal:       types.Int64Value(z.ExplorersTotal),
			TasksActive:          types.Int64Value(z.TasksActive),
			TasksWaiting:         types.Int64Value(z.TasksWaiting),
			OrganizationID:       types.StringValue(z.OrganizationID),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
