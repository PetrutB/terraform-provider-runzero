package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	expectedResourceCount   = 17
	expectedDataSourceCount = 42
)

// Resources that manage objects created outside Terraform and are adopted by ID.
var adoptedResources = map[string]bool{
	"runzero_agent":    true,
	"runzero_explorer": true,
	"runzero_asset":    true,
}

func TestProvider_Metadata(t *testing.T) {
	p := New("1.2.3")()
	resp := &fwprovider.MetadataResponse{}
	p.Metadata(context.Background(), fwprovider.MetadataRequest{}, resp)
	if resp.TypeName != "runzero" {
		t.Errorf("type name: want runzero, got %s", resp.TypeName)
	}
	if resp.Version != "1.2.3" {
		t.Errorf("version: want 1.2.3, got %s", resp.Version)
	}
}

func TestProvider_Schema(t *testing.T) {
	ctx := context.Background()
	resp := &fwprovider.SchemaResponse{}
	New("test")().Schema(ctx, fwprovider.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("invalid provider schema: %v", diags)
	}
	for _, name := range []string{"org_token", "account_token", "client_id", "client_secret"} {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("missing attribute %s", name)
		}
		if !attr.IsSensitive() {
			t.Errorf("%s must be sensitive", name)
		}
	}
}

func TestProvider_ResourceSchemas(t *testing.T) {
	ctx := context.Background()
	factories := New("test")().Resources(ctx)
	if len(factories) != expectedResourceCount {
		t.Errorf("want %d resources, got %d", expectedResourceCount, len(factories))
	}

	seen := map[string]bool{}
	for _, factory := range factories {
		r := factory()
		meta := &resource.MetadataResponse{}
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "runzero"}, meta)
		name := meta.TypeName

		t.Run(name, func(t *testing.T) {
			if !strings.HasPrefix(name, "runzero_") {
				t.Errorf("type name %q must start with runzero_", name)
			}
			if seen[name] {
				t.Errorf("duplicate resource type %s", name)
			}
			seen[name] = true

			resp := &resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
			}
			if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
				t.Fatalf("invalid schema: %v", diags)
			}
			if resp.Schema.Description == "" {
				t.Error("schema description is empty")
			}
			id, ok := resp.Schema.Attributes["id"]
			if !ok {
				t.Fatal("missing id attribute")
			}
			if adoptedResources[name] {
				if !id.IsRequired() {
					t.Error("id of an adopted resource must be required")
				}
			} else if !id.IsComputed() {
				t.Error("id must be computed")
			}
			for attrName, attr := range resp.Schema.Attributes {
				if attr.GetDescription() == "" && attr.GetMarkdownDescription() == "" {
					t.Errorf("attribute %s has no description", attrName)
				}
			}
			if _, ok := r.(resource.ResourceWithConfigure); !ok {
				t.Error("resource does not implement Configure")
			}
		})
	}
}

func TestProvider_ResourceImportSupport(t *testing.T) {
	ctx := context.Background()
	notImportable := map[string]bool{
		"runzero_scan":        true,
		"runzero_org_api_key": true,
	}
	for _, factory := range New("test")().Resources(ctx) {
		r := factory()
		meta := &resource.MetadataResponse{}
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "runzero"}, meta)
		_, importable := r.(resource.ResourceWithImportState)
		if importable == notImportable[meta.TypeName] {
			t.Errorf("%s: importable=%v, expected %v", meta.TypeName, importable, !notImportable[meta.TypeName])
		}
	}
}

func TestProvider_DataSourceSchemas(t *testing.T) {
	ctx := context.Background()
	factories := New("test")().DataSources(ctx)
	if len(factories) != expectedDataSourceCount {
		t.Errorf("want %d data sources, got %d", expectedDataSourceCount, len(factories))
	}

	seen := map[string]bool{}
	for _, factory := range factories {
		d := factory()
		meta := &datasource.MetadataResponse{}
		d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "runzero"}, meta)
		name := meta.TypeName

		t.Run(name, func(t *testing.T) {
			if !strings.HasPrefix(name, "runzero_") {
				t.Errorf("type name %q must start with runzero_", name)
			}
			if seen[name] {
				t.Errorf("duplicate data source type %s", name)
			}
			seen[name] = true

			resp := &datasource.SchemaResponse{}
			d.Schema(ctx, datasource.SchemaRequest{}, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
			}
			if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
				t.Fatalf("invalid schema: %v", diags)
			}
			if resp.Schema.Description == "" {
				t.Error("schema description is empty")
			}
			if _, ok := d.(datasource.DataSourceWithConfigure); !ok {
				t.Error("data source does not implement Configure")
			}
		})
	}
}

func TestResolveString(t *testing.T) {
	const env = "RUNZERO_TEST_RESOLVE"

	t.Setenv(env, "from-env")
	if got := resolveString(types.StringValue("from-config"), env, "default"); got != "from-config" {
		t.Errorf("config should win, got %s", got)
	}
	if got := resolveString(types.StringNull(), env, "default"); got != "from-env" {
		t.Errorf("env should be used when config is null, got %s", got)
	}
	if got := resolveString(types.StringUnknown(), env, "default"); got != "from-env" {
		t.Errorf("env should be used when config is unknown, got %s", got)
	}

	t.Setenv(env, "")
	if got := resolveString(types.StringNull(), env, "default"); got != "default" {
		t.Errorf("default should be used, got %s", got)
	}
}
