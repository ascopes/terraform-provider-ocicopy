package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewOciCopyProvider() provider.Provider {
	return &ociCopyProvider{}
}

type ociCopyProvider struct {
	Registries []registryModel
}

func (provider *ociCopyProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// Nothing to do here
}

func (*ociCopyProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (*ociCopyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "ocicopy"
}

func (provider *ociCopyProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return newRepositoryResource(provider) },
	}
}

func (*ociCopyProvider) Schema(_ context.Context, _ provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Configuration and settings for this provider",
		Blocks: map[string]schema.Block{
			"registry": getRegistryBlockSchema(),
		},
	}
}

func (provider ociCopyProvider) getRegistriesMapping() map[string]registryModel {
	mapping := make(map[string]registryModel, len(provider.Registries))

	for _, registry := range provider.Registries {
		mapping[registry.Url.ValueString()] = registry
	}

	return mapping
}
