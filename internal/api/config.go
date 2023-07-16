package api

import (
	"time"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Registry client configuration interface.
// I might remove this and just pass the struct around later. Depends how easy it is to
// unit test.
type RegistryClientConfig interface {
	SetBasicAuth(userName string, password string)
	SetBearerAuth(token string)
	SetForceAttemptHttp2(bool)
	SetIdleConnectionTimeout(time.Duration)
	SetJobs(int)
	SetKeepAlive(time.Duration)
	SetMaxIdleConnections(int)
	SetTimeout(time.Duration)
	SetTlsHandshakeTimeout(time.Duration)

	GetOptions() []remote.Option
}
