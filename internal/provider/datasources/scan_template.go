package datasources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &ScanTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &ScanTemplateDataSource{}
)

type ScanTemplateDataSource struct {
	client *client.Client
}

type ScanTemplateDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	ClientID            types.String `tfsdk:"client_id"`
	OrganizationID      types.String `tfsdk:"organization_id"`
	AgentID             types.String `tfsdk:"agent_id"`
	SiteID              types.String `tfsdk:"site_id"`
	CreatedAt           types.Int64  `tfsdk:"created_at"`
	CreatedBy           types.String `tfsdk:"created_by"`
	CreatedByUserID     types.String `tfsdk:"created_by_user_id"`
	UpdatedAt           types.Int64  `tfsdk:"updated_at"`
	Type                types.String `tfsdk:"type"`
	Status              types.String `tfsdk:"status"`
	Hidden              types.Bool   `tfsdk:"hidden"`
	Recur               types.Bool   `tfsdk:"recur"`
	RecurFrequency      types.String `tfsdk:"recur_frequency"`
	StartTime           types.Int64  `tfsdk:"start_time"`
	RecurLast           types.Int64  `tfsdk:"recur_last"`
	RecurNext           types.Int64  `tfsdk:"recur_next"`
	RecurLastTaskID     types.String `tfsdk:"recur_last_task_id"`
	GracePeriod         types.String `tfsdk:"grace_period"`
	Global              types.Bool   `tfsdk:"global"`
	ParentID            types.String `tfsdk:"parent_id"`
	TemplateID          types.String `tfsdk:"template_id"`
	HostedZoneID        types.String `tfsdk:"hosted_zone_id"`
	CustomIntegrationID types.String `tfsdk:"custom_integration_id"`
	SourceID            types.String `tfsdk:"source_id"`
	Params              types.Map    `tfsdk:"params"`
}

func NewScanTemplateDataSource() datasource.DataSource {
	return &ScanTemplateDataSource{}
}

func (d *ScanTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scan_template"
}

func (d *ScanTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := scanTemplateDSSchemaAttributes()
	attrs["id"] = schema.StringAttribute{
		Description: "The UUID of the scan template.",
		Required:    true,
	}
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero scan template by ID.",
		Attributes:  attrs,
	}
}

func (d *ScanTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScanTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ScanTemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := d.client.GetScanTemplate(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading scan template", err.Error())
		return
	}

	mapScanTemplateToDS(ctx, template, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func scanTemplateDSSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The UUID of the scan template.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the scan template.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Description of the scan template.",
			Computed:    true,
		},
		"client_id": schema.StringAttribute{
			Description: "The client UUID.",
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization UUID.",
			Computed:    true,
		},
		"agent_id": schema.StringAttribute{
			Description: "The agent UUID.",
			Computed:    true,
		},
		"site_id": schema.StringAttribute{
			Description: "The site UUID.",
			Computed:    true,
		},
		"created_at": schema.Int64Attribute{
			Description: "Unix timestamp when the scan template was created.",
			Computed:    true,
		},
		"created_by": schema.StringAttribute{
			Description: "Email of the user who created the template.",
			Computed:    true,
		},
		"created_by_user_id": schema.StringAttribute{
			Description: "UUID of the user who created the template.",
			Computed:    true,
		},
		"updated_at": schema.Int64Attribute{
			Description: "Unix timestamp when the scan template was last updated.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the scan template.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the scan template.",
			Computed:    true,
		},
		"hidden": schema.BoolAttribute{
			Description: "Whether the scan template is hidden.",
			Computed:    true,
		},
		"recur": schema.BoolAttribute{
			Description: "Whether the scan template is recurring.",
			Computed:    true,
		},
		"recur_frequency": schema.StringAttribute{
			Description: "The recurrence frequency.",
			Computed:    true,
		},
		"start_time": schema.Int64Attribute{
			Description: "Unix timestamp for the scan start time.",
			Computed:    true,
		},
		"recur_last": schema.Int64Attribute{
			Description: "Unix timestamp of the last recurrence.",
			Computed:    true,
		},
		"recur_next": schema.Int64Attribute{
			Description: "Unix timestamp of the next recurrence.",
			Computed:    true,
		},
		"recur_last_task_id": schema.StringAttribute{
			Description: "UUID of the last recurring task.",
			Computed:    true,
		},
		"grace_period": schema.StringAttribute{
			Description: "Grace period for the scan template.",
			Computed:    true,
		},
		"global": schema.BoolAttribute{
			Description: "Whether the scan template is globally accessible.",
			Computed:    true,
		},
		"parent_id": schema.StringAttribute{
			Description: "The parent template UUID.",
			Computed:    true,
		},
		"template_id": schema.StringAttribute{
			Description: "The base template UUID.",
			Computed:    true,
		},
		"hosted_zone_id": schema.StringAttribute{
			Description: "The hosted zone UUID.",
			Computed:    true,
		},
		"custom_integration_id": schema.StringAttribute{
			Description: "The custom integration UUID.",
			Computed:    true,
		},
		"source_id": schema.StringAttribute{
			Description: "The source UUID.",
			Computed:    true,
		},
		"params": schema.MapAttribute{
			Description: "Scan template parameters.",
			Computed:    true,
			ElementType: types.StringType,
		},
	}
}

func mapScanTemplateToDS(ctx context.Context, t *client.ScanTemplate, m *ScanTemplateDataSourceModel, diags *diag.Diagnostics) {
	m.ID = types.StringValue(t.ID)
	m.Name = types.StringValue(t.Name)
	m.Description = types.StringValue(t.Description)
	m.ClientID = types.StringValue(t.ClientID)
	m.OrganizationID = types.StringValue(t.OrganizationID)
	m.AgentID = types.StringValue(t.AgentID)
	m.SiteID = types.StringValue(t.SiteID)
	m.CreatedAt = types.Int64Value(t.CreatedAt)
	m.CreatedBy = types.StringValue(t.CreatedBy)
	m.CreatedByUserID = types.StringValue(t.CreatedByUserID)
	m.UpdatedAt = types.Int64Value(t.UpdatedAt)
	m.Type = types.StringValue(t.Type)
	m.Status = types.StringValue(t.Status)
	m.Hidden = types.BoolValue(t.Hidden)
	m.Recur = types.BoolValue(t.Recur)
	m.RecurFrequency = types.StringValue(t.RecurFrequency)
	m.StartTime = types.Int64Value(t.StartTime)
	m.RecurLast = types.Int64Value(t.RecurLast)
	m.RecurNext = types.Int64Value(t.RecurNext)
	m.RecurLastTaskID = types.StringValue(t.RecurLastTaskID)
	m.GracePeriod = types.StringValue(strconv.FormatInt(t.GracePeriod, 10))
	m.Global = types.BoolValue(t.Global)
	m.ParentID = types.StringValue(t.ParentID)
	m.TemplateID = types.StringValue(t.TemplateID)
	m.HostedZoneID = types.StringValue(t.HostedZoneID)
	m.CustomIntegrationID = types.StringValue(t.CustomIntegrationID)
	m.SourceID = types.StringValue(t.SourceID)

	params, d := types.MapValueFrom(ctx, types.StringType, t.Params)
	diags.Append(d...)
	m.Params = params
}
