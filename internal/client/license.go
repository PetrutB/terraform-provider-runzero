package client

import (
	"context"
	"fmt"
	"net/http"
)

type License struct {
	ID                       string         `json:"id"`
	CreatedAt                int64          `json:"created_at"`
	UpdatedAt                int64          `json:"updated_at"`
	Name                     string         `json:"name"`
	ActivatedAt              int64          `json:"activated_at"`
	LicenseType              string         `json:"license_type"`
	LicenseExpiration        int64          `json:"license_expiration"`
	LicenseMaxAssets         int64          `json:"license_max_assets"`
	LicenseLiveAssetCount    int64          `json:"license_live_asset_count"`
	LicenseProjectAssetCount int64          `json:"license_project_asset_count"`
	ViaReseller              bool           `json:"via_reseller"`
	Phone                    string         `json:"phone"`
	AddressStreet            string         `json:"address_street"`
	AddressCountry           string         `json:"address_country"`
	AddressCity              string         `json:"address_city"`
	AddressRegion            string         `json:"address_region"`
	AddressPostal            string         `json:"address_postal"`
	SubscriptionPeriodStart  int64          `json:"subscription_period_start"`
	SubscriptionPeriodEnd    int64          `json:"subscription_period_end"`
	SubscriptionCancelAt     int64          `json:"subscription_cancel_at"`
	SubscriptionCanceledAt   int64          `json:"subscription_canceled_at"`
	Settings                 map[string]any `json:"settings"`
	SSOType                  string         `json:"sso_type"`
	SSODomain                string         `json:"sso_domain"`
	SSOMode                  string         `json:"sso_mode"`
	SSOLoginMessage          string         `json:"sso_login_message"`
	SSOLoginIssuerURL        string         `json:"sso_login_issuer_url"`
	SSOLoginLoginURL         string         `json:"sso_login_login_url"`
	SSOLoginLogoutURL        string         `json:"sso_login_logout_url"`
	SSODefaultRole           string         `json:"sso_default_role"`
	Partner                  string         `json:"partner"`
}

func (c *Client) GetLicense(ctx context.Context) (*License, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/account/license", AuthAccount, nil)
	if err != nil {
		return nil, err
	}
	var license License
	if err := HandleResponse(resp, &license); err != nil {
		return nil, fmt.Errorf("get license: %w", err)
	}
	return &license, nil
}
