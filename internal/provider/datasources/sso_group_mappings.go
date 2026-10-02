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
	_ datasource.DataSource              = &SSOGroupMappingsDataSource{}
	_ datasource.DataSourceWithConfigure = &SSOGroupMappingsDataSource{}
)

type SSOGroupMappingsDataSource struct {
	client *client.Client
}

type SSOGroupMappingsDataSourceModel struct {
	SSOGroupMappings []SSOGroupMappingDSModel `tfsdk:"sso_group_mappings"`
}

type SSOGroupMappingDSModel struct {
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

func NewSSOGroupMappingsDataSource() datasource.DataSource {
	return &SSOGroupMappingsDataSource{}
}

func (d *SSOGroupMappingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_group_mappings"
}

func (d *SSOGroupMappingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero SSO group mappings.",
		Attributes: map[string]schema.Attribute{
			"sso_group_mappings": schema.ListNestedAttribute{
				Description: "List of SSO group mappings.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the SSO group mapping.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *SSOGroupMappingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SSOGroupMappingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SSOGroupMappingsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mappings, err := d.client.ListSSOGroupMappings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSO group mappings", err.Error())
		return
	}

	state.SSOGroupMappings = make([]SSOGroupMappingDSModel, len(mappings))
	for i, m := range mappings {
		state.SSOGroupMappings[i] = SSOGroupMappingDSModel{
			ID:             types.StringValue(m.ID),
			GroupID:        types.StringValue(m.GroupID),
			GroupName:      types.StringValue(m.GroupName),
			SSOAttribute:   types.StringValue(m.SSOAttribute),
			SSOValue:       types.StringValue(m.SSOValue),
			Description:    types.StringValue(m.Description),
			CreatedByEmail: types.StringValue(m.CreatedByEmail),
			CreatedAt:      types.Int64Value(m.CreatedAt),
			UpdatedAt:      types.Int64Value(m.UpdatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
