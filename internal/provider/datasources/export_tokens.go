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
	_ datasource.DataSource              = &ExportTokensDataSource{}
	_ datasource.DataSourceWithConfigure = &ExportTokensDataSource{}
)

type ExportTokensDataSource struct {
	client *client.Client
}

type ExportTokensDataSourceModel struct {
	OrgID        types.String         `tfsdk:"org_id"`
	ExportTokens []ExportTokenDSModel `tfsdk:"export_tokens"`
}

type ExportTokenDSModel struct {
	ID         types.String `tfsdk:"id"`
	Token      types.String `tfsdk:"token"`
	Comment    types.String `tfsdk:"comment"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	LastUsedAt types.Int64  `tfsdk:"last_used_at"`
	LastUsedBy types.String `tfsdk:"last_used_by"`
	Counter    types.Int64  `tfsdk:"counter"`
}

func NewExportTokensDataSource() datasource.DataSource {
	return &ExportTokensDataSource{}
}

func (d *ExportTokensDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export_tokens"
}

func (d *ExportTokensDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero export tokens.",
		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				Description: "The organization ID to list export tokens for.",
				Required:    true,
			},
			"export_tokens": schema.ListNestedAttribute{
				Description: "List of export tokens.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the export token.",
							Computed:    true,
						},
						"token": schema.StringAttribute{
							Description: "The export token value.",
							Computed:    true,
							Sensitive:   true,
						},
						"comment": schema.StringAttribute{
							Description: "A comment for the export token.",
							Computed:    true,
						},
						"created_at": schema.Int64Attribute{
							Description: "Unix timestamp when the export token was created.",
							Computed:    true,
						},
						"last_used_at": schema.Int64Attribute{
							Description: "Unix timestamp when the token was last used.",
							Computed:    true,
						},
						"last_used_by": schema.StringAttribute{
							Description: "Identifier of the last user of the token.",
							Computed:    true,
						},
						"counter": schema.Int64Attribute{
							Description: "Total usage counter.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ExportTokensDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExportTokensDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ExportTokensDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokens, err := d.client.ListExportTokens(ctx, state.OrgID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading export tokens", err.Error())
		return
	}

	state.ExportTokens = make([]ExportTokenDSModel, len(tokens))
	for i, t := range tokens {
		state.ExportTokens[i] = ExportTokenDSModel{
			ID:         types.StringValue(t.ID),
			Token:      types.StringValue(t.Token),
			Comment:    types.StringValue(t.Comment),
			CreatedAt:  types.Int64Value(t.CreatedAt),
			LastUsedAt: types.Int64Value(t.LastUsedAt),
			LastUsedBy: types.StringValue(t.LastUsedBy),
			Counter:    types.Int64Value(t.Counter),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
