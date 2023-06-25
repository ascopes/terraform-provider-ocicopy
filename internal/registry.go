package internal

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type registryModel struct {
	Url        types.String             `tfsdk:"url"`
	BasicAuth  *registryBasicAuthModel  `tfsdk:"basic_auth"`
	BearerAuth *registryBearerAuthModel `tfsdk:"bearer_auth"`
}

func (registry registryModel) getAuthenticator() authn.Authenticator {
	if registry.BasicAuth != nil {
		return &authn.Basic{
			Username: registry.BasicAuth.Username.ValueString(),
			Password: registry.BasicAuth.Password.ValueString(),
		}
	}

	if registry.BearerAuth != nil {
		return &authn.Bearer{
			Token: registry.BearerAuth.Token.ValueString(),
		}
	}

	return authn.Anonymous
}

type registryBasicAuthModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type registryBearerAuthModel struct {
	Token types.String `tfsdk:"token"`
}

func getRegistryBlockSchema() schema.SetNestedBlock {
	return schema.SetNestedBlock{
		Description: "Configure how the plugin connects to a registry",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					Description: "URL that identifies the registry to provide additional configuration for",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"basic_auth": schema.SingleNestedAttribute{
					Description: "Configuration for Basic authentication",
					Optional:    true,
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
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer_auth")),
					},
				},
				"bearer_auth": schema.SingleNestedAttribute{
					Description: "Configuration for Bearer authentication",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"token": schema.StringAttribute{
							Description: "The registry bearer token to use.",
							Required:    true,
							Sensitive:   true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic_auth")),
					},
				},
			},
		},
	}
}
