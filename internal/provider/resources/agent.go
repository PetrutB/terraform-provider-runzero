package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &AgentResource{}
	_ resource.ResourceWithImportState = &AgentResource{}
	_ resource.ResourceWithConfigure   = &AgentResource{}
)

type AgentResource struct {
	client *client.Client
}

type AgentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	SiteID         types.String `tfsdk:"site_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ClientID       types.String `tfsdk:"client_id"`
	OS             types.String `tfsdk:"os"`
	Arch           types.String `tfsdk:"arch"`
	Version        types.String `tfsdk:"version"`
	ExternalIP     types.String `tfsdk:"external_ip"`
	InternalIP     types.String `tfsdk:"internal_ip"`
	Connected      types.Bool   `tfsdk:"connected"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	LastCheckin    types.Int64  `tfsdk:"last_checkin"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	UpdatedAt      types.Int64  `tfsdk:"updated_at"`
	DeactivatedAt  types.Int64  `tfsdk:"deactivated_at"`
}

func NewAgentResource() resource.Resource {
	return &AgentResource{}
}

func (r *AgentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (r *AgentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero agent. Agents are deployed externally and adopted into Terraform by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the agent.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the agent.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID the agent is associated with.",
				Optional:    true,
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID the agent belongs to.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"os": schema.StringAttribute{
				Description: "The operating system of the agent host.",
				Computed:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the agent host.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of the agent.",
				Computed:    true,
			},
			"external_ip": schema.StringAttribute{
				Description: "The external IP address of the agent.",
				Computed:    true,
			},
			"internal_ip": schema.StringAttribute{
				Description: "The internal IP address of the agent.",
				Computed:    true,
			},
			"connected": schema.BoolAttribute{
				Description: "Whether the agent is currently connected.",
				Computed:    true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the agent is inactive.",
				Computed:    true,
			},
			"last_checkin": schema.Int64Attribute{
				Description: "Unix timestamp of the agent's last check-in.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the agent was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the agent was last updated.",
				Computed:    true,
			},
			"deactivated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the agent was deactivated.",
				Computed:    true,
			},
		},
	}
}

func (r *AgentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := r.client.GetAgent(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error adopting agent", fmt.Sprintf("Agent %s not found: %s", plan.ID.ValueString(), err.Error()))
		return
	}

	if !plan.SiteID.IsNull() && !plan.SiteID.IsUnknown() && plan.SiteID.ValueString() != agent.SiteID {
		agent, err = r.client.UpdateAgent(ctx, agent.ID, client.UpdateAgentRequest{
			SiteID: plan.SiteID.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating agent site", err.Error())
			return
		}
	}

	r.setStateFromAgent(&plan, agent)

	tflog.Trace(ctx, "adopted agent", map[string]interface{}{"id": agent.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := r.client.GetAgent(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading agent", err.Error())
		return
	}

	r.setStateFromAgent(&state, agent)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateAgentRequest{}
	if !plan.SiteID.IsNull() && plan.SiteID.ValueString() != state.SiteID.ValueString() {
		updateReq.SiteID = plan.SiteID.ValueString()
	}

	agent, err := r.client.UpdateAgent(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating agent", err.Error())
		return
	}

	r.setStateFromAgent(&plan, agent)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *AgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAgent(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting agent", err.Error())
		return
	}
}

func (r *AgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AgentResource) setStateFromAgent(model *AgentResourceModel, agent *client.Agent) {
	model.ID = types.StringValue(agent.ID)
	model.Name = types.StringValue(agent.Name)
	model.SiteID = types.StringValue(agent.SiteID)
	model.OrganizationID = types.StringValue(agent.OrganizationID)
	model.ClientID = types.StringValue(agent.ClientID)
	model.OS = types.StringValue(agent.OS)
	model.Arch = types.StringValue(agent.Arch)
	model.Version = types.StringValue(agent.Version)
	model.ExternalIP = types.StringValue(agent.ExternalIP)
	model.InternalIP = types.StringValue(agent.InternalIP)
	model.Connected = types.BoolValue(agent.Connected)
	model.Inactive = types.BoolValue(agent.Inactive)
	model.LastCheckin = types.Int64Value(agent.LastCheckin)
	model.CreatedAt = types.Int64Value(agent.CreatedAt)
	model.UpdatedAt = types.Int64Value(agent.UpdatedAt)
	model.DeactivatedAt = types.Int64Value(agent.DeactivatedAt)
}
