package config

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BearerAuthModel struct {
	Token types.String `tfsdk:"token"`
}

func bearerAuthModelSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "The bearer authentication token",
				Required:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Description: "Configure the provider to use bearer credentials to authenticate with the registry",
		Optional:    true,
		Validators: []validator.Object{
			objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic_authentication")),
		},
	}
}
