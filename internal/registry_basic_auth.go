package internal

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type basicAuthRegistryModel struct {
	RegistryUrl types.String `tfsdk:"registry_url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
}

func (basicAuth basicAuthRegistryModel) GetAuthenticator() authn.Authenticator {
	return &authn.Basic{
		Username: basicAuth.Username.ValueString(),
		Password: basicAuth.Password.ValueString(),
	}
}

func (basicAuth basicAuthRegistryModel) GetRegistryUrl() string {
	return basicAuth.RegistryUrl.ValueString()
}

func getBasicBlockSchema() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "The registry username to use.",
				Required:    true,
				Sensitive:   true,
			},
			"password": schema.StringAttribute{
				Description: "The registry password to use.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}
