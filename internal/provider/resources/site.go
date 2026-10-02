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
	_ resource.Resource                = &SiteResource{}
	_ resource.ResourceWithImportState = &SiteResource{}
	_ resource.ResourceWithConfigure   = &SiteResource{}
)

type SiteResource struct {
	client *client.Client
}

type SiteResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Scope       types.String `tfsdk:"scope"`
	Excludes    types.String `tfsdk:"excludes"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

func (r *SiteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero site.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the site.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the site.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the site.",
				Optional:    true,
				Computed:    true,
			},
			"scope": schema.StringAttribute{
				Description: "The scan scope for this site (CIDR ranges, etc.).",
				Optional:    true,
				Computed:    true,
			},
			"excludes": schema.StringAttribute{
				Description: "Excluded scan targets.",
				Optional:    true,
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the site was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the site was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *SiteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.CreateSite(ctx, client.CreateSiteRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Scope:       plan.Scope.ValueString(),
		Excludes:    plan.Excludes.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating site", err.Error())
		return
	}

	mapSiteToState(site, &plan)
	tflog.Trace(ctx, "created site", map[string]interface{}{"id": site.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.GetSite(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}

	mapSiteToState(site, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.UpdateSite(ctx, state.ID.ValueString(), client.UpdateSiteRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Scope:       plan.Scope.ValueString(),
		Excludes:    plan.Excludes.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating site", err.Error())
		return
	}

	mapSiteToState(site, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSite(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting site", err.Error())
		return
	}
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapSiteToState(site *client.Site, state *SiteResourceModel) {
	state.ID = types.StringValue(site.ID)
	state.Name = types.StringValue(site.Name)
	state.Description = types.StringValue(site.Description)
	state.Scope = types.StringValue(site.Scope)
	state.Excludes = types.StringValue(site.Excludes)
	state.CreatedAt = types.Int64Value(site.CreatedAt)
	state.UpdatedAt = types.Int64Value(site.UpdatedAt)
}
