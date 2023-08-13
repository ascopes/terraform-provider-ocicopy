package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type providerImpl struct {
	model *providerModel
}

func NewProvider() provider.Provider {
	return &providerImpl{
		model: &providerModel{},
	}
}

func (provider *providerImpl) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	diags := req.Config.Get(ctx, provider.model)
	res.Diagnostics.Append(diags...)
}

func (provider *providerImpl) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (provider *providerImpl) Metadata(ctx context.Context, req provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "ocicopy"
}

func (provider *providerImpl) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//func() resource.Resource { return NewImageResource(provider.model) },
	}
}

func (provider *providerImpl) Schema(ctx context.Context, req provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = provider.model.schema()
}
