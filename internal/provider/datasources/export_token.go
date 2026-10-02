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
	_ datasource.DataSource              = &ExportTokenDataSource{}
	_ datasource.DataSourceWithConfigure = &ExportTokenDataSource{}
)

type ExportTokenDataSource struct {
	client *client.Client
}

type ExportTokenDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	OrgID      types.String `tfsdk:"org_id"`
	Token      types.String `tfsdk:"token"`
	Comment    types.String `tfsdk:"comment"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	LastUsedAt types.Int64  `tfsdk:"last_used_at"`
	LastUsedBy types.String `tfsdk:"last_used_by"`
	Counter    types.Int64  `tfsdk:"counter"`
}

func NewExportTokenDataSource() datasource.DataSource {
	return &ExportTokenDataSource{}
}

func (d *ExportTokenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export_token"
}

func (d *ExportTokenDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero export token by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the export token.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The organization ID for this export token.",
				Required:    true,
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
	}
}

func (d *ExportTokenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExportTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ExportTokenDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := d.client.GetExportToken(ctx, state.OrgID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading export token", err.Error())
		return
	}

	state.ID = types.StringValue(token.ID)
	state.Token = types.StringValue(token.Token)
	state.Comment = types.StringValue(token.Comment)
	state.CreatedAt = types.Int64Value(token.CreatedAt)
	state.LastUsedAt = types.Int64Value(token.LastUsedAt)
	state.LastUsedBy = types.StringValue(token.LastUsedBy)
	state.Counter = types.Int64Value(token.Counter)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
