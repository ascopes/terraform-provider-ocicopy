package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type imageResource struct {
	provider *providerModel
}

func NewImageResource(provider *providerModel) resource.Resource {
	return &imageResource{
		provider: provider,
	}
}

func (*imageResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_image"
}

func (*imageResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	panic("unimplemented")
}

func (*imageResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	panic("unimplemented")
}

func (*imageResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	panic("unimplemented")
}

func (*imageResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	panic("unimplemented")
}

func (*imageResource) Schema(context.Context, resource.SchemaRequest, *resource.SchemaResponse) {
	panic("unimplemented")
}
