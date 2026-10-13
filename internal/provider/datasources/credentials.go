package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &CredentialsDataSource{}
	_ datasource.DataSourceWithConfigure = &CredentialsDataSource{}
)

type CredentialsDataSource struct {
	client *client.Client
}

type CredentialsDataSourceModel struct {
	Credentials []CredentialDSModel `tfsdk:"credentials"`
}

type CredentialDSModel struct {
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

func NewCredentialsDataSource() datasource.DataSource {
	return &CredentialsDataSource{}
}

func (d *CredentialsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credentials"
}

func (d *CredentialsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of RunZero credentials.",
		Attributes: map[string]schema.Attribute{
			"credentials": schema.ListNestedAttribute{
				Description: "List of credentials.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: credentialDSSchemaAttributes(false),
				},
			},
		},
	}
}

func (d *CredentialsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state CredentialsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentials, err := d.client.ListCredentials(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading credentials", err.Error())
		return
	}

	state.Credentials = make([]CredentialDSModel, len(credentials))
	for i, c := range credentials {
		var acl types.Map
		if c.ACL != nil {
			aclStrMap := make(map[string]string, len(c.ACL))
			for k, v := range c.ACL {
				aclStrMap[k] = fmt.Sprintf("%v", v)
			}
			aclVal, d := types.MapValueFrom(ctx, types.StringType, aclStrMap)
			resp.Diagnostics.Append(d...)
			acl = aclVal
		} else {
			acl = types.MapNull(types.StringType)
		}

		var cidrs types.List
		if c.CIDRs != nil {
			cidrsVal, d := types.ListValueFrom(ctx, types.StringType, c.CIDRs)
			resp.Diagnostics.Append(d...)
			cidrs = cidrsVal
		} else {
			cidrs = types.ListNull(types.StringType)
		}

		state.Credentials[i] = CredentialDSModel{
			ID:             types.StringValue(c.ID),
			ClientID:       types.StringValue(c.ClientID),
			Name:           types.StringValue(c.Name),
			Type:           types.StringValue(c.Type),
			CreatedAt:      types.Int64Value(c.CreatedAt),
			CreatedByID:    types.StringValue(c.CreatedByID),
			CreatedByEmail: types.StringValue(c.CreatedByEmail),
			ACL:            acl,
			Global:         types.BoolValue(c.Global),
			CIDRs:          cidrs,
			LastUsedAt:     types.Int64Value(c.LastUsedAt),
			LastUsedByID:   types.StringValue(c.LastUsedByID),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func credentialDSSchemaAttributes(idRequired bool) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The UUID of the credential.",
			Computed:    !idRequired,
			Required:    idRequired,
		},
		"client_id": schema.StringAttribute{
			Description: "The client UUID.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the credential.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the credential.",
			Computed:    true,
		},
		"created_at": schema.Int64Attribute{
			Description: "Unix timestamp when the credential was created.",
			Computed:    true,
		},
		"created_by_id": schema.StringAttribute{
			Description: "UUID of the user who created the credential.",
			Computed:    true,
		},
		"created_by_email": schema.StringAttribute{
			Description: "Email of the user who created the credential.",
			Computed:    true,
		},
		"acl": schema.MapAttribute{
			Description: "Access control list mapping organization UUIDs to permission levels.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"global": schema.BoolAttribute{
			Description: "Whether the credential is globally available.",
			Computed:    true,
		},
		"cidrs": schema.ListAttribute{
			Description: "List of allowed CIDR ranges.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"last_used_at": schema.Int64Attribute{
			Description: "Unix timestamp when the credential was last used.",
			Computed:    true,
		},
		"last_used_by_id": schema.StringAttribute{
			Description: "UUID of the user who last used the credential.",
			Computed:    true,
		},
	}
	return attrs
}

func mapCredentialToDS(ctx context.Context, c *client.Credential, state *CredentialDataSourceModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(c.ID)
	state.ClientID = types.StringValue(c.ClientID)
	state.Name = types.StringValue(c.Name)
	state.Type = types.StringValue(c.Type)
	state.CreatedAt = types.Int64Value(c.CreatedAt)
	state.CreatedByID = types.StringValue(c.CreatedByID)
	state.CreatedByEmail = types.StringValue(c.CreatedByEmail)
	state.Global = types.BoolValue(c.Global)
	state.LastUsedAt = types.Int64Value(c.LastUsedAt)
	state.LastUsedByID = types.StringValue(c.LastUsedByID)

	if c.ACL != nil {
		aclStrMap := make(map[string]string, len(c.ACL))
		for k, v := range c.ACL {
			aclStrMap[k] = fmt.Sprintf("%v", v)
		}
		acl, d := types.MapValueFrom(ctx, types.StringType, aclStrMap)
		diags.Append(d...)
		state.ACL = acl
	} else {
		state.ACL = types.MapNull(types.StringType)
	}

	if c.CIDRs != nil {
		cidrs, d := types.ListValueFrom(ctx, types.StringType, c.CIDRs)
		diags.Append(d...)
		state.CIDRs = cidrs
	} else {
		state.CIDRs = types.ListNull(types.StringType)
	}
}
