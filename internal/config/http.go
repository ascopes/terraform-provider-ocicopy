package config

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Configuration for HTTP transports.
type RegistryHttpConfigurationModel struct {
	ForceHttp2          *types.Bool   `tfsdk:"force_http2"`
	Jobs                *types.Int64  `tfsdk:"jobs"`
	KeepAlive           *DurationType `tfsdk:"keep_alive"`
	Timeout             *DurationType `tfsdk:"timeout"`
	TlsHandshakeTimeout *DurationType `tfsdk:"tls_handshake_timeout"`
}

// Schema for the HTTP configuration attribute.
func registryHttpConfigurationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"force_http2": schema.BoolAttribute{
				Description: "Set whether to explicitly force attempts to use HTTP/2",
				Optional:    true,
			},
			"jobs": schema.Int64Attribute{
				Description: "Set the number of concurrent jobs to execute each request with",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"keep_alive": schema.StringAttribute{
				Description: "Set the connection keepalive duration",
				CustomType:  DurationType{},
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"timeout": schema.StringAttribute{
				Description: "Set the request keepalive duration",
				CustomType:  DurationType{},
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"tls_handshake_timeout": schema.StringAttribute{
				Description: "Set the TLS handshake timeout duration",
				CustomType:  DurationType{},
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Description: "Override the underlying HTTP settings",
		Optional:    true,
	}
}
