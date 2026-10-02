package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource              = &OrgAPIKeyResource{}
	_ resource.ResourceWithConfigure = &OrgAPIKeyResource{}
)

type OrgAPIKeyResource struct {
	client *client.Client
}

type OrgAPIKeyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Rotate         types.String `tfsdk:"rotate"`
	ClientID       types.String `tfsdk:"client_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Comment        types.String `tfsdk:"comment"`
	Token          types.String `tfsdk:"token"`
	Type           types.String `tfsdk:"type"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	Counter        types.Int64  `tfsdk:"counter"`
	UsageToday     types.Int64  `tfsdk:"usage_today"`
	UsageLimit     types.Int64  `tfsdk:"usage_limit"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	LastUsedAt     types.Int64  `tfsdk:"last_used_at"`
	LastUsedIP     types.String `tfsdk:"last_used_ip"`
	LastUsedUA     types.String `tfsdk:"last_used_ua"`
}

func NewOrgAPIKeyResource() resource.Resource {
	return &OrgAPIKeyResource{}
}

func (r *OrgAPIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_api_key"
}

func (r *OrgAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the API key of the organization tied to the configured organization token. " +
			"The key is adopted on create and removed on destroy. Change the `rotate` value to rotate the key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the organization API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rotate": schema.StringAttribute{
				Description: "Arbitrary keeper value. Changing it rotates the organization API key.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client UUID that owns the key.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization UUID the key belongs to.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment associated with the key.",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The organization API key token.",
				Computed:    true,
				Sensitive:   true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the key.",
				Computed:    true,
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the key is inactive.",
				Computed:    true,
			},
			"counter": schema.Int64Attribute{
				Description: "Total usage counter for the key.",
				Computed:    true,
			},
			"usage_today": schema.Int64Attribute{
				Description: "Number of times the key was used today.",
				Computed:    true,
			},
			"usage_limit": schema.Int64Attribute{
				Description: "Daily usage limit for the key.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the key was created.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Identifier of who created the key.",
				Computed:    true,
			},
			"last_used_at": schema.Int64Attribute{
				Description: "Unix timestamp when the key was last used.",
				Computed:    true,
			},
			"last_used_ip": schema.StringAttribute{
				Description: "IP address that last used the key.",
				Computed:    true,
			},
			"last_used_ua": schema.StringAttribute{
				Description: "User agent that last used the key.",
				Computed:    true,
			},
		},
	}
}

func (r *OrgAPIKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func mapOrgAPIKey(k *client.APIKey, m *OrgAPIKeyResourceModel) {
	m.ID = types.StringValue(k.ID)
	m.ClientID = types.StringValue(k.ClientID)
	m.OrganizationID = types.StringValue(k.OrganizationID)
	m.Comment = types.StringValue(k.Comment)
	m.Token = types.StringValue(k.Token)
	m.Type = types.StringValue(k.Type)
	m.Inactive = types.BoolValue(k.Inactive)
	m.Counter = types.Int64Value(k.Counter)
	m.UsageToday = types.Int64Value(k.UsageToday)
	m.UsageLimit = types.Int64Value(k.UsageLimit)
	m.CreatedAt = types.Int64Value(k.CreatedAt)
	m.CreatedBy = types.StringValue(k.CreatedBy)
	m.LastUsedAt = types.Int64Value(k.LastUsedAt)
	m.LastUsedIP = types.StringValue(k.LastUsedIP)
	m.LastUsedUA = types.StringValue(k.LastUsedUA)
}

func (r *OrgAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrgAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetOrgKey(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading org API key", err.Error())
		return
	}

	mapOrgAPIKey(key, &plan)
	tflog.Trace(ctx, "adopted org API key", map[string]interface{}{"id": key.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OrgAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrgAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetOrgKey(ctx)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading org API key", err.Error())
		return
	}

	mapOrgAPIKey(key, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *OrgAPIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state OrgAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Rotate only when the keeper value changes; otherwise refresh from the API.
	if !plan.Rotate.Equal(state.Rotate) {
		key, err := r.client.RotateOrgKey(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error rotating org API key", err.Error())
			return
		}
		mapOrgAPIKey(key, &plan)
	} else {
		key, err := r.client.GetOrgKey(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error reading org API key", err.Error())
			return
		}
		mapOrgAPIKey(key, &plan)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OrgAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrgAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteOrgKey(ctx); err != nil {
		resp.Diagnostics.AddError("Error deleting org API key", err.Error())
		return
	}
}
