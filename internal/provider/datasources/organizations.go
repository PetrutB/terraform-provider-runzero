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
	_ datasource.DataSource              = &OrganizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationsDataSource{}
)

type OrganizationsDataSource struct {
	client *client.Client
}

type OrganizationsDataSourceModel struct {
	Organizations []OrganizationDSModel `tfsdk:"organizations"`
}

type OrganizationDSModel struct {
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

func NewOrganizationsDataSource() datasource.DataSource {
	return &OrganizationsDataSource{}
}

func (d *OrganizationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *OrganizationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero organizations.",
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Description: "List of organizations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: orgDSSchemaAttributes(),
				},
			},
		},
	}
}

func orgDSSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The UUID of the organization.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the organization.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "Description of the organization.",
			Computed:    true,
		},
		"client_id": schema.StringAttribute{
			Description: "The client ID of the organization.",
			Computed:    true,
		},
		"inactive": schema.BoolAttribute{
			Description: "Whether the organization is inactive.",
			Computed:    true,
		},
		"project": schema.BoolAttribute{
			Description: "Whether this is a project organization.",
			Computed:    true,
		},
		"parent_id": schema.StringAttribute{
			Description: "The parent organization ID.",
			Computed:    true,
		},
		"asset_count": schema.Int64Attribute{
			Description: "Number of assets in the organization.",
			Computed:    true,
		},
		"service_count": schema.Int64Attribute{
			Description: "Number of services in the organization.",
			Computed:    true,
		},
		"expiration_assets_stale": schema.Int64Attribute{
			Description: "Duration in seconds before assets are marked stale.",
			Computed:    true,
		},
		"expiration_assets_offline": schema.Int64Attribute{
			Description: "Duration in seconds before assets are marked offline.",
			Computed:    true,
		},
		"expiration_integration_attributes": schema.Int64Attribute{
			Description: "Duration in seconds before integration attributes expire.",
			Computed:    true,
		},
		"expiration_scans": schema.Int64Attribute{
			Description: "Duration in seconds before scans expire.",
			Computed:    true,
		},
		"expiration_vulnerabilities": schema.Int64Attribute{
			Description: "Duration in seconds before vulnerabilities expire.",
			Computed:    true,
		},
		"keep_latest_integration_attributes": schema.BoolAttribute{
			Description: "Whether to keep the latest integration attributes.",
			Computed:    true,
		},
		"export_token": schema.StringAttribute{
			Description: "The organization data export token.",
			Computed:    true,
			Sensitive:   true,
		},
		"export_token_created_at": schema.Int64Attribute{
			Description: "Unix timestamp when the export token was created.",
			Computed:    true,
		},
		"export_token_last_used_at": schema.Int64Attribute{
			Description: "Unix timestamp when the export token was last used.",
			Computed:    true,
		},
		"export_token_last_used_by": schema.StringAttribute{
			Description: "Identifier of who last used the export token.",
			Computed:    true,
		},
		"export_token_counter": schema.Int64Attribute{
			Description: "Number of times the export token has been used.",
			Computed:    true,
		},
		"created_at": schema.Int64Attribute{
			Description: "Unix timestamp when the organization was created.",
			Computed:    true,
		},
		"updated_at": schema.Int64Attribute{
			Description: "Unix timestamp when the organization was last updated.",
			Computed:    true,
		},
	}
}

func (d *OrganizationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OrganizationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgs, err := d.client.ListOrganizations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading organizations", err.Error())
		return
	}

	state.Organizations = make([]OrganizationDSModel, len(orgs))
	for i, o := range orgs {
		state.Organizations[i] = mapOrgToDS(&o)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func mapOrgToDS(o *client.Organization) OrganizationDSModel {
	return OrganizationDSModel{
		ID:                              types.StringValue(o.ID),
		Name:                            types.StringValue(o.Name),
		Description:                     types.StringValue(o.Description),
		ClientID:                        types.StringValue(o.ClientID),
		Inactive:                        types.BoolValue(o.Inactive),
		Project:                         types.BoolValue(o.Project),
		ParentID:                        types.StringValue(o.ParentID),
		AssetCount:                      types.Int64Value(o.AssetCount),
		ServiceCount:                    types.Int64Value(o.ServiceCount),
		ExpirationAssetsStale:           types.Int64Value(o.ExpirationAssetsStale),
		ExpirationAssetsOffline:         types.Int64Value(o.ExpirationAssetsOffline),
		ExpirationIntegrationAttributes: types.Int64Value(o.ExpirationIntegrationAttributes),
		ExpirationScans:                 types.Int64Value(o.ExpirationScans),
		ExpirationVulnerabilities:       types.Int64Value(o.ExpirationVulnerabilities),
		KeepLatestIntegrationAttributes: types.BoolValue(o.KeepLatestIntegrationAttributes),
		ExportToken:                     types.StringValue(o.ExportToken),
		ExportTokenCreatedAt:            types.Int64Value(o.ExportTokenCreatedAt),
		ExportTokenLastUsedAt:           types.Int64Value(o.ExportTokenLastUsedAt),
		ExportTokenLastUsedBy:           types.StringValue(o.ExportTokenLastUsedBy),
		ExportTokenCounter:              types.Int64Value(o.ExportTokenCounter),
		CreatedAt:                       types.Int64Value(o.CreatedAt),
		UpdatedAt:                       types.Int64Value(o.UpdatedAt),
	}
}
