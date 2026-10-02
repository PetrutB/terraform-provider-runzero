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
	_ datasource.DataSource              = &AssetOwnershipTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &AssetOwnershipTypesDataSource{}
)

type AssetOwnershipTypesDataSource struct {
	client *client.Client
}

type AssetOwnershipTypesDataSourceModel struct {
	AssetOwnershipTypes []AssetOwnershipTypeDSModel `tfsdk:"asset_ownership_types"`
}

type AssetOwnershipTypeDSModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Reference types.Int64  `tfsdk:"reference"`
	Order     types.Int64  `tfsdk:"order"`
	Hidden    types.Bool   `tfsdk:"hidden"`
}

func NewAssetOwnershipTypesDataSource() datasource.DataSource {
	return &AssetOwnershipTypesDataSource{}
}

func (d *AssetOwnershipTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_ownership_types"
}

func (d *AssetOwnershipTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero asset ownership types.",
		Attributes: map[string]schema.Attribute{
			"asset_ownership_types": schema.ListNestedAttribute{
				Description: "List of asset ownership types.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the asset ownership type.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the asset ownership type.",
							Computed:    true,
						},
						"reference": schema.Int64Attribute{
							Description: "The reference number of the asset ownership type.",
							Computed:    true,
						},
						"order": schema.Int64Attribute{
							Description: "The display order of the asset ownership type.",
							Computed:    true,
						},
						"hidden": schema.BoolAttribute{
							Description: "Whether the asset ownership type is hidden.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AssetOwnershipTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AssetOwnershipTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AssetOwnershipTypesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipTypes, err := d.client.ListAssetOwnershipTypes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading asset ownership types", err.Error())
		return
	}

	state.AssetOwnershipTypes = make([]AssetOwnershipTypeDSModel, len(ownershipTypes))
	for i, ot := range ownershipTypes {
		state.AssetOwnershipTypes[i] = AssetOwnershipTypeDSModel{
			ID:        types.StringValue(ot.ID),
			Name:      types.StringValue(ot.Name),
			Reference: types.Int64Value(ot.Reference),
			Order:     types.Int64Value(ot.Order),
			Hidden:    types.BoolValue(ot.Hidden),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
