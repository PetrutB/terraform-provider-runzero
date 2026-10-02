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
	_ datasource.DataSource              = &OrganizationDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationDataSource{}
)

type OrganizationDataSource struct {
	client *client.Client
}

type OrganizationDataSourceModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	ClientID                        types.String `tfsdk:"client_id"`
	Inactive                        types.Bool   `tfsdk:"inactive"`
	Project                         types.Bool   `tfsdk:"project"`
	ParentID                        types.String `tfsdk:"parent_id"`
	AssetCount                      types.Int64  `tfsdk:"asset_count"`
	ServiceCount                    types.Int64  `tfsdk:"service_count"`
	ExpirationAssetsStale           types.Int64  `tfsdk:"expiration_assets_stale"`
	ExpirationAssetsOffline         types.Int64  `tfsdk:"expiration_assets_offline"`
	ExpirationIntegrationAttributes types.Int64  `tfsdk:"expiration_integration_attributes"`
	ExpirationScans                 types.Int64  `tfsdk:"expiration_scans"`
	ExpirationVulnerabilities       types.Int64  `tfsdk:"expiration_vulnerabilities"`
	KeepLatestIntegrationAttributes types.Bool   `tfsdk:"keep_latest_integration_attributes"`
	ExportToken                     types.String `tfsdk:"export_token"`
	ExportTokenCreatedAt            types.Int64  `tfsdk:"export_token_created_at"`
	ExportTokenLastUsedAt           types.Int64  `tfsdk:"export_token_last_used_at"`
	ExportTokenLastUsedBy           types.String `tfsdk:"export_token_last_used_by"`
	ExportTokenCounter              types.Int64  `tfsdk:"export_token_counter"`
	CreatedAt                       types.Int64  `tfsdk:"created_at"`
	UpdatedAt                       types.Int64  `tfsdk:"updated_at"`
}

func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

func (d *OrganizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := orgDSSchemaAttributes()
	attrs["id"] = schema.StringAttribute{
		Description: "The UUID of the organization.",
		Required:    true,
	}
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero organization by ID.",
		Attributes:  attrs,
	}
}

func (d *OrganizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OrganizationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := d.client.GetOrganization(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading organization", err.Error())
		return
	}

	ds := mapOrgToDS(org)
	state.ID = ds.ID
	state.Name = ds.Name
	state.Description = ds.Description
	state.ClientID = ds.ClientID
	state.Inactive = ds.Inactive
	state.Project = ds.Project
	state.ParentID = ds.ParentID
	state.AssetCount = ds.AssetCount
	state.ServiceCount = ds.ServiceCount
	state.ExpirationAssetsStale = ds.ExpirationAssetsStale
	state.ExpirationAssetsOffline = ds.ExpirationAssetsOffline
	state.ExpirationIntegrationAttributes = ds.ExpirationIntegrationAttributes
	state.ExpirationScans = ds.ExpirationScans
	state.ExpirationVulnerabilities = ds.ExpirationVulnerabilities
	state.KeepLatestIntegrationAttributes = ds.KeepLatestIntegrationAttributes
	state.ExportToken = ds.ExportToken
	state.ExportTokenCreatedAt = ds.ExportTokenCreatedAt
	state.ExportTokenLastUsedAt = ds.ExportTokenLastUsedAt
	state.ExportTokenLastUsedBy = ds.ExportTokenLastUsedBy
	state.ExportTokenCounter = ds.ExportTokenCounter
	state.CreatedAt = ds.CreatedAt
	state.UpdatedAt = ds.UpdatedAt

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
