package registryapi

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Representation of credentials for registry level authentication.
type Authenticator interface {
	createDockerAuthenticator() authn.Authenticator
}

// Create an unauthenticated registry authenticator.
func NewAnonymousAuthenticator() Authenticator {
	return anonymousAuthenticatorInstance
}

// Create a basic auth authenticator.
func NewBasicAuthenticator(username string, password string) Authenticator {
	return basicAuthenticator{username, password}
}

// Create a bearer auth token.
func NewBearerAuthenticator(token string) Authenticator {
	return bearerAuthenticator{token}
}

var anonymousAuthenticatorInstance = anonymousAuthenticator{}

type anonymousAuthenticator struct{}

func (auth anonymousAuthenticator) createDockerAuthenticator() authn.Authenticator {
	return authn.Anonymous
}

type basicAuthenticator struct {
	username string
	password string
}

func (auth basicAuthenticator) createDockerAuthenticator() authn.Authenticator {
	return &authn.Basic{Username: auth.username, Password: auth.password}
}

type bearerAuthenticator struct {
	token string
}

func (auth bearerAuthenticator) createDockerAuthenticator() authn.Authenticator {
	return &authn.Bearer{Token: auth.token}
}
