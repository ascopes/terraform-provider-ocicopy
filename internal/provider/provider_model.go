package provider

import (
	"github.com/ascopes/terraform-provider-ocicopy/internal/duration_type"
	"github.com/ascopes/terraform-provider-ocicopy/internal/registry_api"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type providerModel struct {
	// Used as a repeatable block.
	Registries []registryConfigModel `tfsdk:"registry"`
}

func (model *providerModel) schema() schema.Schema {
	return schema.Schema{
		Description: "Global configuration for image copy operations",
		Blocks: map[string]schema.Block{
			"registry": schema.SetNestedBlock{
				Description: "Override the default configuration for a specific registry",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"basic_auth": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"username": schema.StringAttribute{
									Description: "The username to use",
									Required:    true,
									Sensitive:   true,
								},
								"password": schema.StringAttribute{
									Description: "The password to use",
									Required:    true,
									Sensitive:   true,
								},
							},
							Description: "Configure basic authentication credentials to use with this registry",
							Optional:    true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer_auth")),
							},
						},
						"bearer_auth": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"token": schema.StringAttribute{
									Description: "The bearer authentication token to use",
									Required:    true,
									Sensitive:   true,
								},
							},
							Description: "Configure bearer authentication credentials to use with this registry",
							Optional:    true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic_auth")),
							},
						},
						"concurrent_jobs": schema.Int64Attribute{
							Description: "Number of concurrent HTTP jobs to allow to run for this registry",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						"connect_timeout": schema.StringAttribute{
							CustomType:  duration_type.DurationType{},
							Description: "The maximum duration to wait for a HTTP connection to complete before timing out",
							Optional:    true,
						},
						"force_attempt_http2": schema.BoolAttribute{
							Description: "If true, then force attempting to communicate over HTTP/2",
							Optional:    true,
						},
						"idle_connection_timeout": schema.StringAttribute{
							CustomType:  duration_type.DurationType{},
							Description: "The maximum duration to allow a HTTP connection to be idle before renewing it",
							Optional:    true,
						},
						"insecure": schema.BoolAttribute{
							Description: "Set to true to disable SSL and fall back to plain HTTP rather than HTTPS",
							Optional:    true,
						},
						"keep_alive": schema.StringAttribute{
							CustomType:  duration_type.DurationType{},
							Description: "The HTTP connection keep-alive timeout",
							Optional:    true,
						},
						"max_idle_connections": schema.Int64Attribute{
							Description: "Number of connections that can remain open when idle",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						"response_timeout": schema.StringAttribute{
							CustomType:  duration_type.DurationType{},
							Description: "The maximum duration to wait for an HTTP response to be received before timing out",
							Optional:    true,
						},
						"tls_handshake_timeout": schema.StringAttribute{
							CustomType:  duration_type.DurationType{},
							Description: "The maximum duration to wait for a TLS handshake to complete before timing out",
							Optional:    true,
						},
						"url": schema.StringAttribute{
							Description: "The URL of the registry to configure, minus any protocol",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (model *providerModel) getRegistryConfig(registryUrl string) registry_api.RegistryConfig {
	for _, registryConfigModel := range model.Registries {
		// TODO: could this become "unknown" prior to apply? How do we deal with that..?
		if registryConfigModel.Url.ValueString() == registryUrl {
			return registryConfigModel.toRegistryConfig()
		}
	}

	// Return defaults for anything else.
	return registry_api.NewRegistryConfig()
}

type registryConfigModel struct {
	BasicAuth             *basicAuthModel             `tfsdk:"basic_auth"`
	BearerAuth            *bearerAuthModel            `tfsdk:"bearer_auth"`
	ConcurrentJobs        types.Int64                 `tfsdk:"concurrent_jobs"`
	ConnectTimeout        duration_type.DurationValue `tfsdk:"connect_timeout"`
	ForceAttemptHttp2     types.Bool                  `tfsdk:"force_attempt_http2"`
	IdleConnectionTimeout duration_type.DurationValue `tfsdk:"idle_connection_timeout"`
	Insecure              types.Bool                  `tfsdk:"insecure"`
	KeepAlive             duration_type.DurationValue `tfsdk:"keep_alive"`
	MaxIdleConnections    types.Int64                 `tfsdk:"max_idle_connections"`
	ResponseTimeout       duration_type.DurationValue `tfsdk:"response_timeout"`
	TlsHandshakeTimeout   duration_type.DurationValue `tfsdk:"tls_handshake_timeout"`
	Url                   types.String                `tfsdk:"url"`
}

// Create a registry_api-compatible registry configuration object from the given
// registry config Terraform model.
func (model *registryConfigModel) toRegistryConfig() registry_api.RegistryConfig {
	config := registry_api.NewRegistryConfig()

	if value := model.BasicAuth; value != nil {
		config.Authenticator = value.toAuthenticator()
	}

	if value := model.BearerAuth; value != nil {
		config.Authenticator = value.toAuthenticator()
	}

	if value := model.ConcurrentJobs; isDefined(value) {
		config.ConcurrentJobs = int(value.ValueInt64())
	}

	if value := model.ConnectTimeout; isDefined(value) {
		config.ConnectTimeout = value.ValueDuration()
	}

	if value := model.ForceAttemptHttp2; isDefined(value) {
		config.ForceAttemptHttp2 = value.ValueBool()
	}

	if value := model.IdleConnectionTimeout; isDefined(value) {
		config.IdleConnectionTimeout = value.ValueDuration()
	}

	if value := model.Insecure; isDefined(value) {
		config.Insecure = value.ValueBool()
	}

	if value := model.KeepAlive; isDefined(value) {
		config.KeepAlive = value.ValueDuration()
	}

	if value := model.MaxIdleConnections; isDefined(value) {
		config.MaxIdleConnections = int(value.ValueInt64())
	}

	if value := model.ResponseTimeout; isDefined(value) {
		config.ResponseTimeout = value.ValueDuration()
	}

	if value := model.TlsHandshakeTimeout; isDefined(value) {
		config.TlsHandshakeTimeout = value.ValueDuration()
	}

	return config
}

type basicAuthModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (authModel *basicAuthModel) toAuthenticator() registry_api.Authenticator {
	return registry_api.NewBasicAuthenticator(
		authModel.Username.ValueString(),
		authModel.Password.ValueString(),
	)
}

type bearerAuthModel struct {
	Token types.String `tfsdk:"token"`
}

func (authModel *bearerAuthModel) toAuthenticator() registry_api.Authenticator {
	return registry_api.NewBearerAuthenticator(
		authModel.Token.ValueString(),
	)
}

func isDefined(value attr.Value) bool {
	return !value.IsUnknown() && !value.IsNull()
}
