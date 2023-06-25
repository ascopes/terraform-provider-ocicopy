package internal

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

type registryModel interface {
	GetRegistryUrl() string
	GetAuthenticator() authn.Authenticator
}
