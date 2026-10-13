package client

import (
	"context"
	"fmt"
	"net/http"
)

type CreateScanRequest struct {
	Targets              string `json:"targets"`
	ScanName             string `json:"scan-name,omitempty"`
	ScanDescription      string `json:"scan-description,omitempty"`
	ScanTemplate         string `json:"scan-template,omitempty"`
	ScanFrequency        string `json:"scan-frequency,omitempty"`
	ScanStart            string `json:"scan-start,omitempty"`
	ScanTags             string `json:"scan-tags,omitempty"`
	ScanGracePeriod      string `json:"scan-grace-period,omitempty"`
	Agent                string `json:"agent,omitempty"`
	Explorer             string `json:"explorer,omitempty"`
	ExplorerGroupID      string `json:"explorer-group-id,omitempty"`
	HostedZoneID         string `json:"hosted-zone-id,omitempty"`
	HostedZoneName       string `json:"hosted-zone-name,omitempty"`
	Rate                 string `json:"rate,omitempty"`
	MaxHostRate          string `json:"max-host-rate,omitempty"`
	Passes               string `json:"passes,omitempty"`
	MaxAttempts          string `json:"max-attempts,omitempty"`
	MaxSockets           string `json:"max-sockets,omitempty"`
	MaxGroupSize         string `json:"max-group-size,omitempty"`
	MaxTTL               string `json:"max-ttl,omitempty"`
	TCPPorts             string `json:"tcp-ports,omitempty"`
	TCPExcludes          string `json:"tcp-excludes,omitempty"`
	Excludes             string `json:"excludes,omitempty"`
	Probes               string `json:"probes,omitempty"`
	Screenshots          string `json:"screenshots,omitempty"`
	Nameservers          string `json:"nameservers,omitempty"`
	SubnetPing           string `json:"subnet-ping,omitempty"`
	SubnetPingNetSize    string `json:"subnet-ping-net-size,omitempty"`
	SubnetPingProbes     string `json:"subnet-ping-probes,omitempty"`
	SubnetPingSampleRate string `json:"subnet-ping-sample-rate,omitempty"`
	HostPing             string `json:"host-ping,omitempty"`
	HostPingProbes       string `json:"host-ping-probes,omitempty"`
	TOS                  string `json:"tos,omitempty"`
}

type Scan struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	TemplateID          string            `json:"template_id"`
	ClientID            string            `json:"client_id"`
	OrganizationID      string            `json:"organization_id"`
	SiteID              string            `json:"site_id"`
	AgentID             string            `json:"agent_id"`
	HostedZoneID        string            `json:"hosted_zone_id"`
	CruncherID          string            `json:"cruncher_id"`
	Type                string            `json:"type"`
	Status              string            `json:"status"`
	Error               string            `json:"error"`
	Params              map[string]string `json:"params"`
	Stats               map[string]any    `json:"stats"`
	Hidden              bool              `json:"hidden"`
	ParentID            string            `json:"parent_id"`
	Recur               bool              `json:"recur"`
	RecurFrequency      string            `json:"recur_frequency"`
	CreatedAt           int64             `json:"created_at"`
	CreatedBy           string            `json:"created_by"`
	CreatedByUserID     string            `json:"created_by_user_id"`
	CustomIntegrationID string            `json:"custom_integration_id"`
	SourceID            int64             `json:"source_id"`
	UpdatedAt           int64             `json:"updated_at"`
	StartTime           int64             `json:"start_time"`
	RecurLast           int64             `json:"recur_last"`
	RecurNext           int64             `json:"recur_next"`
	RecurLastTaskID     string            `json:"recur_last_task_id"`
}

func (c *Client) CreateScan(ctx context.Context, siteID string, req CreateScanRequest) (*Scan, error) {
	resp, err := c.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/org/sites/%s/scan", siteID), AuthOrg, req)
	if err != nil {
		return nil, err
	}
	var scan Scan
	if err := HandleResponse(resp, &scan); err != nil {
		return nil, fmt.Errorf("create scan: %w", err)
	}
	return &scan, nil
}
