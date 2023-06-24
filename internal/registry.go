package internal

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

type RegistryModel interface {
	GetRegistryUrl() string
	GetAuthenticator() authn.Authenticator
}
