package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ServerEnrichedMap() planmodifier.Map {
	return serverEnrichedMapModifier{unknownOnCreate: false}
}

func ServerEnrichedMapUnknownOnCreate() planmodifier.Map {
	return serverEnrichedMapModifier{unknownOnCreate: true}
}

type serverEnrichedMapModifier struct {
	unknownOnCreate bool
}

func (m serverEnrichedMapModifier) Description(_ context.Context) string {
	return "Merges user-configured map keys with server-enriched state to prevent drift from API-added defaults."
}

func (m serverEnrichedMapModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m serverEnrichedMapModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.StateValue.IsNull() {
		if m.unknownOnCreate && !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
			resp.PlanValue = types.MapUnknown(types.StringType)
		}
		return
	}

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	stateElements := req.StateValue.Elements()
	configElements := req.ConfigValue.Elements()

	merged := make(map[string]attr.Value, len(stateElements))
	for k, v := range stateElements {
		merged[k] = v
	}
	for k, v := range configElements {
		merged[k] = v
	}

	mergedMap, diags := types.MapValue(types.StringType, merged)
	resp.Diagnostics.Append(diags...)
	if !resp.Diagnostics.HasError() {
		resp.PlanValue = mergedMap
	}
}
