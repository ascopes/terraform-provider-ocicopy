package internal

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type registriesModel struct {
	BasicAuth  []basicAuthRegistryModel  `tfsdk:"basic_auth"`
	BearerAuth []bearerAuthRegistryModel `tfsdk:"bearer_auth"`
}

type registriesMapping map[string]registryModel

func getRegistriesBlockSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Configure registries with authentication details",
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The registry URL to use.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"basic_auth": schema.SetNestedBlock{
				Description:  "Configure the registry to use Basic authentication",
				NestedObject: getBasicBlockSchema(),
			},
			"bearer_auth": schema.SetNestedBlock{
				Description:  "Configure the registry to use Bearer authentication",
				NestedObject: GetBearerBlockSchema(),
			},
		},
	}
}

func (registries registriesModel) getRegistriesMapping() registriesMapping {
	mapping := make(registriesMapping)
	for _, registryModel := range registries.BasicAuth {
		mapping[registryModel.GetRegistryUrl()] = registryModel
	}
	for _, registryModel := range registries.BearerAuth {
		mapping[registryModel.GetRegistryUrl()] = registryModel
	}
	return mapping
}
