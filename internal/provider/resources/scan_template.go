package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
	"github.com/petrutb/terraform-provider-runzero/internal/provider/planmodifiers"
)

var (
	_ resource.Resource                = &ScanTemplateResource{}
	_ resource.ResourceWithImportState = &ScanTemplateResource{}
	_ resource.ResourceWithConfigure   = &ScanTemplateResource{}
)

type ScanTemplateResource struct {
	client *client.Client
}

type ScanTemplateResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Params          types.Map    `tfsdk:"params"`
	OrganizationID  types.String `tfsdk:"organization_id"`
	ClientID        types.String `tfsdk:"client_id"`
	AgentID         types.String `tfsdk:"agent_id"`
	SiteID          types.String `tfsdk:"site_id"`
	CreatedAt       types.Int64  `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	UpdatedAt       types.Int64  `tfsdk:"updated_at"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	Hidden          types.Bool   `tfsdk:"hidden"`
	Recur           types.Bool   `tfsdk:"recur"`
	RecurFrequency  types.String `tfsdk:"recur_frequency"`
	StartTime       types.Int64  `tfsdk:"start_time"`
	RecurLast       types.Int64  `tfsdk:"recur_last"`
	RecurNext       types.Int64  `tfsdk:"recur_next"`
	GracePeriod     types.String `tfsdk:"grace_period"`
	Global          types.Bool   `tfsdk:"global"`
	ACL             types.Map    `tfsdk:"acl"`
	LinkedTaskCount types.Int64  `tfsdk:"linked_task_count"`
}

func NewScanTemplateResource() resource.Resource {
	return &ScanTemplateResource{}
}

func (r *ScanTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scan_template"
}

func (r *ScanTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RunZero scan template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The UUID of the scan template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the scan template.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the scan template.",
				Optional:    true,
				Computed:    true,
			},
			"params": schema.MapAttribute{
				Description: "Scan parameters for the template.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					planmodifiers.ServerEnrichedMapUnknownOnCreate(),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID for the scan template.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The account ID that owns the template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agent_id": schema.StringAttribute{
				Description: "The explorer ID that may execute the template.",
				Computed:    true,
			},
			"site_id": schema.StringAttribute{
				Description: "The site ID for the template.",
				Computed:    true,
			},
			"created_at": schema.Int64Attribute{
				Description: "Unix timestamp when the scan template was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Description: "The user who created the template.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Description: "Unix timestamp when the scan template was last updated.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of task the template creates.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the last task using the template.",
				Computed:    true,
			},
			"hidden": schema.BoolAttribute{
				Description: "Whether the template is hidden from common view.",
				Optional:    true,
				Computed:    true,
			},
			"recur": schema.BoolAttribute{
				Description: "Whether tasks derived from this template repeat.",
				Optional:    true,
				Computed:    true,
			},
			"recur_frequency": schema.StringAttribute{
				Description: "Execution frequency if recurring (once, hourly, daily, weekly, monthly, continuous).",
				Optional:    true,
				Computed:    true,
			},
			"start_time": schema.Int64Attribute{
				Description: "Unix timestamp for the next execution time.",
				Optional:    true,
				Computed:    true,
			},
			"recur_last": schema.Int64Attribute{
				Description: "Unix timestamp of the last recurring execution.",
				Computed:    true,
			},
			"recur_next": schema.Int64Attribute{
				Description: "Unix timestamp of the next recurring execution.",
				Computed:    true,
			},
			"grace_period": schema.StringAttribute{
				Description: "Additional time beyond hard expiration deadline.",
				Optional:    true,
				Computed:    true,
			},
			"global": schema.BoolAttribute{
				Description: "Whether the template is globally available to all organizations.",
				Required:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"acl": schema.MapAttribute{
				Description: "Access control list mapping organization UUIDs to permission levels.",
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"linked_task_count": schema.Int64Attribute{
				Description: "Number of tasks derived from the template.",
				Computed:    true,
			},
		},
	}
}

