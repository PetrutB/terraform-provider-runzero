package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &AssetDataSource{}
	_ datasource.DataSourceWithConfigure = &AssetDataSource{}
)

type AssetDataSource struct {
	client *client.Client
}

type AssetDataSourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	OrganizationID     types.String  `tfsdk:"organization_id"`
	SiteID             types.String  `tfsdk:"site_id"`
	Scanned            types.Bool    `tfsdk:"scanned"`
	Alive              types.Bool    `tfsdk:"alive"`
	FirstSeen          types.Int64   `tfsdk:"first_seen"`
	LastSeen           types.Int64   `tfsdk:"last_seen"`
	DetectedBy         types.String  `tfsdk:"detected_by"`
	Type               types.String  `tfsdk:"type"`
	OS                 types.String  `tfsdk:"os"`
	OSVersion          types.String  `tfsdk:"os_version"`
	HW                 types.String  `tfsdk:"hw"`
	Addresses          types.List    `tfsdk:"addresses"`
	AddressesExtra     types.List    `tfsdk:"addresses_extra"`
	MACs               types.List    `tfsdk:"macs"`
	MACVendors         types.List    `tfsdk:"mac_vendors"`
	Names              types.List    `tfsdk:"names"`
	Domains            types.List    `tfsdk:"domains"`
	ServiceCount       types.Int64   `tfsdk:"service_count"`
	ServiceCountTCP    types.Int64   `tfsdk:"service_count_tcp"`
	ServiceCountUDP    types.Int64   `tfsdk:"service_count_udp"`
	ServiceCountARP    types.Int64   `tfsdk:"service_count_arp"`
	ServiceCountICMP   types.Int64   `tfsdk:"service_count_icmp"`
	LowestTTL          types.Int64   `tfsdk:"lowest_ttl"`
	LowestRTT          types.Int64   `tfsdk:"lowest_rtt"`
	LastAgentID        types.String  `tfsdk:"last_agent_id"`
	LastTaskID         types.String  `tfsdk:"last_task_id"`
	NewestMAC          types.String  `tfsdk:"newest_mac"`
	NewestMACVendor    types.String  `tfsdk:"newest_mac_vendor"`
	NewestMACAge       types.Int64   `tfsdk:"newest_mac_age"`
	Comments           types.String  `tfsdk:"comments"`
	Criticality        types.String  `tfsdk:"criticality"`
	ServicePortsTCP    types.List    `tfsdk:"service_ports_tcp"`
	ServicePortsUDP    types.List    `tfsdk:"service_ports_udp"`
	ServiceProtocols   types.List    `tfsdk:"service_protocols"`
	ServiceProducts    types.List    `tfsdk:"service_products"`
	OrgName            types.String  `tfsdk:"org_name"`
	SiteName           types.String  `tfsdk:"site_name"`
	AgentName          types.String  `tfsdk:"agent_name"`
	Tags               types.Map     `tfsdk:"tags"`
	Attributes         types.Map     `tfsdk:"attributes"`
	Credentials        types.Map     `tfsdk:"credentials"`
	Risk               types.String  `tfsdk:"risk"`
	RiskRank           types.Float64 `tfsdk:"risk_rank"`
	VulnerabilityCount types.Float64 `tfsdk:"vulnerability_count"`
	FindingCount       types.Float64 `tfsdk:"finding_count"`
	EolOS              types.Float64 `tfsdk:"eol_os"`
	EolOSExt           types.Float64 `tfsdk:"eol_os_ext"`
	Sources            types.List    `tfsdk:"sources"`
	CreatedAt          types.Int64   `tfsdk:"created_at"`
	UpdatedAt          types.Int64   `tfsdk:"updated_at"`
}

func NewAssetDataSource() datasource.DataSource {
	return &AssetDataSource{}
}

func (d *AssetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (d *AssetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := assetDSSchemaAttributes()
	attrs["id"] = schema.StringAttribute{
		Description: "The UUID of the asset.",
		Required:    true,
	}
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero asset by ID.",
		Attributes:  attrs,
	}
}

