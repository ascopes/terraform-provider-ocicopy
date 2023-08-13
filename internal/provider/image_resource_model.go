package provider

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

type imageResourceModel struct {
}

func (model *imageResourceModel) schema() schema.Schema {
	return schema.Schema{}
}
