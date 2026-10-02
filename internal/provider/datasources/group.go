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
	_ datasource.DataSource              = &GroupDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupDataSource{}
)

type GroupDataSource struct {
	client *client.Client
}

type GroupDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	RoleSummary    types.String `tfsdk:"role_summary"`
	UserCount      types.Int64  `tfsdk:"user_count"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	ExpiresAt      types.Int64  `tfsdk:"expires_at"`
	OrgDefaultRole types.String `tfsdk:"org_default_role"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

func (d *GroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero group by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the group.",
				Required:    true,
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
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the group was created.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the group was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *GroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetGroup(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading group", err.Error())
		return
	}

	state.ID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)
	state.RoleSummary = types.StringValue(group.RoleSummary)
	state.UserCount = types.Int64Value(group.UserCount)
	state.CreatedByEmail = types.StringValue(group.CreatedByEmail)
	state.ExpiresAt = types.Int64Value(group.ExpiresAt)
	state.OrgDefaultRole = types.StringValue(group.OrgDefaultRole)
	state.CreatedAt = types.Int64Value(group.CreatedAt)
	state.UpdatedAt = types.Int64Value(group.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