func (d *AssetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AssetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AssetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := d.client.GetAsset(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading asset", err.Error())
		return
	}

	mapAssetToDS(ctx, asset, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func assetDSSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The UUID of the asset.",
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization UUID.",
			Computed:    true,
		},
		"site_id": schema.StringAttribute{
			Description: "The site ID the asset belongs to.",
			Computed:    true,
		},
		"scanned": schema.BoolAttribute{
			Description: "Whether the asset has been scanned.",
			Computed:    true,
		},
		"alive": schema.BoolAttribute{
			Description: "Whether the asset is alive.",
			Computed:    true,
		},
		"first_seen": schema.Int64Attribute{
			Description: "Unix timestamp when the asset was first seen.",
			Computed:    true,
		},
		"last_seen": schema.Int64Attribute{
			Description: "Unix timestamp when the asset was last seen.",
			Computed:    true,
		},
		"detected_by": schema.StringAttribute{
			Description: "How the asset was detected.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the asset.",
			Computed:    true,
		},
		"os": schema.StringAttribute{
			Description: "The operating system of the asset.",
			Computed:    true,
		},
		"os_version": schema.StringAttribute{
			Description: "The OS version of the asset.",
			Computed:    true,
		},
		"hw": schema.StringAttribute{
			Description: "The hardware description of the asset.",
			Computed:    true,
		},
		"addresses": schema.ListAttribute{
			Description: "IP addresses associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"addresses_extra": schema.ListAttribute{
			Description: "Additional addresses associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"macs": schema.ListAttribute{
			Description: "MAC addresses associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"mac_vendors": schema.ListAttribute{
			Description: "MAC address vendors.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"names": schema.ListAttribute{
			Description: "Hostnames associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"domains": schema.ListAttribute{
			Description: "Domains associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"service_count": schema.Int64Attribute{
			Description: "Total number of services on the asset.",
			Computed:    true,
		},
		"service_count_tcp": schema.Int64Attribute{
			Description: "Number of TCP services.",
			Computed:    true,
		},
		"service_count_udp": schema.Int64Attribute{
			Description: "Number of UDP services.",
			Computed:    true,
		},
		"service_count_arp": schema.Int64Attribute{
			Description: "Number of ARP services.",
			Computed:    true,
		},
		"service_count_icmp": schema.Int64Attribute{
			Description: "Number of ICMP services.",
			Computed:    true,
		},
		"lowest_ttl": schema.Int64Attribute{
			Description: "Lowest TTL observed.",
			Computed:    true,
		},
		"lowest_rtt": schema.Int64Attribute{
			Description: "Lowest round-trip time observed.",
			Computed:    true,
		},
		"last_agent_id": schema.StringAttribute{
			Description: "UUID of the last agent that scanned the asset.",
			Computed:    true,
		},
		"last_task_id": schema.StringAttribute{
			Description: "UUID of the last task that scanned the asset.",
			Computed:    true,
		},
		"newest_mac": schema.StringAttribute{
			Description: "Most recently observed MAC address.",
			Computed:    true,
		},
		"newest_mac_vendor": schema.StringAttribute{
			Description: "Vendor of the most recently observed MAC address.",
			Computed:    true,
		},
		"newest_mac_age": schema.Int64Attribute{
			Description: "Age of the most recently observed MAC address.",
			Computed:    true,
		},
		"comments": schema.StringAttribute{
			Description: "Comments on the asset.",
			Computed:    true,
		},
		"criticality": schema.StringAttribute{
			Description: "The criticality level of the asset.",
			Computed:    true,
		},
		"service_ports_tcp": schema.ListAttribute{
			Description: "TCP service ports.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"service_ports_udp": schema.ListAttribute{
			Description: "UDP service ports.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"service_protocols": schema.ListAttribute{
			Description: "Service protocols detected.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"service_products": schema.ListAttribute{
			Description: "Service products detected.",
			Computed:    true,
			ElementType: types.StringType,
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
		"tags": schema.MapAttribute{
			Description: "Tags associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"attributes": schema.MapAttribute{
			Description: "Attributes associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"credentials": schema.MapAttribute{
			Description: "Credentials associated with the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"risk": schema.StringAttribute{
			Description: "The risk level of the asset.",
			Computed:    true,
		},
		"risk_rank": schema.Float64Attribute{
			Description: "The numeric risk rank.",
			Computed:    true,
		},
		"vulnerability_count": schema.Float64Attribute{
			Description: "Number of vulnerabilities found.",
			Computed:    true,
		},
		"finding_count": schema.Float64Attribute{
			Description: "Number of findings.",
			Computed:    true,
		},
		"eol_os": schema.Float64Attribute{
			Description: "End-of-life date for the operating system.",
			Computed:    true,
		},
		"eol_os_ext": schema.Float64Attribute{
			Description: "Extended end-of-life date for the operating system.",
			Computed:    true,
		},
		"sources": schema.ListAttribute{
			Description: "Data sources for the asset.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"created_at": schema.Int64Attribute{
			Description: "Unix timestamp when the asset was created.",
			Computed:    true,
		},
		"updated_at": schema.Int64Attribute{
			Description: "Unix timestamp when the asset was last updated.",
			Computed:    true,
		},
	}
}

