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
	_ datasource.DataSource              = &TasksDataSource{}
	_ datasource.DataSourceWithConfigure = &TasksDataSource{}
)

type TasksDataSource struct {
	client *client.Client
}

type TasksDataSourceModel struct {
	Status types.String  `tfsdk:"status"`
	Tasks  []TaskDSModel `tfsdk:"tasks"`
}

type TaskDSModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OrganizationID types.String `tfsdk:"organization_id"`
	SiteID         types.String `tfsdk:"site_id"`
	Type           types.String `tfsdk:"type"`
	Status         types.String `tfsdk:"status"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
	Hidden         types.Bool   `tfsdk:"hidden"`
	Recur          types.Bool   `tfsdk:"recur"`
}

func NewTasksDataSource() datasource.DataSource {
	return &TasksDataSource{}
}

func (d *TasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tasks"
}

func (d *TasksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero tasks.",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Description: "Optional status filter for tasks.",
				Optional:    true,
			},
			"tasks": schema.ListNestedAttribute{
				Description: "List of tasks.",
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

func (d *TasksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TasksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tasks, err := d.client.ListTasks(ctx, state.Status.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading tasks", err.Error())
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
