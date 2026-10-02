package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource                = &APIKeyResource{}
	_ resource.ResourceWithImportState = &APIKeyResource{}
	_ resource.ResourceWithConfigure   = &APIKeyResource{}
)

type APIKeyResource struct {
	client *client.Client
}

type APIKeyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Rotate         types.String `tfsdk:"rotate"`
	ClientID       types.String `tfsdk:"client_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Comment        types.String `tfsdk:"comment"`
	LastUsedAt     types.Int64  `tfsdk:"last_used_at"`
	LastUsedIP     types.String `tfsdk:"last_used_ip"`
	LastUsedUA     types.String `tfsdk:"last_used_ua"`
	Counter        types.Int64  `tfsdk:"counter"`
	UsageToday     types.Int64  `tfsdk:"usage_today"`
	UsageLimit     types.Int64  `tfsdk:"usage_limit"`
	Token          types.String `tfsdk:"token"`
	Inactive       types.Bool   `tfsdk:"inactive"`
	Type           types.String `tfsdk:"type"`
}

func NewAPIKeyResource() resource.Resource {
	return &APIKeyResource{}
}

func (r *APIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *APIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization UUID to scope the key to.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the API key was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Description: "Email of the user who created the key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rotate": schema.StringAttribute{
				Description: "Arbitrary keeper value. Changing it rotates the API key token.",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the API key.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_used_at": schema.Int64Attribute{
				Description: "Unix timestamp when the key was last used.",
				Computed:    true,
			},
			"last_used_ip": schema.StringAttribute{
				Description: "IP address from which the key was last used.",
				Computed:    true,
			},
			"last_used_ua": schema.StringAttribute{
				Description: "User agent string from the last use.",
				Computed:    true,
			},
			"counter": schema.Int64Attribute{
				Description: "Total usage counter.",
				Computed:    true,
			},
			"usage_today": schema.Int64Attribute{
				Description: "Usage count for today.",
				Computed:    true,
			},
			"usage_limit": schema.Int64Attribute{
				Description: "Maximum allowed usage.",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The API key token value.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inactive": schema.BoolAttribute{
				Description: "Whether the API key is inactive.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the API key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *APIKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *APIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan APIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.CreateAPIKey(ctx, client.CreateAPIKeyRequest{
		Comment:        plan.Comment.ValueString(),
		OrganizationID: plan.OrganizationID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating API key", err.Error())
		return
	}

	mapAPIKeyToState(apiKey, &plan)

	tflog.Trace(ctx, "created api key", map[string]interface{}{"id": apiKey.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *APIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.GetAPIKey(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading API key", err.Error())
		return
	}

	mapAPIKeyToState(apiKey, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *APIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state APIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Rotate.Equal(state.Rotate) {
		key, err := r.client.RotateAPIKey(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error rotating API key", err.Error())
			return
		}
		mapAPIKeyToState(key, &plan)
	} else {
		key, err := r.client.GetAPIKey(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error reading API key", err.Error())
			return
		}
		mapAPIKeyToState(key, &plan)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *APIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state APIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAPIKey(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting API key", err.Error())
		return
	}
}

func (r *APIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapAPIKeyToState(apiKey *client.APIKey, state *APIKeyResourceModel) {
	state.ID = types.StringValue(apiKey.ID)
	state.ClientID = types.StringValue(apiKey.ClientID)
	state.OrganizationID = types.StringValue(apiKey.OrganizationID)
	state.CreatedAt = types.Int64Value(apiKey.CreatedAt)
	state.CreatedBy = types.StringValue(apiKey.CreatedBy)
	state.Comment = types.StringValue(apiKey.Comment)
	state.LastUsedAt = types.Int64Value(apiKey.LastUsedAt)
	state.LastUsedIP = types.StringValue(apiKey.LastUsedIP)
	state.LastUsedUA = types.StringValue(apiKey.LastUsedUA)
	state.Counter = types.Int64Value(apiKey.Counter)
	state.UsageToday = types.Int64Value(apiKey.UsageToday)
	state.UsageLimit = types.Int64Value(apiKey.UsageLimit)
	state.Token = types.StringValue(apiKey.Token)
	state.Inactive = types.BoolValue(apiKey.Inactive)
	state.Type = types.StringValue(apiKey.Type)
}
