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
	_ datasource.DataSource              = &WirelessLANDataSource{}
	_ datasource.DataSourceWithConfigure = &WirelessLANDataSource{}
)

type WirelessLANDataSource struct {
	client *client.Client
}

type WirelessLANDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	SiteID         types.String `tfsdk:"site_id"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	LastSeen       types.Int64  `tfsdk:"last_seen"`
	LastAgentID    types.String `tfsdk:"last_agent_id"`
	LastTaskID     types.String `tfsdk:"last_task_id"`
	ESSID          types.String `tfsdk:"essid"`
	BSSID          types.String `tfsdk:"bssid"`
	Type           types.String `tfsdk:"type"`
	Authentication types.String `tfsdk:"authentication"`
	Encryption     types.String `tfsdk:"encryption"`
	Signal         types.Int64  `tfsdk:"signal"`
	Channels       types.String `tfsdk:"channels"`
	Interface      types.String `tfsdk:"interface"`
	Vendor         types.String `tfsdk:"vendor"`
	Family         types.String `tfsdk:"family"`
	OrgName        types.String `tfsdk:"org_name"`
	SiteName       types.String `tfsdk:"site_name"`
	AgentName      types.String `tfsdk:"agent_name"`
}

func NewWirelessLANDataSource() datasource.DataSource {
	return &WirelessLANDataSource{}
}

func (d *WirelessLANDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_lan"
}

func (d *WirelessLANDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero wireless LAN by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the wireless LAN.",
				Required:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Creation timestamp (Unix epoch).",
				Computed:    true,
			},
			"last_seen": schema.Int64Attribute{
				Description: "Last seen timestamp (Unix epoch).",
				Computed:    true,
			},
			"last_agent_id": schema.StringAttribute{
				Description: "The ID of the last agent that observed this network.",
				Computed:    true,
			},
			"last_task_id": schema.StringAttribute{
				Description: "The ID of the last task that observed this network.",
				Computed:    true,
			},
			"essid": schema.StringAttribute{
				Description: "The ESSID (network name) of the wireless LAN.",
				Computed:    true,
			},
			"bssid": schema.StringAttribute{
				Description: "The BSSID (MAC address) of the access point.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the wireless network.",
				Computed:    true,
			},
			"authentication": schema.StringAttribute{
				Description: "The authentication method.",
				Computed:    true,
			},
			"encryption": schema.StringAttribute{
				Description: "The encryption type.",
				Computed:    true,
			},
			"signal": schema.Int64Attribute{
				Description: "The signal strength (0-100).",
				Computed:    true,
			},
			"channels": schema.StringAttribute{
				Description: "The channels used by the wireless LAN.",
				Computed:    true,
			},
			"interface": schema.StringAttribute{
				Description: "The network interface.",
				Computed:    true,
			},
			"vendor": schema.StringAttribute{
				Description: "The vendor of the access point.",
				Computed:    true,
			},
			"family": schema.StringAttribute{
				Description: "The device family.",
				Computed:    true,
			},
			"org_name": schema.StringAttribute{
				Description: "The organization name.",
				Computed:    true,
			},
			"site_name": schema.StringAttribute{
				Description: "The site name.",
				Computed:    true,
			},
			"agent_name": schema.StringAttribute{
				Description: "The agent name.",
				Computed:    true,
			},
		},
	}
}

func (d *WirelessLANDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WirelessLANDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state WirelessLANDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lan, err := d.client.GetWirelessLAN(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading wireless LAN", err.Error())
		return
	}

	state.ID = types.StringValue(lan.ID)
	state.OrganizationID = types.StringValue(lan.OrganizationID)
	state.SiteID = types.StringValue(lan.SiteID)
	state.CreatedAt = types.Int64Value(lan.CreatedAt)
	state.LastSeen = types.Int64Value(lan.LastSeen)
	state.LastAgentID = types.StringValue(lan.LastAgentID)
	state.LastTaskID = types.StringValue(lan.LastTaskID)
	state.ESSID = types.StringValue(lan.ESSID)
	state.BSSID = types.StringValue(lan.BSSID)
	state.Type = types.StringValue(lan.Type)
	state.Authentication = types.StringValue(lan.Authentication)
	state.Encryption = types.StringValue(lan.Encryption)
	state.Signal = types.Int64Value(int64(lan.Signal))
	state.Channels = types.StringValue(lan.Channels)
	state.Interface = types.StringValue(lan.Interface)
	state.Vendor = types.StringValue(lan.Vendor)
	state.Family = types.StringValue(lan.Family)
	state.OrgName = types.StringValue(lan.OrgName)
	state.SiteName = types.StringValue(lan.SiteName)
	state.AgentName = types.StringValue(lan.AgentName)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
