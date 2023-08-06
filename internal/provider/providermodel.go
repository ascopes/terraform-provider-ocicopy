package provider

import (
	"github.com/ascopes/terraform-provider-ocicopy/internal/durationtype"
	"github.com/ascopes/terraform-provider-ocicopy/internal/mapping"
	"github.com/ascopes/terraform-provider-ocicopy/internal/registryapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type providerModel struct {
	// Used as a repeatable block.
	Registries []registryConfigModel `tfsdk:"registry"`
}

func (configModel *providerModel) getRegistryConfig(registryUrl string) registryapi.RegistryConfig {
	for _, registryConfigModel := range configModel.Registries {
		// TODO: could this become "unknown" prior to apply? How do we deal with that..?
		if registryConfigModel.Url.ValueString() == registryUrl {
			return registryConfigModel.toRegistryConfig()
		}
	}

	// Return defaults for anything else.
	return registryapi.NewRegistryConfig()
}

type registryConfigModel struct {
	BasicAuth             *basicAuthModel            `tfsdk:"basic_auth"`
	BearerAuth            *bearerAuthModel           `tfsdk:"bearer_auth"`
	ConcurrentJobs        types.Int64                `tfsdk:"concurrent_jobs"`
	ConnectTimeout        durationtype.DurationValue `tfsdk:"connect_timeout"`
	ForceAttemptHttp2     types.Bool                 `tfsdk:"force_attempt_http2"`
	IdleConnectionTimeout durationtype.DurationValue `tfsdk:"idle_connection_timeout"`
	Insecure              types.Bool                 `tfsdk:"insecure"`
	KeepAlive             durationtype.DurationValue `tfsdk:"keep_alive"`
	MaxIdleConnections    types.Int64                `tfsdk:"max_idle_connections"`
	ResponseTimeout       durationtype.DurationValue `tfsdk:"response_timeout"`
	TlsHandshakeTimeout   durationtype.DurationValue `tfsdk:"tls_handshake_timeout"`
	Url                   types.String               `tfsdk:"url"`
}

// Create a registryapi-compatible registry configuration object from the given
// registry config Terraform model.
func (configModel *registryConfigModel) toRegistryConfig() registryapi.RegistryConfig {
	config := registryapi.NewRegistryConfig()

	mapping.IfNotNil(
		configModel.BasicAuth,
		func(value *basicAuthModel) { config.Authenticator = value.toAuthenticator() },
	)
	mapping.IfNotNil(
		configModel.BearerAuth,
		func(value *bearerAuthModel) { config.Authenticator = value.toAuthenticator() },
	)
	mapping.IfPresent(
		configModel.ConcurrentJobs,
		func(value types.Int64) { config.ConcurrentJobs = int(value.ValueInt64()) },
	)
	mapping.IfPresent(
		configModel.ConnectTimeout,
		func(value durationtype.DurationValue) { config.ConnectTimeout = value.ValueDuration() },
	)
	mapping.IfPresent(
		configModel.ForceAttemptHttp2,
		func(value types.Bool) { config.ForceAttemptHttp2 = value.ValueBool() },
	)
	mapping.IfPresent(
		configModel.IdleConnectionTimeout,
		func(value durationtype.DurationValue) { config.IdleConnectionTimeout = value.ValueDuration() },
	)
	mapping.IfPresent(
		configModel.Insecure,
		func(value types.Bool) { config.Insecure = value.ValueBool() },
	)
	mapping.IfPresent(
		configModel.KeepAlive,
		func(value durationtype.DurationValue) { config.KeepAlive = value.ValueDuration() },
	)
	mapping.IfPresent(
		configModel.MaxIdleConnections,
		func(value types.Int64) { config.MaxIdleConnections = int(value.ValueInt64()) },
	)
	mapping.IfPresent(
		configModel.ResponseTimeout,
		func(value durationtype.DurationValue) { config.ResponseTimeout = value.ValueDuration() },
	)
	mapping.IfPresent(
		configModel.TlsHandshakeTimeout,
		func(value durationtype.DurationValue) { config.TlsHandshakeTimeout = value.ValueDuration() },
	)

	return config
}

type basicAuthModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (authModel *basicAuthModel) toAuthenticator() registryapi.Authenticator {
	return registryapi.NewBasicAuthenticator(
		authModel.Username.ValueString(),
		authModel.Password.ValueString(),
	)
}

type bearerAuthModel struct {
	Token types.String `tfsdk:"token"`
}

func (authModel *bearerAuthModel) toAuthenticator() registryapi.Authenticator {
	return registryapi.NewBearerAuthenticator(
		authModel.Token.ValueString(),
	)
}

func providerConfigModelSchema() schema.Schema {
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
							CustomType:  durationtype.DurationType{},
							Description: "The maximum duration to wait for a HTTP connection to complete before timing out",
							Optional:    true,
						},
						"force_attempt_http2": schema.BoolAttribute{
							Description: "If true, then force attempting to communicate over HTTP/2",
							Optional:    true,
						},
						"idle_connection_timeout": schema.StringAttribute{
							CustomType:  durationtype.DurationType{},
							Description: "The maximum duration to allow a HTTP connection to be idle before renewing it",
							Optional:    true,
						},
						"insecure": schema.BoolAttribute{
							Description: "Set to true to disable SSL and fall back to plain HTTP rather than HTTPS",
							Optional:    true,
						},
						"keep_alive": schema.StringAttribute{
							CustomType:  durationtype.DurationType{},
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
							CustomType:  durationtype.DurationType{},
							Description: "The maximum duration to wait for an HTTP response to be received before timing out",
							Optional:    true,
						},
						"tls_handshake_timeout": schema.StringAttribute{
							CustomType:  durationtype.DurationType{},
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
