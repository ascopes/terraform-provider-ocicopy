package registryapi

import (
	"reflect"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"gotest.tools/v3/assert"
)

func Test_NewAnonymousAuthenticator_Returns_anonymousAuthenticator(t *testing.T) {

	// When
	got := NewAnonymousAuthenticator()

	// Then
	_, ok := got.(anonymousAuthenticator)
	assert.Check(t, ok, "Expected anonymousAuthenticator, got %#v", reflect.TypeOf(got).Name())
}

func Test_NewBasicAuthenticator_Returns_basicAuthenticator(t *testing.T) {
	// Given
	wantedUsername := "foo-username"
	wantedPassword := "bar-password"

	// When
	got := NewBasicAuthenticator(wantedUsername, wantedPassword)

	// Then
	gotImpl, ok := got.(basicAuthenticator)
	assert.Check(t, ok, "Expected basicAuthenticator, got %#v", reflect.TypeOf(got).Name())
	assert.Equal(t, wantedUsername, gotImpl.username)
	assert.Equal(t, wantedPassword, gotImpl.password)
}

func Test_NewBearerAuthenticator_Returns_basicAuthenticator(t *testing.T) {
	// Given
	wantedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// When
	got := NewBearerAuthenticator(wantedToken)

	// Then
	gotImpl, ok := got.(bearerAuthenticator)
	assert.Check(t, ok, "Expected bearerAuthenticator, got %#v", reflect.TypeOf(got).Name())
	assert.Equal(t, wantedToken, gotImpl.token)
}

func Test_anonymousAuthenticator_createDockerAuthenticator_ReturnsExpectedValue(t *testing.T) {
	// Given
	authenticator := anonymousAuthenticator{}

	// When
	got := authenticator.createDockerAuthenticator()
	wanted := authn.Anonymous

	// Then
	assert.Equal(t, wanted, got, "Unexpected anonymous authenticator result")
}

func Test_basicAuthenticator_createDockerAuthenticator_ReturnsExpectedValue(t *testing.T) {
	// Given
	wantedUsername := "bar-username"
	wantedPassword := "baz-password-123"
	authenticator := basicAuthenticator{username: wantedUsername, password: wantedPassword}

	// When
	got, gotErr := authenticator.createDockerAuthenticator().Authorization()
	wanted, wantedErr := (&authn.Basic{Username: wantedUsername, Password: wantedPassword}).Authorization()

	// Then
	assert.NilError(t, gotErr, "Failed to init generated 'got' authenticator")
	assert.NilError(t, wantedErr, "Failed to init generated 'wanted' authenticator")
	assert.Equal(t, *wanted, *got, "Unexpected basic authenticator result")
}

func Test_bearerAuthenticator_createDockerAuthenticator_ReturnsExpectedValue(t *testing.T) {
	// Given
	wantedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	authenticator := bearerAuthenticator{token: wantedToken}

	// When
	got, gotErr := authenticator.createDockerAuthenticator().Authorization()
	wanted, wantedErr := (&authn.Bearer{Token: wantedToken}).Authorization()

	// Then
	assert.NilError(t, gotErr, "Failed to init generated 'got' authenticator")
	assert.NilError(t, wantedErr, "Failed to init generated 'wanted' authenticator")
	assert.Equal(t, *wanted, *got, "Unexpected bearer authenticator result")
}
