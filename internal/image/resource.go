package image

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/internal/api"
	"github.com/ascopes/terraform-provider-ocicopy/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewImageCopyResource(config *config.ConfigModel) resource.Resource {
	return &imageCopyResource{config: config}
}

type imageCopyResource struct {
	config *config.ConfigModel
}

func (*imageCopyResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	model := &imageCopyModel{}

	diags := req.Plan.Get(ctx, model)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, "Creating new tag in repository", map[string]any{
		"id":                "tbc",
		"source_registry":   model.Source.Registry,
		"source_repository": model.Source.Repository,
		"source_digest":     model.Source.Digest,
		"target_registry":   model.Target.Registry,
		"target_repository": model.Target.Repository,
		"target_tag":        model.Target.Tag,
		"target_digest":     "tbc",
	})

	// TODO: handle registry configuration here.
	srcImage := api.NewImage(
		api.NewRegistryClientConfig(model.Source.Registry.ValueString()),
		model.Source.Registry.ValueString(),
		model.Source.Repository.ValueString(),
	)

	destImage := api.NewImage(
		api.NewRegistryClientConfig(model.Target.Registry.ValueString()),
		model.Target.Registry.ValueString(),
		model.Target.Repository.ValueString(),
	)

	srcImageObj, err := srcImage.GetImageDescriptorFromDigest(
		ctx,
		model.Source.Digest.ValueString(),
	)

	if err != nil {
		res.Diagnostics.AddError("Failed to fetch source image descriptor", err.Error())
		return
	}

	digest, err := destImage.UploadImageAsTag(
		ctx,
		model.Target.Tag.ValueString(),
		srcImageObj,
	)

	if err != nil {
		res.Diagnostics.AddError("Failed to upload source image to target registry", err.Error())
		return
	}

	// Update state
	model.Target.Digest = types.StringValue(digest)
	model.computeId()

	tflog.Info(ctx, "Created new tag in repository", map[string]any{
		"id":                model.Id,
		"source_registry":   model.Source.Registry,
		"source_repository": model.Source.Repository,
		"source_digest":     model.Source.Digest,
		"target_registry":   model.Target.Registry,
		"target_repository": model.Target.Repository,
		"target_tag":        model.Target.Tag,
		"target_digest":     model.Target.Digest,
	})

	// Store state
	diags = res.State.Set(ctx, model)
	res.Diagnostics.Append(diags...)
}

func (*imageCopyResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	model := &imageCopyModel{}

	diags := req.State.Get(ctx, model)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting existing digest in repository", map[string]any{
		"id":                model.Id,
		"target_registry":   model.Target.Registry,
		"target_repository": model.Target.Repository,
		"target_tag":        model.Target.Tag,
		"target_digest":     model.Target.Digest,
	})

	destImage := api.NewImage(
		api.NewRegistryClientConfig(model.Target.Registry.ValueString()),
		model.Target.Registry.ValueString(),
		model.Target.Repository.ValueString(),
	)

	err := destImage.DeleteImage(ctx, model.Target.Digest.ValueString())
	if err != nil {
		res.Diagnostics.AddError("Failed to delete tag from target repository", err.Error())
	}
}

func (*imageCopyResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_image"
}

func (*imageCopyResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	model := &imageCopyModel{}

	diags := req.State.Get(ctx, model)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	model.computeId()
	diags = res.State.Set(ctx, model)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Computed state for registry", map[string]any{
		"id":                model.Id,
		"source_registry":   model.Source.Registry,
		"source_repository": model.Source.Repository,
		"source_digest":     model.Source.Digest,
		"target_registry":   model.Target.Registry,
		"target_repository": model.Target.Repository,
		"target_tag":        model.Target.Tag,
		"target_digest":     model.Target.Digest,
	})
}

func (*imageCopyResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = imageCopyModelSchema()
}

func (*imageCopyResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	panic("unimplemented")
}
