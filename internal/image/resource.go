package image

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewImageCopyResource(config *config.ConfigModel) resource.Resource {
	return &imageCopyResource{config: config}
}

type imageCopyResource struct {
	config *config.ConfigModel
}

func (*imageCopyResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	panic("unimplemented")
}

func (*imageCopyResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	panic("unimplemented")
}

func (*imageCopyResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_image"
}

func (*imageCopyResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	panic("unimplemented")
}

func (*imageCopyResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = imageCopyModelSchema()
}

func (*imageCopyResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	panic("unimplemented")
}
