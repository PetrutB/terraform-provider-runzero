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
	_ datasource.DataSource              = &AssetsDataSource{}
	_ datasource.DataSourceWithConfigure = &AssetsDataSource{}
)

type AssetsDataSource struct {
	client *client.Client
}

type AssetsDataSourceModel struct {
	Search types.String           `tfsdk:"search"`
	Assets []AssetDataSourceModel `tfsdk:"assets"`
}

func NewAssetsDataSource() datasource.DataSource {
	return &AssetsDataSource{}
}

func (d *AssetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assets"
}

func (d *AssetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero assets.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Description: "Optional search query to filter assets.",
				Optional:    true,
			},
			"assets": schema.ListNestedAttribute{
				Description: "List of assets.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: assetDSSchemaAttributes(),
				},
			},
		},
	}
}

func (d *AssetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AssetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AssetsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assets, err := d.client.ListAssets(ctx, state.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading assets", err.Error())
		return
	}

	state.Assets = make([]AssetDataSourceModel, len(assets))
	for i, a := range assets {
		mapAssetToDS(ctx, &a, &state.Assets[i], &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
