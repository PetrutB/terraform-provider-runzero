package resources

import (
	"context"

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
