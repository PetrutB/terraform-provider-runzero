package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource              = &ScanResource{}
	_ resource.ResourceWithConfigure = &ScanResource{}
)

type ScanResource struct {
	client *client.Client
}

type ScanResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	SiteID               types.String `tfsdk:"site_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Targets              types.String `tfsdk:"targets"`
	Excludes             types.String `tfsdk:"excludes"`
	Agent                types.String `tfsdk:"agent"`
	Explorer             types.String `tfsdk:"explorer"`
	ScanTemplate         types.String `tfsdk:"scan_template"`
	ScanFrequency        types.String `tfsdk:"scan_frequency"`
	Rate                 types.String `tfsdk:"rate"`
	MaxHostRate          types.String `tfsdk:"max_host_rate"`
	TCPPorts             types.String `tfsdk:"tcp_ports"`
	Probes               types.String `tfsdk:"probes"`
	ScanStart            types.String `tfsdk:"scan_start"`
	ScanTags             types.String `tfsdk:"scan_tags"`
	ScanGracePeriod      types.String `tfsdk:"scan_grace_period"`
	ExplorerGroupID      types.String `tfsdk:"explorer_group_id"`
	HostedZoneID         types.String `tfsdk:"hosted_zone_id"`
	HostedZoneName       types.String `tfsdk:"hosted_zone_name"`
	Passes               types.String `tfsdk:"passes"`
	MaxAttempts          types.String `tfsdk:"max_attempts"`
	MaxSockets           types.String `tfsdk:"max_sockets"`
	MaxGroupSize         types.String `tfsdk:"max_group_size"`
	MaxTTL               types.String `tfsdk:"max_ttl"`
	TCPExcludes          types.String `tfsdk:"tcp_excludes"`
	Screenshots          types.String `tfsdk:"screenshots"`
	Nameservers          types.String `tfsdk:"nameservers"`
	SubnetPing           types.String `tfsdk:"subnet_ping"`
	SubnetPingNetSize    types.String `tfsdk:"subnet_ping_net_size"`
	SubnetPingProbes     types.String `tfsdk:"subnet_ping_probes"`
	SubnetPingSampleRate types.String `tfsdk:"subnet_ping_sample_rate"`
	HostPing             types.String `tfsdk:"host_ping"`
	HostPingProbes       types.String `tfsdk:"host_ping_probes"`
	TOS                  types.String `tfsdk:"tos"`
	Status               types.String `tfsdk:"status"`
	CreatedAt            types.Int64  `tfsdk:"created_at"`
	UpdatedAt            types.Int64  `tfsdk:"updated_at"`
}

func NewScanResource() resource.Resource {
	return &ScanResource{}
}

func (r *ScanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scan"
}

func (r *ScanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Triggers a RunZero scan on a site. Scans cannot be updated; changing any attribute forces a new scan.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the scan task.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID to scan.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the scan.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the scan.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"targets": schema.StringAttribute{
				Description: "Scan targets as a string (IP addresses, CIDR ranges, hostnames, separated by spaces or commas).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"excludes": schema.StringAttribute{
				Description: "Targets to exclude from the scan.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"agent": schema.StringAttribute{
				Description: "The UUID of the agent to run the scan.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"explorer": schema.StringAttribute{
				Description: "The UUID of the explorer to run the scan.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scan_template": schema.StringAttribute{
				Description: "The UUID of the scan template to use.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scan_frequency": schema.StringAttribute{
				Description: "Scan frequency: once, hourly, daily, weekly, monthly, continuous.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rate": schema.StringAttribute{
				Description: "Scan rate (packets per second).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_host_rate": schema.StringAttribute{
				Description: "Maximum rate per host.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tcp_ports": schema.StringAttribute{
				Description: "TCP ports to scan (e.g. 1-1000,5000-6000).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"probes": schema.StringAttribute{
				Description: "Comma-separated list of probes to use.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scan_start":              scanOptionalStringAttr("Scheduled start time for the scan (RFC3339 or relative)."),
			"scan_tags":               scanOptionalStringAttr("Tags to apply to discovered assets."),
			"scan_grace_period":       scanOptionalStringAttr("Grace period before a recurring scan is considered late."),
			"explorer_group_id":       scanOptionalStringAttr("The UUID of the explorer group to run the scan."),
			"hosted_zone_id":          scanOptionalStringAttr("The UUID of the hosted zone to run the scan."),
			"hosted_zone_name":        scanOptionalStringAttr("The name of the hosted zone to run the scan."),
			"passes":                  scanOptionalStringAttr("Number of scan passes."),
			"max_attempts":            scanOptionalStringAttr("Maximum probe attempts per host."),
			"max_sockets":             scanOptionalStringAttr("Maximum concurrent sockets."),
			"max_group_size":          scanOptionalStringAttr("Maximum target group size."),
			"max_ttl":                 scanOptionalStringAttr("Maximum TTL (hop limit) for probes."),
			"tcp_excludes":            scanOptionalStringAttr("TCP ports to exclude from the scan."),
			"screenshots":             scanOptionalStringAttr("Whether to capture screenshots (true/false)."),
			"nameservers":             scanOptionalStringAttr("Comma-separated list of nameservers to use."),
			"subnet_ping":             scanOptionalStringAttr("Whether to perform subnet ping (true/false)."),
			"subnet_ping_net_size":    scanOptionalStringAttr("Subnet size for subnet ping sampling."),
			"subnet_ping_probes":      scanOptionalStringAttr("Probes to use for subnet ping."),
			"subnet_ping_sample_rate": scanOptionalStringAttr("Sample rate for subnet ping."),
			"host_ping":               scanOptionalStringAttr("Whether to perform host ping (true/false)."),
			"host_ping_probes":        scanOptionalStringAttr("Probes to use for host ping."),
			"tos":                     scanOptionalStringAttr("IP type-of-service (ToS) byte for probes."),
			"status": schema.StringAttribute{
				Description: "The status of the scan.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the scan was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the scan was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *ScanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *ScanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scanReq := client.CreateScanRequest{
		Targets:  plan.Targets.ValueString(),
		ScanName: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		scanReq.ScanDescription = plan.Description.ValueString()
	}
	if !plan.Excludes.IsNull() {
		scanReq.Excludes = plan.Excludes.ValueString()
	}
	if !plan.Agent.IsNull() {
		scanReq.Agent = plan.Agent.ValueString()
	}
	if !plan.Explorer.IsNull() {
		scanReq.Explorer = plan.Explorer.ValueString()
	}
	if !plan.ScanTemplate.IsNull() {
		scanReq.ScanTemplate = plan.ScanTemplate.ValueString()
	}
	if !plan.ScanFrequency.IsNull() {
		scanReq.ScanFrequency = plan.ScanFrequency.ValueString()
	}
	if !plan.Rate.IsNull() {
		scanReq.Rate = plan.Rate.ValueString()
	}
	if !plan.MaxHostRate.IsNull() {
		scanReq.MaxHostRate = plan.MaxHostRate.ValueString()
	}
	if !plan.TCPPorts.IsNull() {
		scanReq.TCPPorts = plan.TCPPorts.ValueString()
	}
	if !plan.Probes.IsNull() {
		scanReq.Probes = plan.Probes.ValueString()
	}
	scanReq.ScanStart = plan.ScanStart.ValueString()
	scanReq.ScanTags = plan.ScanTags.ValueString()
	scanReq.ScanGracePeriod = plan.ScanGracePeriod.ValueString()
	scanReq.ExplorerGroupID = plan.ExplorerGroupID.ValueString()
	scanReq.HostedZoneID = plan.HostedZoneID.ValueString()
	scanReq.HostedZoneName = plan.HostedZoneName.ValueString()
	scanReq.Passes = plan.Passes.ValueString()
	scanReq.MaxAttempts = plan.MaxAttempts.ValueString()
	scanReq.MaxSockets = plan.MaxSockets.ValueString()
	scanReq.MaxGroupSize = plan.MaxGroupSize.ValueString()
	scanReq.MaxTTL = plan.MaxTTL.ValueString()
	scanReq.TCPExcludes = plan.TCPExcludes.ValueString()
	scanReq.Screenshots = plan.Screenshots.ValueString()
	scanReq.Nameservers = plan.Nameservers.ValueString()
	scanReq.SubnetPing = plan.SubnetPing.ValueString()
	scanReq.SubnetPingNetSize = plan.SubnetPingNetSize.ValueString()
	scanReq.SubnetPingProbes = plan.SubnetPingProbes.ValueString()
	scanReq.SubnetPingSampleRate = plan.SubnetPingSampleRate.ValueString()
	scanReq.HostPing = plan.HostPing.ValueString()
	scanReq.HostPingProbes = plan.HostPingProbes.ValueString()
	scanReq.TOS = plan.TOS.ValueString()

	scan, err := r.client.CreateScan(ctx, plan.SiteID.ValueString(), scanReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating scan", err.Error())
		return
	}

	plan.ID = types.StringValue(scan.ID)
	plan.Status = types.StringValue(scan.Status)
	plan.CreatedAt = types.Int64Value(scan.CreatedAt)
	plan.UpdatedAt = types.Int64Value(scan.UpdatedAt)

	tflog.Trace(ctx, "created scan", map[string]interface{}{"id": scan.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	task, err := r.client.GetTask(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading scan task", err.Error())
		return
	}

	state.Status = types.StringValue(task.Status)
	state.UpdatedAt = types.Int64Value(task.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ScanResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Scans cannot be updated", "All scan attributes use RequiresReplace. This method should not be called.")
}

func (r *ScanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_ = r.client.StopTask(ctx, state.ID.ValueString())
}

// scanOptionalStringAttr builds an optional, replace-on-change string scan option.
func scanOptionalStringAttr(desc string) schema.StringAttribute {
	return schema.StringAttribute{
		Description: desc,
		Optional:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}
