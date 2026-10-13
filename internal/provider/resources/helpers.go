package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func extractMapString(ctx context.Context, m types.Map, diags *diag.Diagnostics) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	elements := make(map[string]types.String, len(m.Elements()))
	d := m.ElementsAs(ctx, &elements, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil
	}
	result := make(map[string]string, len(elements))
	for k, v := range elements {
		result[k] = v.ValueString()
	}
	return result
}

func filterMapToKeys(_ context.Context, full types.Map, reference types.Map, diags *diag.Diagnostics) types.Map {
	if full.IsNull() || full.IsUnknown() {
		return full
	}
	refElements := reference.Elements()
	fullElements := full.Elements()
	filtered := make(map[string]attr.Value, len(refElements))
	for k := range refElements {
		if v, ok := fullElements[k]; ok {
			filtered[k] = v
		}
	}
	result, d := types.MapValue(types.StringType, filtered)
	diags.Append(d...)
	return result
}
