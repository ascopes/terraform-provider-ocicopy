package config

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

// Model for global configuration that is passed between providers, datasources, and resources.
type ConfigModel struct {
	Registries []RegistryConfigurationModel `tfsdk:"registry"`
}

func ConfigSchema() schema.Schema {
	return schema.Schema{
		Blocks: map[string]schema.Block{
			"registry": schema.SetNestedBlock{
				Description:  "Specify additional configuration for a specific registry",
				NestedObject: registryConfigurationSchema(),
			},
		},
		Description: "Global provider configuration that affects all resources and datasources",
	}
}
