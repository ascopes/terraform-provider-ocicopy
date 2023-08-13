package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type imageResource struct {
	provider *providerModel
	model    *imageResourceModel
}

func NewImageResource(provider *providerModel) resource.Resource {
	return &imageResource{
		provider: provider,
		model:    &imageResourceModel{},
	}
}

func (resource *imageResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_image"
}

func (resource *imageResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	panic("unimplemented")
}

func (resource *imageResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	panic("unimplemented")
}

func (resource *imageResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	panic("unimplemented")
}

func (resource *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	panic("unimplemented")
}

func (resource *imageResource) Schema(context.Context, resource.SchemaRequest, *resource.SchemaResponse) {
	return
}
