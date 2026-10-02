package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
	"github.com/petrutb/terraform-provider-runzero/internal/provider/datasources"
	"github.com/petrutb/terraform-provider-runzero/internal/provider/resources"
)

var _ provider.Provider = &RunZeroProvider{}

type RunZeroProvider struct {
	version string
}

type RunZeroProviderModel struct {
	APIURL       types.String `tfsdk:"api_url"`
	OrgToken     types.String `tfsdk:"org_token"`
	AccountToken types.String `tfsdk:"account_token"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RunZeroProvider{version: version}
	}
}

func (p *RunZeroProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "runzero"
	resp.Version = p.version
}

func (p *RunZeroProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with the RunZero platform API.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Description: "RunZero API base URL. Defaults to https://console.runzero.com/api/v1.0. Can also be set via RUNZERO_API_URL.",
				Optional:    true,
			},
			"org_token": schema.StringAttribute{
				Description: "RunZero Organization API token (OT prefix). Can also be set via RUNZERO_ORG_TOKEN.",
				Optional:    true,
				Sensitive:   true,
			},
			"account_token": schema.StringAttribute{
				Description: "RunZero Account API token (CT prefix). Can also be set via RUNZERO_ACCOUNT_TOKEN.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_id": schema.StringAttribute{
				Description: "RunZero OAuth2 client ID. Can also be set via RUNZERO_CLIENT_ID.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_secret": schema.StringAttribute{
				Description: "RunZero OAuth2 client secret. Can also be set via RUNZERO_CLIENT_SECRET.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *RunZeroProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config RunZeroProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := resolveString(config.APIURL, "RUNZERO_API_URL", "https://console.runzero.com/api/v1.0")
	orgToken := resolveString(config.OrgToken, "RUNZERO_ORG_TOKEN", "")
	accountToken := resolveString(config.AccountToken, "RUNZERO_ACCOUNT_TOKEN", "")
	clientID := resolveString(config.ClientID, "RUNZERO_CLIENT_ID", "")
	clientSecret := resolveString(config.ClientSecret, "RUNZERO_CLIENT_SECRET", "")

	if orgToken == "" && accountToken == "" && (clientID == "" || clientSecret == "") {
		resp.Diagnostics.AddError(
			"Missing RunZero Credentials",
			"At least one of org_token, account_token, or client_id+client_secret must be configured.",
		)
		return
	}

	c, err := client.NewClient(client.Config{
		BaseURL:      apiURL,
		OrgToken:     orgToken,
		AccountToken: accountToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		UserAgent:    "terraform-provider-runzero/" + p.version,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create RunZero client", err.Error())
		return
	}

	resp.ResourceData = c
	resp.DataSourceData = c
}

func resolveString(tfVal types.String, envVar, defaultVal string) string {
	if !tfVal.IsNull() && !tfVal.IsUnknown() {
		return tfVal.ValueString()
	}
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultVal
}

func (p *RunZeroProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewOrganizationResource,
		resources.NewSiteResource,
		resources.NewUserResource,
		resources.NewUserInviteResource,
		resources.NewGroupResource,
		resources.NewAPIKeyResource,
		resources.NewCredentialResource,
		resources.NewCustomIntegrationResource,
		resources.NewScanTemplateResource,
		resources.NewExportTokenResource,
		resources.NewAssetOwnershipTypeResource,
		resources.NewSSOGroupMappingResource,
		resources.NewAgentResource,
		resources.NewExplorerResource,
		resources.NewAssetResource,
		resources.NewScanResource,
		resources.NewOrgAPIKeyResource,
	}
}

func (p *RunZeroProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewAssetsDataSource,
		datasources.NewAssetDataSource,
		datasources.NewAgentsDataSource,
		datasources.NewAgentDataSource,
		datasources.NewServicesDataSource,
		datasources.NewServiceDataSource,
		datasources.NewExplorersDataSource,
		datasources.NewExplorerDataSource,
		datasources.NewTasksDataSource,
		datasources.NewTaskDataSource,
		datasources.NewHostedZonesDataSource,
		datasources.NewHostedZoneDataSource,
		datasources.NewWirelessLANsDataSource,
		datasources.NewWirelessLANDataSource,
		datasources.NewOrganizationDataSource,
		datasources.NewOrganizationsDataSource,
		datasources.NewOrgDataSource,
		datasources.NewSiteDataSource,
		datasources.NewSitesDataSource,
		datasources.NewUserDataSource,
		datasources.NewUsersDataSource,
		datasources.NewGroupDataSource,
		datasources.NewGroupsDataSource,
		datasources.NewAPIKeyDataSource,
		datasources.NewAPIKeysDataSource,
		datasources.NewExportTokenDataSource,
		datasources.NewExportTokensDataSource,
		datasources.NewOrgAPIKeyDataSource,
		datasources.NewScanTemplateDataSource,
		datasources.NewScanTemplatesDataSource,
		datasources.NewSSOGroupMappingDataSource,
		datasources.NewSSOGroupMappingsDataSource,
		datasources.NewAssetOwnershipTypeDataSource,
		datasources.NewAssetOwnershipTypesDataSource,
		datasources.NewCredentialDataSource,
		datasources.NewCredentialsDataSource,
		datasources.NewCustomIntegrationDataSource,
		datasources.NewCustomIntegrationsDataSource,
		datasources.NewLicenseDataSource,
		datasources.NewHealthDataSource,
		datasources.NewAccountSitesDataSource,
		datasources.NewAccountTasksDataSource,
	}
}