func (r *ScanTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func mapScanTemplateToState(ctx context.Context, template *client.ScanTemplate, model *ScanTemplateResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(template.ID)
	model.Name = types.StringValue(template.Name)
	model.Description = types.StringValue(template.Description)
	model.OrganizationID = types.StringValue(template.OrganizationID)
	model.ClientID = types.StringValue(template.ClientID)
	model.AgentID = types.StringValue(template.AgentID)
	model.SiteID = types.StringValue(template.SiteID)
	model.CreatedAt = types.Int64Value(template.CreatedAt)
	model.CreatedBy = types.StringValue(template.CreatedBy)
	model.UpdatedAt = types.Int64Value(template.UpdatedAt)
	model.Type = types.StringValue(template.Type)
	model.Status = types.StringValue(template.Status)
	model.Hidden = types.BoolValue(template.Hidden)
	model.Recur = types.BoolValue(template.Recur)
	model.RecurFrequency = types.StringValue(template.RecurFrequency)
	model.StartTime = types.Int64Value(template.StartTime)
	model.RecurLast = types.Int64Value(template.RecurLast)
	model.RecurNext = types.Int64Value(template.RecurNext)
	model.GracePeriod = types.StringValue(strconv.FormatInt(template.GracePeriod, 10))
	model.Global = types.BoolValue(template.Global)
	model.LinkedTaskCount = types.Int64Value(int64(template.LinkedTaskCount))

	params, d := types.MapValueFrom(ctx, types.StringType, template.Params)
	diags.Append(d...)
	model.Params = params

	aclStrMap := make(map[string]string, len(template.ACL))
	for k, v := range template.ACL {
		aclStrMap[k] = fmt.Sprintf("%v", v)
	}
	acl, d := types.MapValueFrom(ctx, types.StringType, aclStrMap)
	diags.Append(d...)
	model.ACL = acl
}

func (r *ScanTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScanTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMap := make(map[string]string)
	if !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		resp.Diagnostics.Append(plan.Params.ElementsAs(ctx, &paramsMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	aclMap := make(map[string]string)
	if !plan.ACL.IsNull() && !plan.ACL.IsUnknown() {
		resp.Diagnostics.Append(plan.ACL.ElementsAs(ctx, &aclMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	aclAny := make(map[string]any, len(aclMap))
	for k, v := range aclMap {
		aclAny[k] = v
	}

	template, err := r.client.CreateScanTemplate(ctx, client.CreateScanTemplateRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		OrganizationID: plan.OrganizationID.ValueString(),
		Params:         paramsMap,
		Global:         plan.Global.ValueBool(),
		ACL:            aclAny,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating scan template", err.Error())
		return
	}

	mapScanTemplateToState(ctx, template, &plan, &resp.Diagnostics)
	tflog.Trace(ctx, "created scan template", map[string]interface{}{"id": template.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScanTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScanTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := r.client.GetScanTemplate(ctx, state.ID.ValueString())
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading scan template", err.Error())
		return
	}

	mapScanTemplateToState(ctx, template, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ScanTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScanTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paramsMap := make(map[string]string)
	if !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		resp.Diagnostics.Append(plan.Params.ElementsAs(ctx, &paramsMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	aclMap := make(map[string]string)
	if !plan.ACL.IsNull() && !plan.ACL.IsUnknown() {
		resp.Diagnostics.Append(plan.ACL.ElementsAs(ctx, &aclMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	aclAny := make(map[string]any, len(aclMap))
	for k, v := range aclMap {
		aclAny[k] = v
	}

	template, err := r.client.UpdateScanTemplate(ctx, client.UpdateScanTemplateRequest{
		ID:             plan.ID.ValueString(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		OrganizationID: plan.OrganizationID.ValueString(),
		Params:         paramsMap,
		Global:         plan.Global.ValueBool(),
		ACL:            aclAny,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating scan template", err.Error())
		return
	}

	mapScanTemplateToState(ctx, template, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ScanTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScanTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteScanTemplate(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting scan template", err.Error())
		return
	}
}

func (r *ScanTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
