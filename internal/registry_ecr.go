package internal

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ RegistryModel = &BasicAuthRegistryModel{}

type EcrRegistryModel struct {
	RegistryUrl types.String `tfsdk:"registry_url"`
	Token       types.String `tfsdk:"token"`
}

func (ecr EcrRegistryModel) GetAuthenticator() authn.Authenticator {
	return &authn.Basic{
		Username: "AWS",
		Password: ecr.Token.String(),
	}
}

func (ecr EcrRegistryModel) GetRegistryUrl() string {
	return ecr.RegistryUrl.ValueString()
}

func GetEcrBlockSchema() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The AWS ECR URL to use.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"token": schema.StringAttribute{
				Description: "The AWS ECR token to use.",
				Required:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}
