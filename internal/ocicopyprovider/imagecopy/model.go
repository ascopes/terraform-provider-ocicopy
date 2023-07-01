package imagecopy

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
}

func sourceImageModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The registry URL to copy the image from",
				Required:    true,
			},
		},
		Description: "A description of the source of a given image",
		Required:    true,
	}
}

type targetImageModel struct {
	RegistryUrl types.String `tfsdk:"registry_url"`
}

func targetImageModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The registry URL to copy the image to",
				Required:    true,
			},
		},
		Description: "A description of the destination of the copy of the given image",
		Required:    true,
	}
}
