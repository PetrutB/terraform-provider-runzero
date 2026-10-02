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
	_ datasource.DataSource              = &AccountTasksDataSource{}
	_ datasource.DataSourceWithConfigure = &AccountTasksDataSource{}
)

type AccountTasksDataSource struct {
	client *client.Client
}

type AccountTasksDataSourceModel struct {
	Search types.String  `tfsdk:"search"`
	Tasks  []TaskDSModel `tfsdk:"tasks"`
}

func NewAccountTasksDataSource() datasource.DataSource {
	return &AccountTasksDataSource{}
}

func (d *AccountTasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_tasks"
}

func (d *AccountTasksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all tasks across the RunZero account.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Description: "Optional search query to filter tasks.",
				Optional:    true,
			},
			"tasks": schema.ListNestedAttribute{
				Description: "List of tasks across all organizations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the task.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the task.",
							Computed:    true,
						},
						"organization_id": schema.StringAttribute{
							Description: "The organization ID.",
							Computed:    true,
						},
						"site_id": schema.StringAttribute{
							Description: "The site ID the task belongs to.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the task (e.g. scan).",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the task.",
							Computed:    true,
						},
						"created_at": schema.Int64Attribute{
							Description: "Creation timestamp (Unix epoch).",
							Computed:    true,
						},
						"updated_at": schema.Int64Attribute{
							Description: "Last update timestamp (Unix epoch).",
							Computed:    true,
						},
						"hidden": schema.BoolAttribute{
							Description: "Whether the task is hidden.",
							Computed:    true,
						},
						"recur": schema.BoolAttribute{
							Description: "Whether the task recurs.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AccountTasksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AccountTasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AccountTasksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tasks, err := d.client.ListAccountTasks(ctx, state.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading account tasks", err.Error())
		return
	}

	state.Tasks = make([]TaskDSModel, len(tasks))
	for i, t := range tasks {
		state.Tasks[i] = TaskDSModel{
			ID:             types.StringValue(t.ID),
			Name:           types.StringValue(t.Name),
			OrganizationID: types.StringValue(t.OrganizationID),
			SiteID:         types.StringValue(t.SiteID),
			Type:           types.StringValue(t.Type),
			Status:         types.StringValue(t.Status),
			CreatedAt:      types.Int64Value(t.CreatedAt),
			UpdatedAt:      types.Int64Value(t.UpdatedAt),
			Hidden:         types.BoolValue(t.Hidden),
			Recur:          types.BoolValue(t.Recur),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
