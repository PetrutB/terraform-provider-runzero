package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ resource.Resource                = &CustomIntegrationResource{}
	_ resource.ResourceWithImportState = &CustomIntegrationResource{}
	_ resource.ResourceWithConfigure   = &CustomIntegrationResource{}
)

type CustomIntegrationResource struct {
	client *client.Client
}

type CustomIntegrationResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ClientID    types.String `tfsdk:"client_id"`
	CreatedByID types.String `tfsdk:"created_by_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	DestroyedAt types.String `tfsdk:"destroyed_at"`
	Name        types.String `tfsdk:"name"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}

func NewCustomIntegrationResource() resource.Resource {
	return &CustomIntegrationResource{}
}

func (r *CustomIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_integration"
}

func (r *CustomIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero custom integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the custom integration.",
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
			"created_by_id": schema.StringAttribute{
				Description: "UUID of the user who created the integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the integration was created (RFC 3339).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the integration was last updated (RFC 3339).",
				Computed:    true,
			},
			"destroyed_at": schema.StringAttribute{
				Description: "Timestamp when the integration was destroyed (RFC 3339).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the custom integration.",
				Required:    true,
			},
			"icon": schema.StringAttribute{
				Description: "Base64-encoded icon for the integration.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the custom integration.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *CustomIntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CustomIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.CreateCustomIntegration(ctx, client.CreateCustomIntegrationRequest{
		Name:        plan.Name.ValueString(),
		Icon:        plan.Icon.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating custom integration", err.Error())
		return
	}

	mapCustomIntegrationToState(integration, &plan)

	tflog.Trace(ctx, "created custom integration", map[string]interface{}{"id": integration.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CustomIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.GetCustomIntegration(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading custom integration", err.Error())
		return
	}

	mapCustomIntegrationToState(integration, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CustomIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CustomIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.UpdateCustomIntegration(ctx, state.ID.ValueString(), client.UpdateCustomIntegrationRequest{
		Name:        plan.Name.ValueString(),
		Icon:        plan.Icon.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating custom integration", err.Error())
		return
	}

	mapCustomIntegrationToState(integration, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *CustomIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCustomIntegration(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting custom integration", err.Error())
		return
	}
}

func (r *CustomIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapCustomIntegrationToState(ci *client.CustomIntegration, state *CustomIntegrationResourceModel) {
	state.ID = types.StringValue(ci.ID)
	state.ClientID = types.StringValue(ci.ClientID)
	state.CreatedByID = types.StringValue(ci.CreatedByID)
	state.CreatedAt = types.StringValue(ci.CreatedAt)
	state.UpdatedAt = types.StringValue(ci.UpdatedAt)
	state.DestroyedAt = types.StringValue(ci.DestroyedAt)
	state.Name = types.StringValue(ci.Name)
	state.Icon = types.StringValue(ci.Icon)
	state.Description = types.StringValue(ci.Description)
}
