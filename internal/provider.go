package internal

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/internal/config"
	"github.com/ascopes/terraform-provider-ocicopy/internal/image"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewOciCopyProvider() provider.Provider {
	return &OciCopyProvider{}
}

type OciCopyProvider struct {
	config *config.ConfigModel
}

func (p *OciCopyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	config := &config.ConfigModel{}
	diags := req.Config.Get(ctx, config)
	res.Diagnostics.Append(diags...)

	p.config = config
}

func (*OciCopyProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (*OciCopyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "ocicopy"
}

func (provider *OciCopyProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return image.NewImageCopyResource(provider.config) },
	}
}

func (*OciCopyProvider) Schema(_ context.Context, _ provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = config.ConfigSchema()
}
