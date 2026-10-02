package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	_ resource.Resource                = &ExportTokenResource{}
	_ resource.ResourceWithImportState = &ExportTokenResource{}
	_ resource.ResourceWithConfigure   = &ExportTokenResource{}
)

type ExportTokenResource struct {
	client *client.Client
}

type ExportTokenResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Rotate     types.String `tfsdk:"rotate"`
	Token      types.String `tfsdk:"token"`
	Comment    types.String `tfsdk:"comment"`
	OrgID      types.String `tfsdk:"org_id"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	LastUsedAt types.Int64  `tfsdk:"last_used_at"`
	LastUsedBy types.String `tfsdk:"last_used_by"`
	Counter    types.Int64  `tfsdk:"counter"`
}

func NewExportTokenResource() resource.Resource {
	return &ExportTokenResource{}
}

func (r *ExportTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export_token"
}

func (r *ExportTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero export token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the export token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token": schema.StringAttribute{
				Description: "The export token value.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the export token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rotate": schema.StringAttribute{
				Description: "Arbitrary keeper value. Changing it rotates the export token.",
				Optional:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The organization ID for this export token.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the export token was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_used_at": schema.Int64Attribute{
				Description: "Unix timestamp when the token was last used.",
				Computed:    true,
			},
			"last_used_by": schema.StringAttribute{
				Description: "Identifier of the last user of the token.",
				Computed:    true,
			},
			"counter": schema.Int64Attribute{
				Description: "Total usage counter.",
				Computed:    true,
			},
		},
	}
}

func (r *ExportTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExportTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ExportTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.CreateExportToken(ctx, plan.OrgID.ValueString(), client.CreateExportTokenRequest{
		Comment: plan.Comment.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating export token", err.Error())
		return
	}

	mapExportTokenToState(token, &plan)

	tflog.Trace(ctx, "created export token", map[string]interface{}{"id": token.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ExportTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ExportTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetExportToken(ctx, state.OrgID.ValueString(), state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading export token", err.Error())
		return
	}

	mapExportTokenToState(token, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ExportTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ExportTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Rotate.Equal(state.Rotate) {
		token, err := r.client.RotateExportToken(ctx, state.OrgID.ValueString(), state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error rotating export token", err.Error())
			return
		}
		mapExportTokenToState(token, &plan)
	} else {
		token, err := r.client.GetExportToken(ctx, state.OrgID.ValueString(), state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error reading export token", err.Error())
			return
		}
		mapExportTokenToState(token, &plan)
	}
	plan.OrgID = state.OrgID

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ExportTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ExportTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteExportToken(ctx, state.OrgID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting export token", err.Error())
		return
	}
}

func (r *ExportTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: org_id/token_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), types.StringValue(parts[0]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(parts[1]))...)
}

func mapExportTokenToState(token *client.ExportToken, state *ExportTokenResourceModel) {
	state.ID = types.StringValue(token.ID)
	state.Token = types.StringValue(token.Token)
	state.Comment = types.StringValue(token.Comment)
	state.CreatedAt = types.Int64Value(token.CreatedAt)
	state.LastUsedAt = types.Int64Value(token.LastUsedAt)
	state.LastUsedBy = types.StringValue(token.LastUsedBy)
	state.Counter = types.Int64Value(token.Counter)
}
