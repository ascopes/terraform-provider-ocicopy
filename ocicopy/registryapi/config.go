package registryapi

import (
	"net"
	"net/http"
	"time"

	v1remote "github.com/google/go-containerregistry/pkg/v1/remote"
)

// Registry access configuration.
type RegistryConfig struct {
	Authenticator         Authenticator
	ConcurrentJobs        int
	ConnectTimeout        time.Duration
	ForceAttemptHttp2     bool
	IdleConnectionTimeout time.Duration
	Insecure              bool
	KeepAlive             time.Duration
	MaxIdleConnections    int
	ResponseTimeout       time.Duration
	TlsHandshakeTimeout   time.Duration
}

// Create a new registry configuration containing default settings.
func NewRegistryConfig() RegistryConfig {
	return RegistryConfig{
		Authenticator:         NewAnonymousAuthenticator(),
		ConcurrentJobs:        4,
		ConnectTimeout:        16 * time.Second,
		ForceAttemptHttp2:     true,
		IdleConnectionTimeout: 1 * time.Hour,
		Insecure:              false,
		KeepAlive:             1 * time.Hour,
		MaxIdleConnections:    10,
		ResponseTimeout:       16 * time.Second,
		TlsHandshakeTimeout:   16 * time.Second,
	}
}

func (config *RegistryConfig) createDialer() *net.Dialer {
	return &net.Dialer{
		KeepAlive: config.KeepAlive,
		Timeout:   config.ConnectTimeout,
	}
}

func (config *RegistryConfig) createRoundTripper() *http.Transport {
	return &http.Transport{
		DialContext:           config.createDialer().DialContext,
		ForceAttemptHTTP2:     config.ForceAttemptHttp2,
		IdleConnTimeout:       config.IdleConnectionTimeout,
		MaxIdleConns:          config.MaxIdleConnections,
		ResponseHeaderTimeout: config.ResponseTimeout,
		TLSHandshakeTimeout:   config.TlsHandshakeTimeout,
	}
}

func (config *RegistryConfig) createDockerOptions() []v1remote.Option {
	return []v1remote.Option{
		v1remote.WithAuth(config.Authenticator.createDockerAuthenticator()),
		v1remote.WithJobs(config.ConcurrentJobs),
		v1remote.WithTransport(config.createRoundTripper()),
	}
}
