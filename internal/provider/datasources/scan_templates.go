package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/petrutb/terraform-provider-runzero/internal/client"
)

var (
	_ datasource.DataSource              = &ScanTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &ScanTemplatesDataSource{}
)

type ScanTemplatesDataSource struct {
	client *client.Client
}

type ScanTemplatesDataSourceModel struct {
	ScanTemplates []ScanTemplateDataSourceModel `tfsdk:"scan_templates"`
}

func NewScanTemplatesDataSource() datasource.DataSource {
	return &ScanTemplatesDataSource{}
}

func (d *ScanTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scan_templates"
}

func (d *ScanTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all RunZero scan templates.",
		Attributes: map[string]schema.Attribute{
			"scan_templates": schema.ListNestedAttribute{
				Description: "List of scan templates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: scanTemplateDSSchemaAttributes(),
				},
			},
		},
	}
}

func (d *ScanTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *ScanTemplatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	templates, err := d.client.ListScanTemplates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing scan templates", err.Error())
		return
	}

	var state ScanTemplatesDataSourceModel
	for _, t := range templates {
		var m ScanTemplateDataSourceModel
		mapScanTemplateToDS(ctx, &t, &m, &resp.Diagnostics)
		state.ScanTemplates = append(state.ScanTemplates, m)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
