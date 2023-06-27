package ocicopyprovider

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/internal/ocicopyprovider/registrymodel"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewOciCopyProvider() provider.Provider {
	return &OciCopyProvider{}
}

type OciCopyProvider struct {
	// Attributes

	// Blocks
	Registries []registrymodel.RegistryConfigurationModel `tfsdk:"registry"`
}

func (p *OciCopyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	diags := req.Config.Get(ctx, p)
	res.Diagnostics.Append(diags...)
}

func (*OciCopyProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (*OciCopyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "ocicopy"
}

func (*OciCopyProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (*OciCopyProvider) Schema(_ context.Context, _ provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"registry": schema.SetNestedBlock{
				Description:  "Specify additional configuration for a specific registry",
				NestedObject: registrymodel.RegistryConfigurationSchema(),
			},
		},
		Description: "Global provider configuration that affects all resources and datasources",
	}
}
