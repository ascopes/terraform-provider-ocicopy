package internal

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

type registriesModel struct {
	BasicAuth  []basicAuthRegistryModel  `tfsdk:"basic_auth"`
	BearerAuth []bearerAuthRegistryModel `tfsdk:"bearer_auth"`
	Ecr        []ecrRegistryModel        `tfsdk:"ecr"`
}

type registriesMapping map[string]registryModel

func getRegistriesBlockSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Configure registries with authentication details",
		Blocks: map[string]schema.Block{
			"basic_auth": schema.SetNestedBlock{
				Description:  "Configure registries to authenticate using basic authentication credentials",
				NestedObject: getBasicBlockSchema(),
			},
			"bearer_auth": schema.SetNestedBlock{
				Description:  "Configure registries to authenticate using bearer authentication credentials",
				NestedObject: GetBearerBlockSchema(),
			},
			"ecr": schema.SetNestedBlock{
				Description:  "Configure registries that are hosted on Amazon ECR",
				NestedObject: GetEcrBlockSchema(),
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
	for _, registryModel := range registries.Ecr {
		mapping[registryModel.GetRegistryUrl()] = registryModel
	}
	return mapping
}
