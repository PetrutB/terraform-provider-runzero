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
	_ datasource.DataSource              = &ServicesDataSource{}
	_ datasource.DataSourceWithConfigure = &ServicesDataSource{}
)

type ServicesDataSource struct {
	client *client.Client
}

type ServicesDataSourceModel struct {
	Search   types.String     `tfsdk:"search"`
	Services []ServiceDSModel `tfsdk:"services"`
}

type ServiceDSModel struct {
	ID               types.String `tfsdk:"id"`
	OrganizationID   types.String `tfsdk:"organization_id"`
	SiteID           types.String `tfsdk:"site_id"`
	CreatedAt        types.Int64  `tfsdk:"created_at"`
	UpdatedAt        types.Int64  `tfsdk:"updated_at"`
	ServiceID        types.String `tfsdk:"service_id"`
	ServiceAssetID   types.String `tfsdk:"service_asset_id"`
	ServiceAddress   types.String `tfsdk:"service_address"`
	ServiceTransport types.String `tfsdk:"service_transport"`
	ServicePort      types.String `tfsdk:"service_port"`
	ServiceProtocol  types.String `tfsdk:"service_protocol"`
	ServiceSummary   types.String `tfsdk:"service_summary"`
}

func NewServicesDataSource() datasource.DataSource {
	return &ServicesDataSource{}
}

func (d *ServicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *ServicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero services.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Description: "Optional search query to filter services.",
				Optional:    true,
			},
			"services": schema.ListNestedAttribute{
				Description: "List of services.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the asset.",
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
							Description: "Asset creation timestamp (Unix epoch).",
							Computed:    true,
						},
						"updated_at": schema.Int64Attribute{
							Description: "Asset last update timestamp (Unix epoch).",
							Computed:    true,
						},
						"service_id": schema.StringAttribute{
							Description: "The UUID of the service.",
							Computed:    true,
						},
						"service_asset_id": schema.StringAttribute{
							Description: "The asset ID associated with the service.",
							Computed:    true,
						},
						"service_address": schema.StringAttribute{
							Description: "The address of the service.",
							Computed:    true,
						},
						"service_transport": schema.StringAttribute{
							Description: "The transport protocol (tcp, udp).",
							Computed:    true,
						},
						"service_port": schema.StringAttribute{
							Description: "The port of the service.",
							Computed:    true,
						},
						"service_protocol": schema.StringAttribute{
							Description: "The application protocol of the service.",
							Computed:    true,
						},
						"service_summary": schema.StringAttribute{
							Description: "A summary of the service.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ServicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServicesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	services, err := d.client.ListServices(ctx, state.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading services", err.Error())
		return
	}

	state.Services = make([]ServiceDSModel, len(services))
	for i, s := range services {
		state.Services[i] = ServiceDSModel{
			ID:               types.StringValue(s.ID),
			OrganizationID:   types.StringValue(s.OrganizationID),
			SiteID:           types.StringValue(s.SiteID),
			CreatedAt:        types.Int64Value(s.CreatedAt),
			UpdatedAt:        types.Int64Value(s.UpdatedAt),
			ServiceID:        types.StringValue(s.ServiceID),
			ServiceAssetID:   types.StringValue(s.ServiceAssetID),
			ServiceAddress:   types.StringValue(s.ServiceAddress),
			ServiceTransport: types.StringValue(s.ServiceTransport),
			ServicePort:      types.StringValue(s.ServicePort),
			ServiceProtocol:  types.StringValue(s.ServiceProtocol),
			ServiceSummary:   types.StringValue(s.ServiceSummary),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
