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
	_ datasource.DataSource              = &ExplorersDataSource{}
	_ datasource.DataSourceWithConfigure = &ExplorersDataSource{}
)

type ExplorersDataSource struct {
	client *client.Client
}

type ExplorersDataSourceModel struct {
	Explorers []ExplorerDSModel `tfsdk:"explorers"`
}

type ExplorerDSModel struct {
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

func NewExplorersDataSource() datasource.DataSource {
	return &ExplorersDataSource{}
}

func (d *ExplorersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_explorers"
}

func (d *ExplorersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero explorers.",
		Attributes: map[string]schema.Attribute{
			"explorers": schema.ListNestedAttribute{
				Description: "List of explorers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the explorer.",
							Computed:    true,
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
							Description: "The organization ID.",
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
				},
			},
		},
	}
}

func (d *ExplorersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExplorersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ExplorersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	explorers, err := d.client.ListExplorers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading explorers", err.Error())
		return
	}

	state.Explorers = make([]ExplorerDSModel, len(explorers))
	for i, e := range explorers {
		state.Explorers[i] = ExplorerDSModel{
			ID:             types.StringValue(e.ID),
			Name:           types.StringValue(e.Name),
			SiteID:         types.StringValue(e.SiteID),
			OrganizationID: types.StringValue(e.OrganizationID),
			OS:             types.StringValue(e.OS),
			Arch:           types.StringValue(e.Arch),
			Version:        types.StringValue(e.Version),
			ExternalIP:     types.StringValue(e.ExternalIP),
			InternalIP:     types.StringValue(e.InternalIP),
			Connected:      types.BoolValue(e.Connected),
			Inactive:       types.BoolValue(e.Inactive),
			LastCheckin:    types.Int64Value(e.LastCheckin),
			CreatedAt:      types.Int64Value(e.CreatedAt),
			UpdatedAt:      types.Int64Value(e.UpdatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
