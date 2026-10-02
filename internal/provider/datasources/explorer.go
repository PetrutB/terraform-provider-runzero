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
	_ datasource.DataSource              = &ExplorerDataSource{}
	_ datasource.DataSourceWithConfigure = &ExplorerDataSource{}
)

type ExplorerDataSource struct {
	client *client.Client
}

type ExplorerDataSourceModel struct {
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

func NewExplorerDataSource() datasource.DataSource {
	return &ExplorerDataSource{}
}

func (d *ExplorerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_explorer"
}

func (d *ExplorerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero explorer by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the explorer.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the explorer.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID the explorer is associated with.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID the explorer belongs to.",
				Computed:    true,
			},
			"os": schema.StringAttribute{
				Description: "The operating system of the explorer host.",
				Computed:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the explorer host.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of the explorer.",
				Computed:    true,
			},
			"external_ip": schema.StringAttribute{
				Description: "The external IP address of the explorer.",
				Computed:    true,
			},
			"internal_ip": schema.StringAttribute{
				Description: "The internal IP address of the explorer.",
				Computed:    true,
			},
			"connected": schema.BoolAttribute{
				Description: "Whether the explorer is currently connected.",
				Computed:    true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the explorer is inactive.",
				Computed:    true,
			},
			"last_checkin": schema.Int64Attribute{
				Description: "Unix timestamp of the explorer's last check-in.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the explorer was created.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the explorer was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *ExplorerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExplorerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ExplorerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	explorer, err := d.client.GetExplorer(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading explorer", err.Error())
		return
	}

	state.ID = types.StringValue(explorer.ID)
	state.Name = types.StringValue(explorer.Name)
	state.SiteID = types.StringValue(explorer.SiteID)
	state.OrganizationID = types.StringValue(explorer.OrganizationID)
	state.OS = types.StringValue(explorer.OS)
	state.Arch = types.StringValue(explorer.Arch)
	state.Version = types.StringValue(explorer.Version)
	state.ExternalIP = types.StringValue(explorer.ExternalIP)
	state.InternalIP = types.StringValue(explorer.InternalIP)
	state.Connected = types.BoolValue(explorer.Connected)
	state.Inactive = types.BoolValue(explorer.Inactive)
	state.LastCheckin = types.Int64Value(explorer.LastCheckin)
	state.CreatedAt = types.Int64Value(explorer.CreatedAt)
	state.UpdatedAt = types.Int64Value(explorer.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
