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
	_ datasource.DataSource              = &AgentDataSource{}
	_ datasource.DataSourceWithConfigure = &AgentDataSource{}
)

type AgentDataSource struct {
	client *client.Client
}

type AgentDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SiteID         types.String `tfsdk:"site_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	OS             types.String `tfsdk:"os"`
	Arch           types.String `tfsdk:"arch"`
	Version        types.String `tfsdk:"version"`
	ExternalIP     types.String `tfsdk:"external_ip"`
	InternalIP     types.String `tfsdk:"internal_ip"`
	Connected      types.Bool   `tfsdk:"connected"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	LastCheckin    types.Int64  `tfsdk:"last_checkin"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewAgentDataSource() datasource.DataSource {
	return &AgentDataSource{}
}

func (d *AgentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (d *AgentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero agent by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the agent.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the agent.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID the agent is associated with.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID the agent belongs to.",
				Computed:    true,
			},
			"os": schema.StringAttribute{
				Description: "The operating system of the agent host.",
				Computed:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the agent host.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of the agent.",
				Computed:    true,
			},
			"external_ip": schema.StringAttribute{
				Description: "The external IP address of the agent.",
				Computed:    true,
			},
			"internal_ip": schema.StringAttribute{
				Description: "The internal IP address of the agent.",
				Computed:    true,
			},
			"connected": schema.BoolAttribute{
				Description: "Whether the agent is currently connected.",
				Computed:    true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the agent is inactive.",
				Computed:    true,
			},
			"last_checkin": schema.Int64Attribute{
				Description: "Unix timestamp of the agent's last check-in.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the agent was created.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the agent was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *AgentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AgentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AgentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := d.client.GetAgent(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading agent", err.Error())
		return
	}

	state.ID = types.StringValue(agent.ID)
	state.Name = types.StringValue(agent.Name)
	state.SiteID = types.StringValue(agent.SiteID)
	state.OrganizationID = types.StringValue(agent.OrganizationID)
	state.OS = types.StringValue(agent.OS)
	state.Arch = types.StringValue(agent.Arch)
	state.Version = types.StringValue(agent.Version)
	state.ExternalIP = types.StringValue(agent.ExternalIP)
	state.InternalIP = types.StringValue(agent.InternalIP)
	state.Connected = types.BoolValue(agent.Connected)
	state.Inactive = types.BoolValue(agent.Inactive)
	state.LastCheckin = types.Int64Value(agent.LastCheckin)
	state.CreatedAt = types.Int64Value(agent.CreatedAt)
	state.UpdatedAt = types.Int64Value(agent.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
