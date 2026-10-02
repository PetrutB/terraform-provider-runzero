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
	_ datasource.DataSource              = &SSOGroupMappingDataSource{}
	_ datasource.DataSourceWithConfigure = &SSOGroupMappingDataSource{}
)

type SSOGroupMappingDataSource struct {
	client *client.Client
}

type SSOGroupMappingDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	GroupID        types.String `tfsdk:"group_id"`
	GroupName      types.String `tfsdk:"group_name"`
	SSOAttribute   types.String `tfsdk:"sso_attribute"`
	SSOValue       types.String `tfsdk:"sso_value"`
	Description    types.String `tfsdk:"description"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewSSOGroupMappingDataSource() datasource.DataSource {
	return &SSOGroupMappingDataSource{}
}

func (d *SSOGroupMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_group_mapping"
}

func (d *SSOGroupMappingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero SSO group mapping by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the SSO group mapping.",
				Required:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "The UUID of the group.",
				Computed:    true,
			},
			"group_name": schema.StringAttribute{
				Description: "The name of the group.",
				Computed:    true,
			},
			"sso_attribute": schema.StringAttribute{
				Description: "The SSO attribute used for mapping.",
				Computed:    true,
			},
			"sso_value": schema.StringAttribute{
				Description: "The SSO attribute value used for mapping.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the SSO group mapping.",
				Computed:    true,
			},
			"created_by_email": schema.StringAttribute{
				Description: "Email of the user who created the mapping.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the mapping was created.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the mapping was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *SSOGroupMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SSOGroupMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SSOGroupMappingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := d.client.GetSSOGroupMapping(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSO group mapping", err.Error())
		return
	}

	state.ID = types.StringValue(mapping.ID)
	state.GroupID = types.StringValue(mapping.GroupID)
	state.GroupName = types.StringValue(mapping.GroupName)
	state.SSOAttribute = types.StringValue(mapping.SSOAttribute)
	state.SSOValue = types.StringValue(mapping.SSOValue)
	state.Description = types.StringValue(mapping.Description)
	state.CreatedByEmail = types.StringValue(mapping.CreatedByEmail)
	state.CreatedAt = types.Int64Value(mapping.CreatedAt)
	state.UpdatedAt = types.Int64Value(mapping.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
