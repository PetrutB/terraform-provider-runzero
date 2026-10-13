package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &GroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupsDataSource{}
)

type GroupsDataSource struct {
	client *client.Client
}

type GroupsDataSourceModel struct {
	Groups []GroupDSModel `tfsdk:"groups"`
}

type GroupDSModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	RoleSummary    types.String `tfsdk:"role_summary"`
	UserCount      types.Int64  `tfsdk:"user_count"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	ExpiresAt      types.Int64  `tfsdk:"expires_at"`
	OrgDefaultRole types.String `tfsdk:"org_default_role"`
	OrgRoles       types.Map    `tfsdk:"org_roles"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewGroupsDataSource() datasource.DataSource {
	return &GroupsDataSource{}
}

func (d *GroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *GroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero groups.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Description: "List of groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The UUID of the group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the group.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the group.",
							Computed:    true,
						},
						"role_summary": schema.StringAttribute{
							Description: "Summary of roles assigned to the group.",
							Computed:    true,
						},
						"user_count": schema.Int64Attribute{
							Description: "Number of users in the group.",
							Computed:    true,
						},
						"created_by_email": schema.StringAttribute{
							Description: "Email of the user who created the group.",
							Computed:    true,
						},
						"expires_at": schema.Int64Attribute{
							Description: "Unix timestamp when the group expires.",
							Computed:    true,
						},
						"org_default_role": schema.StringAttribute{
							Description: "The default organization role for the group.",
							Computed:    true,
						},
						"org_roles": schema.MapAttribute{
							Description: "Map of organization IDs to roles.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"created_at": schema.Int64Attribute{
							Description: "Unix timestamp when the group was created.",
							Computed:    true,
						},
						"updated_at": schema.Int64Attribute{
							Description: "Unix timestamp when the group was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *GroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := d.client.ListGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading groups", err.Error())
		return
	}

	state.Groups = make([]GroupDSModel, len(groups))
	for i, g := range groups {
		var orgRoles types.Map
		if g.OrgRoles != nil {
			elements := make(map[string]attr.Value, len(g.OrgRoles))
			for k, v := range g.OrgRoles {
				elements[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
			orgRoles = types.MapValueMust(types.StringType, elements)
		} else {
			orgRoles = types.MapNull(types.StringType)
		}
		state.Groups[i] = GroupDSModel{
			ID:             types.StringValue(g.ID),
			Name:           types.StringValue(g.Name),
			Description:    types.StringValue(g.Description),
			RoleSummary:    types.StringValue(g.RoleSummary),
			UserCount:      types.Int64Value(g.UserCount),
			CreatedByEmail: types.StringValue(g.CreatedByEmail),
			ExpiresAt:      types.Int64Value(g.ExpiresAt),
			OrgDefaultRole: types.StringValue(g.OrgDefaultRole),
			OrgRoles:       orgRoles,
			CreatedAt:      types.Int64Value(g.CreatedAt),
			UpdatedAt:      types.Int64Value(g.UpdatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
