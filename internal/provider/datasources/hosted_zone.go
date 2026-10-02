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
	_ datasource.DataSource              = &HostedZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &HostedZoneDataSource{}
)

type HostedZoneDataSource struct {
	client *client.Client
}

type HostedZoneDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	ProcessorID          types.String `tfsdk:"processor_id"`
	ExplorersConcurrency types.Int64  `tfsdk:"explorers_concurrency"`
	ExplorersTotal       types.Int64  `tfsdk:"explorers_total"`
	TasksActive          types.Int64  `tfsdk:"tasks_active"`
	TasksWaiting         types.Int64  `tfsdk:"tasks_waiting"`
	OrganizationID       types.String `tfsdk:"organization_id"`
}

func NewHostedZoneDataSource() datasource.DataSource {
	return &HostedZoneDataSource{}
}

func (d *HostedZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosted_zone"
}

func (d *HostedZoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero hosted zone by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the hosted zone.",
				Required:    true,
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
			"processor_id": schema.StringAttribute{
				Description: "The processor ID.",
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
				Description: "The organization ID the hosted zone belongs to.",
				Computed:    true,
			},
		},
	}
}

func (d *HostedZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HostedZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state HostedZoneDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := d.client.GetHostedZone(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading hosted zone", err.Error())
		return
	}

	state.ID = types.StringValue(zone.ID)
	state.Name = types.StringValue(zone.Name)
	state.Enabled = types.BoolValue(zone.Enabled)
	state.UpdatedAt = types.StringValue(zone.UpdatedAt)
	state.ProcessorID = types.StringValue(zone.ProcessorID)
	state.ExplorersConcurrency = types.Int64Value(zone.ExplorersConcurrency)
	state.ExplorersTotal = types.Int64Value(zone.ExplorersTotal)
	state.TasksActive = types.Int64Value(zone.TasksActive)
	state.TasksWaiting = types.Int64Value(zone.TasksWaiting)
	state.OrganizationID = types.StringValue(zone.OrganizationID)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
