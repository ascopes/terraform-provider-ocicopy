package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type providerImpl struct {
	config *providerConfigModel
}

func NewProvider() provider.Provider {
	return &providerImpl{
		config: &providerConfigModel{},
	}
}

func (provider *providerImpl) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	diags := req.Config.Get(ctx, provider.config)
	res.Diagnostics.Append(diags...)
}

func (*providerImpl) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (*providerImpl) Metadata(ctx context.Context, req provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "ocicopy"
}

func (*providerImpl) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (*providerImpl) Schema(ctx context.Context, req provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = providerConfigModelSchema()
}
