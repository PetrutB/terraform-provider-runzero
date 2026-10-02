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
	_ datasource.DataSource              = &AgentsDataSource{}
	_ datasource.DataSourceWithConfigure = &AgentsDataSource{}
)

type AgentsDataSource struct {
	client *client.Client
}

type AgentsDataSourceModel struct {
	Agents []AgentDSModel `tfsdk:"agents"`
}

type AgentDSModel struct {
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

func NewAgentsDataSource() datasource.DataSource {
	return &AgentsDataSource{}
}

func (d *AgentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agents"
}

func (d *AgentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero agents.",
		Attributes: map[string]schema.Attribute{
			"agents": schema.ListNestedAttribute{
				Description: "List of agents.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the agent.",
							Computed:    true,
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
							Description: "The organization ID.",
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
				},
			},
		},
	}
}

func (d *AgentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AgentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AgentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agents, err := d.client.ListAgents(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading agents", err.Error())
		return
	}

	state.Agents = make([]AgentDSModel, len(agents))
	for i, a := range agents {
		state.Agents[i] = AgentDSModel{
			ID:             types.StringValue(a.ID),
			Name:           types.StringValue(a.Name),
			SiteID:         types.StringValue(a.SiteID),
			OrganizationID: types.StringValue(a.OrganizationID),
			OS:             types.StringValue(a.OS),
			Arch:           types.StringValue(a.Arch),
			Version:        types.StringValue(a.Version),
			ExternalIP:     types.StringValue(a.ExternalIP),
			InternalIP:     types.StringValue(a.InternalIP),
			Connected:      types.BoolValue(a.Connected),
			Inactive:       types.BoolValue(a.Inactive),
			LastCheckin:    types.Int64Value(a.LastCheckin),
			CreatedAt:      types.Int64Value(a.CreatedAt),
			UpdatedAt:      types.Int64Value(a.UpdatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
