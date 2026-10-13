package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func ServerEnrichedMap() planmodifier.Map {
	return serverEnrichedMapModifier{}
}

type serverEnrichedMapModifier struct{}

func (m serverEnrichedMapModifier) Description(_ context.Context) string {
	return "Preserves prior state for server-enriched map attributes when config is absent."
}

func (m serverEnrichedMapModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m serverEnrichedMapModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}
}
