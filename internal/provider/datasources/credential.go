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
	_ datasource.DataSource              = &CredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &CredentialDataSource{}
)

type CredentialDataSource struct {
	client *client.Client
}

type CredentialDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	ClientID       types.String `tfsdk:"client_id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedByID    types.String `tfsdk:"created_by_id"`
	CreatedByEmail types.String `tfsdk:"created_by_email"`
	ACL            types.Map    `tfsdk:"acl"`
	Global         types.Bool   `tfsdk:"global"`
	CIDRs          types.List   `tfsdk:"cidrs"`
	LastUsedAt     types.Int64  `tfsdk:"last_used_at"`
	LastUsedByID   types.String `tfsdk:"last_used_by_id"`
}

func NewCredentialDataSource() datasource.DataSource {
	return &CredentialDataSource{}
}

func (d *CredentialDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (d *CredentialDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single RunZero credential by ID.",
		Attributes:  credentialDSSchemaAttributes(true),
	}
}

func (d *CredentialDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CredentialDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credential, err := d.client.GetCredential(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading credential", err.Error())
		return
	}

	mapCredentialToDS(ctx, credential, &state, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
