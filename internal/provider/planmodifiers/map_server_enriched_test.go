package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrutb/terraform-provider-runzero/internal/provider/planmodifiers"
)

func TestServerEnrichedMap_Create(t *testing.T) {
	ctx := context.Background()
	configVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate": "1000",
	})
	req := planmodifier.MapRequest{
		StateValue:  types.MapNull(types.StringType),
		ConfigValue: configVal,
		PlanValue:   configVal,
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(configVal) {
		t.Errorf("expected config passthrough on create, got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMapUnknownOnCreate(t *testing.T) {
	ctx := context.Background()
	configVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate": "1000",
	})
	req := planmodifier.MapRequest{
		StateValue:  types.MapNull(types.StringType),
		ConfigValue: configVal,
		PlanValue:   configVal,
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMapUnknownOnCreate().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.IsUnknown() {
		t.Errorf("expected unknown on create (server will enrich), got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMap_ConfigNull(t *testing.T) {
	ctx := context.Background()
	stateVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "1000",
		"server-key": "default",
	})
	req := planmodifier.MapRequest{
		StateValue:  stateVal,
		ConfigValue: types.MapNull(types.StringType),
		PlanValue:   types.MapNull(types.StringType),
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(stateVal) {
		t.Errorf("expected state preserved when config is null, got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMap_ConfigUnknown(t *testing.T) {
	ctx := context.Background()
	stateVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate": "1000",
	})
	req := planmodifier.MapRequest{
		StateValue:  stateVal,
		ConfigValue: types.MapUnknown(types.StringType),
		PlanValue:   types.MapUnknown(types.StringType),
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(stateVal) {
		t.Errorf("expected state preserved when config is unknown, got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMap_MergeNoChange(t *testing.T) {
	ctx := context.Background()
	stateVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "1000",
		"server-key": "default",
	})
	configVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate": "1000",
	})
	req := planmodifier.MapRequest{
		StateValue:  stateVal,
		ConfigValue: configVal,
		PlanValue:   configVal,
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(stateVal) {
		t.Errorf("expected merged map with server keys preserved, got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMap_MergeOverride(t *testing.T) {
	ctx := context.Background()
	stateVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "1000",
		"server-key": "default",
	})
	configVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate": "2000",
	})
	expected, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "2000",
		"server-key": "default",
	})
	req := planmodifier.MapRequest{
		StateValue:  stateVal,
		ConfigValue: configVal,
		PlanValue:   configVal,
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(expected) {
		t.Errorf("expected user override merged with server keys, got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMap_MergeAddKey(t *testing.T) {
	ctx := context.Background()
	stateVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "1000",
		"server-key": "default",
	})
	configVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":    "1000",
		"new-key": "val",
	})
	expected, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "1000",
		"server-key": "default",
		"new-key":    "val",
	})
	req := planmodifier.MapRequest{
		StateValue:  stateVal,
		ConfigValue: configVal,
		PlanValue:   configVal,
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(expected) {
		t.Errorf("expected new key added with server keys preserved, got %v", resp.PlanValue)
	}
}

func TestServerEnrichedMap_MergeRemoveKey(t *testing.T) {
	ctx := context.Background()
	stateVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{
		"rate":       "1000",
		"server-key": "default",
	})
	configVal, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{})
	req := planmodifier.MapRequest{
		StateValue:  stateVal,
		ConfigValue: configVal,
		PlanValue:   configVal,
	}
	resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
	planmodifiers.ServerEnrichedMap().PlanModifyMap(ctx, req, resp)

	if !resp.PlanValue.Equal(stateVal) {
		t.Errorf("expected state preserved when user removes keys from config, got %v", resp.PlanValue)
	}
}
