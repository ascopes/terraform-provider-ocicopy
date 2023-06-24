package internal

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ RegistryModel = &BearerAuthRegistryModel{}

type BearerAuthRegistryModel struct {
	RegistryUrl types.String `tfsdk:"registry_url"`
	Token       types.String `tfsdk:"token"`
}

func (bearerAuth BearerAuthRegistryModel) GetAuthenticator() authn.Authenticator {
	return &authn.Bearer{Token: bearerAuth.Token.ValueString()}
}

func (bearerAuth BearerAuthRegistryModel) GetRegistryUrl() string {
	return bearerAuth.RegistryUrl.ValueString()
}

func GetBearerBlockSchema() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"registry_url": schema.StringAttribute{
				Description: "The registry URL to use.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"token": schema.StringAttribute{
				Description: "The registry bearer token to use.",
				Required:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}
