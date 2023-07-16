package image

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type imageCopyModel struct {
	Source sourceImageModel `tfsdk:"source"`
	Target sourceImageModel `tfsdk:"target"`
}

func imageCopyModelSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"source": sourceImageModelSchema(),
			"target": targetImageModelSchema(),
		},
		Description: "A description of how to copy a specific image between registries",
	}
}

type sourceImageModel struct {
	RegistryUrl types.String `tfsdk:"registry_url"`
	Image       types.String `tfsdk:"image"`
	Tag         types.String `tfsdk:"tag"`
}

func sourceImageModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The registry URL to copy the image from",
				Required:    true,
			},
			"image": schema.StringAttribute{
				Description: "The image name to copy",
				Required:    true,
			},
			"tag": schema.StringAttribute{
				Description: "The image tag to copy",
				Required:    true,
			},
		},
		Description: "A description of the source of a given image",
		Required:    true,
	}
}

type targetImageModel struct {
	RegistryUrl types.String `tfsdk:"registry_url"`
	Image       types.String `tfsdk:"image"`
	Tag         types.String `tfsdk:"tag"`
}

func targetImageModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The registry URL to copy the image to",
				Required:    true,
			},
			"image": schema.StringAttribute{
				Description: "The image name to store (or omit to use the source image name)",
				Optional:    true,
			},
			"tag": schema.StringAttribute{
				Description: "The image tag to store (or omit to use the source tag)",
				Optional:    true,
			},
		},
		Description: "A description of the destination of the copy of the given image",
		Required:    true,
	}
}
