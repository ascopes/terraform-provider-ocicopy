package registrymodel

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Configuration for a specific registry.
type RegistryConfigurationModel struct {
	BasicAuth  *BasicAuthModel                 `tfsdk:"basic_authentication"`
	BearerAuth *BearerAuthModel                `tfsdk:"bearer_authentication"`
	Http       *RegistryHttpConfigurationModel `tfsdk:"http"`
	Url        string                          `tfsdk:"url"`
}

// Schema for specific registry configuration details.
func RegistryConfigurationSchema() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"basic_authentication":  basicAuthModelSchema(),
			"bearer_authentication": bearerAuthModelSchema(),
			"http":                  registryHttpConfigurationSchema(),
			"url": schema.StringAttribute{
				Description: "Set the registry URL",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Validators: []validator.Object{
			objectvalidator.IsRequired(),
		},
	}
}
