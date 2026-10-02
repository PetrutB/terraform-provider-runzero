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
	_ datasource.DataSource              = &AssetOwnershipTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &AssetOwnershipTypeDataSource{}
)

type AssetOwnershipTypeDataSource struct {
	client *client.Client
}

type AssetOwnershipTypeDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Reference types.Int64  `tfsdk:"reference"`
	Order     types.Int64  `tfsdk:"order"`
	Hidden    types.Bool   `tfsdk:"hidden"`
}

func NewAssetOwnershipTypeDataSource() datasource.DataSource {
	return &AssetOwnershipTypeDataSource{}
}

func (d *AssetOwnershipTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_ownership_type"
}

func (d *AssetOwnershipTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero asset ownership type by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the asset ownership type.",
				Required:    true,
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
	}
}

func (d *AssetOwnershipTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AssetOwnershipTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AssetOwnershipTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ownershipType, err := d.client.GetAssetOwnershipType(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading asset ownership type", err.Error())
		return
	}

	state.ID = types.StringValue(ownershipType.ID)
	state.Name = types.StringValue(ownershipType.Name)
	state.Reference = types.Int64Value(ownershipType.Reference)
	state.Order = types.Int64Value(ownershipType.Order)
	state.Hidden = types.BoolValue(ownershipType.Hidden)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
