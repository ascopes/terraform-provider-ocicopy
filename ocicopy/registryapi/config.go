package registryapi

import (
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Registry access configuration.
type RegistryConfig struct {
	authenticator         Authenticator
	concurrentJobs        int
	connectTimeout        time.Duration
	forceAttemptHttp2     bool
	idleConnectionTimeout time.Duration
	keepAlive             time.Duration
	maxIdleConnections    int
	responseTimeout       time.Duration
	tlsHandshakeTimeout   time.Duration
}

// Create a new registry configuration containing default settings.
func NewRegistryConfig() RegistryConfig {
	return RegistryConfig{
		authenticator:         NewAnonymousAuthenticator(),
		concurrentJobs:        4,
		connectTimeout:        16 * time.Second,
		forceAttemptHttp2:     true,
		idleConnectionTimeout: 1 * time.Hour,
		keepAlive:             1 * time.Hour,
		maxIdleConnections:    10,
		responseTimeout:       16 * time.Second,
		tlsHandshakeTimeout:   16 * time.Second,
	}
}

func (config RegistryConfig) createDialer() *net.Dialer {
	return &net.Dialer{
		KeepAlive: config.keepAlive,
		Timeout:   config.connectTimeout,
	}
}

func (config RegistryConfig) createRoundTripper() *http.Transport {
	return &http.Transport{
		DialContext:           config.createDialer().DialContext,
		ForceAttemptHTTP2:     config.forceAttemptHttp2,
		IdleConnTimeout:       config.idleConnectionTimeout,
		MaxIdleConns:          config.maxIdleConnections,
		ResponseHeaderTimeout: config.responseTimeout,
		TLSHandshakeTimeout:   config.tlsHandshakeTimeout,
	}
}

func (config RegistryConfig) createDockerOptions() []remote.Option {
	return []remote.Option{
		remote.WithAuth(config.authenticator.createDockerAuthenticator()),
		remote.WithJobs(config.concurrentJobs),
		remote.WithTransport(config.createRoundTripper()),
	}
}
