package internal

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

type RegistriesModel struct {
	BasicAuth  []BasicAuthRegistryModel  `tfsdk:"basic_auth"`
	BearerAuth []BearerAuthRegistryModel `tfsdk:"bearer_auth"`
	Ecr        []EcrRegistryModel        `tfsdk:"ecr"`
}

func GetRegistriesBlockSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Configure registries with authentication details.",
		Blocks: map[string]schema.Block{
			"basic_auth": schema.SetNestedBlock{
				Description:  "Configure registries to authenticate using basic authentication credentials.",
				NestedObject: GetBasicBlockSchema(),
			},
			"bearer_auth": schema.SetNestedBlock{
				Description:  "Configure registries to authenticate using bearer authentication credentials.",
				NestedObject: GetBearerBlockSchema(),
			},
			"ecr": schema.SetNestedBlock{
				Description:  "Configure registries that are hosted on AWS ECR.",
				NestedObject: GetEcrBlockSchema(),
			},
		},
	}
}
