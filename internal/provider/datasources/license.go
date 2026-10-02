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
	_ datasource.DataSource              = &LicenseDataSource{}
	_ datasource.DataSourceWithConfigure = &LicenseDataSource{}
)

type LicenseDataSource struct {
	client *client.Client
}

type LicenseDataSourceModel struct {
	ID                       types.String `tfsdk:"id"`
	CreatedAt                types.Int64  `tfsdk:"created_at"`
	UpdatedAt                types.Int64  `tfsdk:"updated_at"`
	Name                     types.String `tfsdk:"name"`
	ActivatedAt              types.Int64  `tfsdk:"activated_at"`
	LicenseType              types.String `tfsdk:"license_type"`
	LicenseExpiration        types.Int64  `tfsdk:"license_expiration"`
	LicenseMaxAssets         types.Int64  `tfsdk:"license_max_assets"`
	LicenseLiveAssetCount    types.Int64  `tfsdk:"license_live_asset_count"`
	LicenseProjectAssetCount types.Int64  `tfsdk:"license_project_asset_count"`
	ViaReseller              types.Bool   `tfsdk:"via_reseller"`
	Phone                    types.String `tfsdk:"phone"`
	AddressStreet            types.String `tfsdk:"address_street"`
	AddressCountry           types.String `tfsdk:"address_country"`
	AddressCity              types.String `tfsdk:"address_city"`
	AddressRegion            types.String `tfsdk:"address_region"`
	AddressPostal            types.String `tfsdk:"address_postal"`
	SubscriptionPeriodStart  types.Int64  `tfsdk:"subscription_period_start"`
	SubscriptionPeriodEnd    types.Int64  `tfsdk:"subscription_period_end"`
	SubscriptionCancelAt     types.Int64  `tfsdk:"subscription_cancel_at"`
	SubscriptionCanceledAt   types.Int64  `tfsdk:"subscription_canceled_at"`
	SSOType                  types.String `tfsdk:"sso_type"`
	SSODomain                types.String `tfsdk:"sso_domain"`
	SSOMode                  types.String `tfsdk:"sso_mode"`
	SSOLoginMessage          types.String `tfsdk:"sso_login_message"`
	SSOLoginIssuerURL        types.String `tfsdk:"sso_login_issuer_url"`
	SSOLoginLoginURL         types.String `tfsdk:"sso_login_login_url"`
	SSOLoginLogoutURL        types.String `tfsdk:"sso_login_logout_url"`
	SSODefaultRole           types.String `tfsdk:"sso_default_role"`
	Partner                  types.String `tfsdk:"partner"`
}

func NewLicenseDataSource() datasource.DataSource {
	return &LicenseDataSource{}
}

