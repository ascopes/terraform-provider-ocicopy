package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewOciCopyProvider() provider.Provider {
	return &OciCopyProvider{}
}

var _ provider.Provider = &OciCopyProvider{}

type OciCopyProvider struct {
	Registries RegistriesModel
}

func (provider *OciCopyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	provider.Registries = RegistriesModel{}
	req.Config.Get(ctx, provider.Registries)
}

func (*OciCopyProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (*OciCopyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "ocicopy"
}

func (provider *OciCopyProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return NewRepositoryResource(provider) },
	}
}

func (*OciCopyProvider) Schema(_ context.Context, _ provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Configure the 'ocicopy' Terraform provider.",
		Blocks: map[string]schema.Block{
			"registries": GetRegistriesBlockSchema(),
		},
	}
}
