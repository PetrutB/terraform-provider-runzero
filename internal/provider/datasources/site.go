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
	_ datasource.DataSource              = &SiteDataSource{}
	_ datasource.DataSourceWithConfigure = &SiteDataSource{}
)

type SiteDataSource struct {
	client *client.Client
}

type SiteDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Scope       types.String `tfsdk:"scope"`
	Excludes    types.String `tfsdk:"excludes"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func NewSiteDataSource() datasource.DataSource {
	return &SiteDataSource{}
}

func (d *SiteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (d *SiteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero site by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the site.",
				Required:    true,
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
	}
}

func (d *SiteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SiteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := d.client.GetSite(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}

	state.ID = types.StringValue(site.ID)
	state.Name = types.StringValue(site.Name)
	state.Description = types.StringValue(site.Description)
	state.Scope = types.StringValue(site.Scope)
	state.Excludes = types.StringValue(site.Excludes)
	state.CreatedAt = types.Int64Value(site.CreatedAt)
	state.UpdatedAt = types.Int64Value(site.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
