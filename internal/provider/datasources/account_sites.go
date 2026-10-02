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
	_ datasource.DataSource              = &AccountSitesDataSource{}
	_ datasource.DataSourceWithConfigure = &AccountSitesDataSource{}
)

type AccountSitesDataSource struct {
	client *client.Client
}

type AccountSitesDataSourceModel struct {
	Search types.String  `tfsdk:"search"`
	Sites  []SiteDSModel `tfsdk:"sites"`
}

func NewAccountSitesDataSource() datasource.DataSource {
	return &AccountSitesDataSource{}
}

func (d *AccountSitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_sites"
}

func (d *AccountSitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all sites across the RunZero account.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Description: "Optional search query to filter sites.",
				Optional:    true,
			},
			"sites": schema.ListNestedAttribute{
				Description: "List of sites across all organizations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the site.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the site.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the site.",
							Computed:    true,
						},
						"scope": schema.StringAttribute{
							Description: "The scan scope for this site.",
							Computed:    true,
						},
						"excludes": schema.StringAttribute{
							Description: "Excluded scan targets.",
							Computed:    true,
						},
						"created_at": schema.Int64Attribute{
							Description: "Unix timestamp when the site was created.",
							Computed:    true,
						},
						"updated_at": schema.Int64Attribute{
							Description: "Unix timestamp when the site was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AccountSitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AccountSitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AccountSitesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sites, err := d.client.ListAccountSites(ctx, state.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading account sites", err.Error())
		return
	}

	state.Sites = make([]SiteDSModel, len(sites))
	for i, s := range sites {
		state.Sites[i] = SiteDSModel{
			ID:          types.StringValue(s.ID),
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			Scope:       types.StringValue(s.Scope),
			Excludes:    types.StringValue(s.Excludes),
			CreatedAt:   types.Int64Value(s.CreatedAt),
			UpdatedAt:   types.Int64Value(s.UpdatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
