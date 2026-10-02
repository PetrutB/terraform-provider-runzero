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
	_ datasource.DataSource              = &TaskDataSource{}
	_ datasource.DataSourceWithConfigure = &TaskDataSource{}
)

type TaskDataSource struct {
	client *client.Client
}

type TaskDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	OrganizationID types.String `tfsdk:"organization_id"`
	AgentID        types.String `tfsdk:"agent_id"`
	SiteID         types.String `tfsdk:"site_id"`
	HostedZoneID   types.String `tfsdk:"hosted_zone_id"`
	TemplateID     types.String `tfsdk:"template_id"`
	Type           types.String `tfsdk:"type"`
	Status         types.String `tfsdk:"status"`
	Error          types.String `tfsdk:"error"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
	Hidden         types.Bool   `tfsdk:"hidden"`
	Recur          types.Bool   `tfsdk:"recur"`
	RecurFrequency types.String `tfsdk:"recur_frequency"`
	StartTime      types.Int64  `tfsdk:"start_time"`
	RecurLast      types.Int64  `tfsdk:"recur_last"`
	RecurNext      types.Int64  `tfsdk:"recur_next"`
}

func NewTaskDataSource() datasource.DataSource {
	return &TaskDataSource{}
}

func (d *TaskDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

func (d *TaskDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero task by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the task.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the task.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the task.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID.",
				Computed:    true,
			},
			"agent_id": schema.StringAttribute{
				Description: "The agent ID executing the task.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID the task belongs to.",
				Computed:    true,
			},
			"hosted_zone_id": schema.StringAttribute{
				Description: "The hosted zone ID executing the task.",
				Computed:    true,
			},
			"template_id": schema.StringAttribute{
				Description: "The scan template ID.",
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
			"error": schema.StringAttribute{
				Description: "Error message if the task failed.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Creation timestamp (Unix epoch).",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user who created the task.",
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
			"recur_frequency": schema.StringAttribute{
				Description: "Recurrence frequency (e.g. hourly).",
				Computed:    true,
			},
			"start_time": schema.Int64Attribute{
				Description: "Start time timestamp (Unix epoch).",
				Computed:    true,
			},
			"recur_last": schema.Int64Attribute{
				Description: "Last recurrence timestamp (Unix epoch).",
				Computed:    true,
			},
			"recur_next": schema.Int64Attribute{
				Description: "Next recurrence timestamp (Unix epoch).",
				Computed:    true,
			},
		},
	}
}

func (d *TaskDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TaskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TaskDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	task, err := d.client.GetTask(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading task", err.Error())
		return
	}

	state.ID = types.StringValue(task.ID)
	state.Name = types.StringValue(task.Name)
	state.Description = types.StringValue(task.Description)
	state.OrganizationID = types.StringValue(task.OrganizationID)
	state.AgentID = types.StringValue(task.AgentID)
	state.SiteID = types.StringValue(task.SiteID)
	state.HostedZoneID = types.StringValue(task.HostedZoneID)
	state.TemplateID = types.StringValue(task.TemplateID)
	state.Type = types.StringValue(task.Type)
	state.Status = types.StringValue(task.Status)
	state.Error = types.StringValue(task.Error)
	state.CreatedAt = types.Int64Value(task.CreatedAt)
	state.CreatedBy = types.StringValue(task.CreatedBy)
	state.UpdatedAt = types.Int64Value(task.UpdatedAt)
	state.Hidden = types.BoolValue(task.Hidden)
	state.Recur = types.BoolValue(task.Recur)
	state.RecurFrequency = types.StringValue(task.RecurFrequency)
	state.StartTime = types.Int64Value(task.StartTime)
	state.RecurLast = types.Int64Value(task.RecurLast)
	state.RecurNext = types.Int64Value(task.RecurNext)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
