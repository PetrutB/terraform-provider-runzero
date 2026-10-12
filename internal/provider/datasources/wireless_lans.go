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
	_ datasource.DataSource              = &WirelessLANsDataSource{}
	_ datasource.DataSourceWithConfigure = &WirelessLANsDataSource{}
)

type WirelessLANsDataSource struct {
	client *client.Client
}

type WirelessLANsDataSourceModel struct {
	WirelessLANs []WirelessLANDSModel `tfsdk:"wireless_lans"`
}

type WirelessLANDSModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	SiteID         types.String `tfsdk:"site_id"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	LastSeen       types.Int64  `tfsdk:"last_seen"`
	ESSID          types.String `tfsdk:"essid"`
	BSSID          types.String `tfsdk:"bssid"`
	Type           types.String `tfsdk:"type"`
	Authentication types.String `tfsdk:"authentication"`
	Encryption     types.String `tfsdk:"encryption"`
	Signal         types.Int64  `tfsdk:"signal"`
	Channels       types.String `tfsdk:"channels"`
	Vendor         types.String `tfsdk:"vendor"`
}

func NewWirelessLANsDataSource() datasource.DataSource {
	return &WirelessLANsDataSource{}
}

func (d *WirelessLANsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_lans"
}

func (d *WirelessLANsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero wireless LANs.",
		Attributes: map[string]schema.Attribute{
			"wireless_lans": schema.ListNestedAttribute{
				Description: "List of wireless LANs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the wireless LAN.",
							Computed:    true,
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
						"essid": schema.StringAttribute{
							Description: "The ESSID (network name).",
							Computed:    true,
						},
						"bssid": schema.StringAttribute{
							Description: "The BSSID (MAC address).",
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
							Description: "The channels used.",
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "The vendor of the access point.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WirelessLANsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WirelessLANsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state WirelessLANsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lans, err := d.client.ListWirelessLANs(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading wireless LANs", err.Error())
		return
	}

	state.WirelessLANs = make([]WirelessLANDSModel, len(lans))
	for i, l := range lans {
		state.WirelessLANs[i] = WirelessLANDSModel{
			ID:             types.StringValue(l.ID),
			OrganizationID: types.StringValue(l.OrganizationID),
			SiteID:         types.StringValue(l.SiteID),
			CreatedAt:      types.Int64Value(l.CreatedAt),
			LastSeen:       types.Int64Value(l.LastSeen),
			ESSID:          types.StringValue(l.ESSID),
			BSSID:          types.StringValue(l.BSSID),
			Type:           types.StringValue(l.Type),
			Authentication: types.StringValue(l.Authentication),
			Encryption:     types.StringValue(l.Encryption),
			Signal:         types.Int64Value(l.Signal),
			Channels:       types.StringValue(l.Channels),
			Vendor:         types.StringValue(l.Vendor),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