func (d *LicenseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

func (d *LicenseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the RunZero account license details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the license.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the license was created.",
				Computed:    true,
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the license was last updated.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name on the license.",
				Computed:    true,
			},
			"activated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the license was activated.",
				Computed:    true,
			},
			"license_type": schema.StringAttribute{
				Description: "The type of license.",
				Computed:    true,
			},
			"license_expiration": schema.Int64Attribute{
				Description: "Unix timestamp when the license expires.",
				Computed:    true,
			},
			"license_max_assets": schema.Int64Attribute{
				Description: "Maximum number of assets allowed by the license.",
				Computed:    true,
			},
			"license_live_asset_count": schema.Int64Attribute{
				Description: "Current count of live assets.",
				Computed:    true,
			},
			"license_project_asset_count": schema.Int64Attribute{
				Description: "Current count of project assets.",
				Computed:    true,
			},
			"via_reseller": schema.BoolAttribute{
				Description: "Whether the license was obtained via a reseller.",
				Computed:    true,
			},
			"phone": schema.StringAttribute{
				Description: "Contact phone number.",
				Computed:    true,
			},
			"address_street": schema.StringAttribute{
				Description: "Street address.",
				Computed:    true,
			},
			"address_country": schema.StringAttribute{
				Description: "Country.",
				Computed:    true,
			},
			"address_city": schema.StringAttribute{
				Description: "City.",
				Computed:    true,
			},
			"address_region": schema.StringAttribute{
				Description: "Region or state.",
				Computed:    true,
			},
			"address_postal": schema.StringAttribute{
				Description: "Postal or ZIP code.",
				Computed:    true,
			},
			"subscription_period_start": schema.Int64Attribute{
				Description: "Unix timestamp of the subscription period start.",
				Computed:    true,
			},
			"subscription_period_end": schema.Int64Attribute{
				Description: "Unix timestamp of the subscription period end.",
				Computed:    true,
			},
			"subscription_cancel_at": schema.Int64Attribute{
				Description: "Unix timestamp when the subscription is scheduled to cancel.",
				Computed:    true,
			},
			"subscription_canceled_at": schema.Int64Attribute{
				Description: "Unix timestamp when the subscription was canceled.",
				Computed:    true,
			},
			"sso_type": schema.StringAttribute{
				Description: "SSO provider type.",
				Computed:    true,
			},
			"sso_domain": schema.StringAttribute{
				Description: "SSO domain.",
				Computed:    true,
			},
			"sso_mode": schema.StringAttribute{
				Description: "SSO mode.",
				Computed:    true,
			},
			"sso_login_message": schema.StringAttribute{
				Description: "SSO login message displayed to users.",
				Computed:    true,
			},
			"sso_login_issuer_url": schema.StringAttribute{
				Description: "SSO login issuer URL.",
				Computed:    true,
			},
			"sso_login_login_url": schema.StringAttribute{
				Description: "SSO login URL.",
				Computed:    true,
			},
			"sso_login_logout_url": schema.StringAttribute{
				Description: "SSO logout URL.",
				Computed:    true,
			},
			"sso_default_role": schema.StringAttribute{
				Description: "Default role assigned to SSO users.",
				Computed:    true,
			},
			"partner": schema.StringAttribute{
				Description: "Partner identifier.",
				Computed:    true,
			},
		},
	}
}

func (d *LicenseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LicenseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	license, err := d.client.GetLicense(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading license", err.Error())
		return
	}

	state := LicenseDataSourceModel{
		ID:                       types.StringValue(license.ID),
		CreatedAt:                types.Int64Value(license.CreatedAt),
		UpdatedAt:                types.Int64Value(license.UpdatedAt),
		Name:                     types.StringValue(license.Name),
		ActivatedAt:              types.Int64Value(license.ActivatedAt),
		LicenseType:              types.StringValue(license.LicenseType),
		LicenseExpiration:        types.Int64Value(license.LicenseExpiration),
		LicenseMaxAssets:         types.Int64Value(license.LicenseMaxAssets),
		LicenseLiveAssetCount:    types.Int64Value(license.LicenseLiveAssetCount),
		LicenseProjectAssetCount: types.Int64Value(license.LicenseProjectAssetCount),
		ViaReseller:              types.BoolValue(license.ViaReseller),
		Phone:                    types.StringValue(license.Phone),
		AddressStreet:            types.StringValue(license.AddressStreet),
		AddressCountry:           types.StringValue(license.AddressCountry),
		AddressCity:              types.StringValue(license.AddressCity),
		AddressRegion:            types.StringValue(license.AddressRegion),
		AddressPostal:            types.StringValue(license.AddressPostal),
		SubscriptionPeriodStart:  types.Int64Value(license.SubscriptionPeriodStart),
		SubscriptionPeriodEnd:    types.Int64Value(license.SubscriptionPeriodEnd),
		SubscriptionCancelAt:     types.Int64Value(license.SubscriptionCancelAt),
		SubscriptionCanceledAt:   types.Int64Value(license.SubscriptionCanceledAt),
		SSOType:                  types.StringValue(license.SSOType),
		SSODomain:                types.StringValue(license.SSODomain),
		SSOMode:                  types.StringValue(license.SSOMode),
		SSOLoginMessage:          types.StringValue(license.SSOLoginMessage),
		SSOLoginIssuerURL:        types.StringValue(license.SSOLoginIssuerURL),
		SSOLoginLoginURL:         types.StringValue(license.SSOLoginLoginURL),
		SSOLoginLogoutURL:        types.StringValue(license.SSOLoginLogoutURL),
		SSODefaultRole:           types.StringValue(license.SSODefaultRole),
		Partner:                  types.StringValue(license.Partner),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
