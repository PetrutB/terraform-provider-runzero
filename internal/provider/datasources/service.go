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
	_ datasource.DataSource              = &ServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &ServiceDataSource{}
)

type ServiceDataSource struct {
	client *client.Client
}

type ServiceDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	OrganizationID   types.String `tfsdk:"organization_id"`
	SiteID           types.String `tfsdk:"site_id"`
	CreatedAt        types.Int64  `tfsdk:"created_at"`
	UpdatedAt        types.Int64  `tfsdk:"updated_at"`
	ServiceID        types.String `tfsdk:"service_id"`
	ServiceAssetID   types.String `tfsdk:"service_asset_id"`
	ServiceCreatedAt types.Int64  `tfsdk:"service_created_at"`
	ServiceUpdatedAt types.Int64  `tfsdk:"service_updated_at"`
	ServiceAddress   types.String `tfsdk:"service_address"`
	ServiceTransport types.String `tfsdk:"service_transport"`
	ServiceVhost     types.String `tfsdk:"service_vhost"`
	ServicePort      types.String `tfsdk:"service_port"`
	ServiceProtocol  types.String `tfsdk:"service_protocol"`
	ServiceSummary   types.String `tfsdk:"service_summary"`
	ServiceLink      types.String `tfsdk:"service_link"`
	Alive            types.Bool   `tfsdk:"alive"`
	FirstSeen        types.Int64  `tfsdk:"first_seen"`
	LastSeen         types.Int64  `tfsdk:"last_seen"`
	Type             types.String `tfsdk:"type"`
	OS               types.String `tfsdk:"os"`
	HW               types.String `tfsdk:"hw"`
	Comments         types.String `tfsdk:"comments"`
}

func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

func (d *ServiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *ServiceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero service by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the asset.",
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
			"service_created_at": schema.Int64Attribute{
				Description: "Service creation timestamp (Unix epoch).",
				Computed:    true,
			},
			"service_updated_at": schema.Int64Attribute{
				Description: "Service last update timestamp (Unix epoch).",
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
			"service_vhost": schema.StringAttribute{
				Description: "The virtual host of the service.",
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
			"service_link": schema.StringAttribute{
				Description: "A link to the service.",
				Computed:    true,
			},
			"alive": schema.BoolAttribute{
				Description: "Whether the asset is alive.",
				Computed:    true,
			},
			"first_seen": schema.Int64Attribute{
				Description: "First seen timestamp (Unix epoch).",
				Computed:    true,
			},
			"last_seen": schema.Int64Attribute{
				Description: "Last seen timestamp (Unix epoch).",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The asset type.",
				Computed:    true,
			},
			"os": schema.StringAttribute{
				Description: "The operating system.",
				Computed:    true,
			},
			"hw": schema.StringAttribute{
				Description: "The hardware identifier.",
				Computed:    true,
			},
			"comments": schema.StringAttribute{
				Description: "Comments on the asset.",
				Computed:    true,
			},
		},
	}
}

func (d *ServiceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServiceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, err := d.client.GetService(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading service", err.Error())
		return
	}

	state.ID = types.StringValue(service.ID)
	state.OrganizationID = types.StringValue(service.OrganizationID)
	state.SiteID = types.StringValue(service.SiteID)
	state.CreatedAt = types.Int64Value(service.CreatedAt)
	state.UpdatedAt = types.Int64Value(service.UpdatedAt)
	state.ServiceID = types.StringValue(service.ServiceID)
	state.ServiceAssetID = types.StringValue(service.ServiceAssetID)
	state.ServiceCreatedAt = types.Int64Value(service.ServiceCreatedAt)
	state.ServiceUpdatedAt = types.Int64Value(service.ServiceUpdatedAt)
	state.ServiceAddress = types.StringValue(service.ServiceAddress)
	state.ServiceTransport = types.StringValue(service.ServiceTransport)
	state.ServiceVhost = types.StringValue(service.ServiceVhost)
	state.ServicePort = types.StringValue(service.ServicePort)
	state.ServiceProtocol = types.StringValue(service.ServiceProtocol)
	state.ServiceSummary = types.StringValue(service.ServiceSummary)
	state.ServiceLink = types.StringValue(service.ServiceLink)
	state.Alive = types.BoolValue(service.Alive)
	state.FirstSeen = types.Int64Value(service.FirstSeen)
	state.LastSeen = types.Int64Value(service.LastSeen)
	state.Type = types.StringValue(service.Type)
	state.OS = types.StringValue(service.OS)
	state.HW = types.StringValue(service.HW)
	state.Comments = types.StringValue(service.Comments)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