func mapAssetToDS(ctx context.Context, a *client.Asset, m *AssetDataSourceModel, diags *diag.Diagnostics) {
	m.ID = types.StringValue(a.ID)
	m.OrganizationID = types.StringValue(a.OrganizationID)
	m.SiteID = types.StringValue(a.SiteID)
	m.Scanned = types.BoolValue(a.Scanned)
	m.Alive = types.BoolValue(a.Alive)
	m.FirstSeen = types.Int64Value(a.FirstSeen)
	m.LastSeen = types.Int64Value(a.LastSeen)
	m.DetectedBy = types.StringValue(a.DetectedBy)
	m.Type = types.StringValue(a.Type)
	m.OS = types.StringValue(a.OS)
	m.OSVersion = types.StringValue(a.OSVersion)
	m.HW = types.StringValue(a.HW)
	m.ServiceCount = types.Int64Value(a.ServiceCount)
	m.ServiceCountTCP = types.Int64Value(a.ServiceCountTCP)
	m.ServiceCountUDP = types.Int64Value(a.ServiceCountUDP)
	m.ServiceCountARP = types.Int64Value(a.ServiceCountARP)
	m.ServiceCountICMP = types.Int64Value(a.ServiceCountICMP)
	m.LowestTTL = types.Int64Value(a.LowestTTL)
	m.LowestRTT = types.Int64Value(a.LowestRTT)
	m.LastAgentID = types.StringValue(a.LastAgentID)
	m.LastTaskID = types.StringValue(a.LastTaskID)
	m.NewestMAC = types.StringValue(a.NewestMAC)
	m.NewestMACVendor = types.StringValue(a.NewestMACVendor)
	m.NewestMACAge = types.Int64Value(a.NewestMACAge)
	m.Comments = types.StringValue(a.Comments)
	m.Criticality = types.StringValue(a.Criticality)
	m.OrgName = types.StringValue(a.OrgName)
	m.SiteName = types.StringValue(a.SiteName)
	m.AgentName = types.StringValue(a.AgentName)
	m.Risk = types.StringValue(a.Risk)
	m.RiskRank = types.Float64Value(a.RiskRank)
	m.VulnerabilityCount = types.Float64Value(a.VulnerabilityCount)
	m.FindingCount = types.Float64Value(a.FindingCount)
	m.EolOS = types.Float64Value(a.EolOS)
	m.EolOSExt = types.Float64Value(a.EolOSExt)
	m.CreatedAt = types.Int64Value(a.CreatedAt)
	m.UpdatedAt = types.Int64Value(a.UpdatedAt)

	listFields := []struct {
		target *types.List
		source []string
	}{
		{&m.Addresses, a.Addresses},
		{&m.AddressesExtra, a.AddressesExtra},
		{&m.MACs, a.MACs},
		{&m.MACVendors, a.MACVendors},
		{&m.Names, a.Names},
		{&m.Domains, a.Domains},
		{&m.ServicePortsTCP, a.ServicePortsTCP},
		{&m.ServicePortsUDP, a.ServicePortsUDP},
		{&m.ServiceProtocols, a.ServiceProtocols},
		{&m.ServiceProducts, a.ServiceProducts},
		{&m.Sources, a.Sources},
	}
	for _, f := range listFields {
		v, d := types.ListValueFrom(ctx, types.StringType, f.source)
		diags.Append(d...)
		*f.target = v
	}

	mapFields := []struct {
		target *types.Map
		source map[string]string
	}{
		{&m.Tags, a.Tags},
		{&m.Attributes, a.Attributes},
	}
	for _, f := range mapFields {
		v, d := types.MapValueFrom(ctx, types.StringType, f.source)
		diags.Append(d...)
		*f.target = v
	}

	credStrMap := make(map[string]string, len(a.Credentials))
	for k, v := range a.Credentials {
		credStrMap[k] = fmt.Sprintf("%v", v)
	}
	credMap, d := types.MapValueFrom(ctx, types.StringType, credStrMap)
	diags.Append(d...)
	m.Credentials = credMap
}
