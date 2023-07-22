package api

import (
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Create a default registry client config which can be further customized with optional
// settings.
func NewRegistryClientConfig(url string) *RegistryClientConfig {
	return &RegistryClientConfig{
		auth:                  authn.Anonymous,
		forceAttemptHttp2:     true,
		idleConnectionTimeout: time.Hour,
		jobs:                  4,
		keepAlive:             time.Hour,
		maxIdleConnections:    10,
		timeout:               16 * time.Second,
		tlsHandshakeTimeout:   16 * time.Second,
		url:                   url,
	}
}

type RegistryClientConfig struct {
	auth                  authn.Authenticator
	forceAttemptHttp2     bool
	idleConnectionTimeout time.Duration
	jobs                  int
	keepAlive             time.Duration
	maxIdleConnections    int
	timeout               time.Duration
	tlsHandshakeTimeout   time.Duration
	url                   string
}

func (cfg *RegistryClientConfig) GetOptions() []remote.Option {
	dialer := &net.Dialer{
		Timeout:   cfg.timeout,
		KeepAlive: cfg.keepAlive,
	}

	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     cfg.forceAttemptHttp2,
		MaxIdleConns:          cfg.maxIdleConnections,
		IdleConnTimeout:       cfg.idleConnectionTimeout,
		TLSHandshakeTimeout:   cfg.tlsHandshakeTimeout,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return []remote.Option{
		remote.WithAuth(cfg.auth),
		remote.WithJobs(cfg.jobs),
		remote.WithTransport(transport),
	}
}

func (cfg *RegistryClientConfig) SetBasicAuth(username string, password string) {
	cfg.auth = &authn.Basic{
		Username: username,
		Password: password,
	}
}

func (cfg *RegistryClientConfig) SetBearerAuth(token string) {
	cfg.auth = &authn.Bearer{Token: token}
}

func (cfg *RegistryClientConfig) SetForceAttemptHttp2(forceAttemptHttp2 bool) {
	cfg.forceAttemptHttp2 = forceAttemptHttp2
}

func (cfg *RegistryClientConfig) SetIdleConnectionTimeout(idleConnectionTimeout time.Duration) {
	cfg.idleConnectionTimeout = idleConnectionTimeout
}

func (cfg *RegistryClientConfig) SetJobs(jobs int) {
	cfg.jobs = jobs
}

func (cfg *RegistryClientConfig) SetKeepAlive(keepAlive time.Duration) {
	cfg.keepAlive = keepAlive
}

func (cfg *RegistryClientConfig) SetMaxIdleConnections(maxIdleConnections int) {
	cfg.maxIdleConnections = maxIdleConnections
}

func (cfg *RegistryClientConfig) SetTimeout(timeout time.Duration) {
	cfg.timeout = timeout
}

func (cfg *RegistryClientConfig) SetTlsHandshakeTimeout(tlsHandshakeTimeout time.Duration) {
	cfg.tlsHandshakeTimeout = tlsHandshakeTimeout
}
