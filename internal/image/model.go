package image

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type imageCopyModel struct {
	Id     types.String     `tfsdk:"id"`
	Source sourceImageModel `tfsdk:"source"`
	Target targetImageModel `tfsdk:"target"`
}

func (image *imageCopyModel) computeId() {
	sourceRepository := image.Source.Repository.ValueString()
	targetRepository := image.Target.Repository.ValueString()
	if len(targetRepository) == 0 {
		// We use the source repository name if we don't override it.
		targetRepository = sourceRepository
	}

	idDesc := idDescriptor{
		SourceRegistry:   image.Source.Registry.ValueString(),
		SourceRepository: sourceRepository,
		SourceDigest:     image.Source.Repository.ValueString(),
		TargetRegistry:   image.Target.Registry.ValueString(),
		TargetRepository: targetRepository,
		TargetTag:        image.Target.Tag.ValueString(),
	}

	image.Id = types.StringValue(idDesc.encode())
}

func imageCopyModelSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An internal ID used by Terraform to uniquely identify this resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": sourceImageModelSchema(),
			"target": targetImageModelSchema(),
		},
		Description: "A description of how to copy a specific image between registries",
	}
}

type sourceImageModel struct {
	Registry   types.String `tfsdk:"registry"`
	Repository types.String `tfsdk:"repository"`
	Digest     types.String `tfsdk:"digest"`
}

func sourceImageModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"registry": schema.StringAttribute{
				Description: "The registry URL to copy the image from",
				Required:    true,
			},
			"repository": schema.StringAttribute{
				Description: "The repository name to copy the image from",
				Required:    true,
			},
			"digest": schema.StringAttribute{
				Description: "The digest to copy",
				Required:    true,
			},
		},
		Description: "A description of the source of a given image",
		Required:    true,
	}
}

type targetImageModel struct {
	Registry   types.String `tfsdk:"registry"`
	Repository types.String `tfsdk:"repository"`
	Tag        types.String `tfsdk:"tag"`
	Digest     types.String `tfsdk:"digest"`
}

func targetImageModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"registry": schema.StringAttribute{
				Description: "The registry URL to copy the image to",
				Required:    true,
			},
			"repository": schema.StringAttribute{
				Description: "The repository to copy the image to",
				Optional:    true,
			},
			"tag": schema.StringAttribute{
				Description: "The image tag to store",
				Required:    true,
			},
			"digest": schema.StringAttribute{
				Description: "The existing digest on the remote repository",
				Computed:    true,
			},
		},
		Description: "A description of the destination of the copy of the given image",
		Required:    true,
	}
}
