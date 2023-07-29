package registryapi

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func Test_NewRegistryConfig_ReturnsExpectedDefaults(t *testing.T) {
	// When
	got := NewRegistryConfig()

	// Then
	t.Run("authenticator", func(t *testing.T) {
		assert.Equal(t, anonymousAuthenticatorInstance, got.authenticator)
	})
	t.Run("concurrentJobs", func(t *testing.T) {
		assert.Equal(t, 4, got.concurrentJobs)
	})
	t.Run("connectTimeout", func(t *testing.T) {
		assert.Equal(t, 16*time.Second, got.connectTimeout)
	})
	t.Run("forceAttemptHttp2", func(t *testing.T) {
		assert.Equal(t, true, got.forceAttemptHttp2)
	})
	t.Run("idleConnectionTimeout", func(t *testing.T) {
		assert.Equal(t, 1*time.Hour, got.idleConnectionTimeout)
	})
	t.Run("keepAlive", func(t *testing.T) {
		assert.Equal(t, 1*time.Hour, got.keepAlive)
	})
	t.Run("maxIdleConnections", func(t *testing.T) {
		assert.Equal(t, 10, got.maxIdleConnections)
	})
	t.Run("responseTimeout", func(t *testing.T) {
		assert.Equal(t, 16*time.Second, got.responseTimeout)
	})
	t.Run("tlsHandshakeTimeout", func(t *testing.T) {
		assert.Equal(t, 16*time.Second, got.tlsHandshakeTimeout)
	})
}

func Test_RegistryConfig_createDialer_ReturnsExpectedResult(t *testing.T) {
	// Given
	config := RegistryConfig{
		connectTimeout: 1 * time.Minute,
		keepAlive:      5 * time.Minute,
	}

	// When
	got := config.createDialer()

	// Then
	assert.Equal(t, 1*time.Minute, got.Timeout)
	assert.Equal(t, 5*time.Minute, got.KeepAlive)
}

func Test_RegistryConfig_createRoundTripper_ReturnsExpectedResult(t *testing.T) {
	// Given
	config := RegistryConfig{
		connectTimeout:        1 * time.Minute,
		forceAttemptHttp2:     true,
		idleConnectionTimeout: 10 * time.Minute,
		keepAlive:             12 * time.Second,
		maxIdleConnections:    32,
		responseTimeout:       31 * time.Second,
		tlsHandshakeTimeout:   13 * time.Second,
	}

	// When
	got := config.createRoundTripper()

	// Then
	t.Run("DialContext", func(t *testing.T) {
		// TODO: find a way to verify this is the result of createDialer.
		assert.Check(t, got.DialContext != nil, "Dialer is not set")
	})
	t.Run("ForceAttemptHTTP2", func(t *testing.T) {
		assert.Equal(t, true, got.ForceAttemptHTTP2)
	})
	t.Run("IdleConnTimeout", func(t *testing.T) {
		assert.Equal(t, 10*time.Minute, got.IdleConnTimeout)
	})
	t.Run("MaxIdleConns", func(t *testing.T) {
		assert.Equal(t, 32, got.MaxIdleConns)
	})
	t.Run("ResponseHeaderTimeout", func(t *testing.T) {
		assert.Equal(t, 31*time.Second, got.ResponseHeaderTimeout)
	})
	t.Run("TLSHandshakeTimeout", func(t *testing.T) {
		assert.Equal(t, 13*time.Second, got.TLSHandshakeTimeout)
	})
}

func Test_RegistryConfig_createDockerOptions_ReturnsExpectedSlice(t *testing.T) {
	// Given
	config := RegistryConfig{
		authenticator:         NewBasicAuthenticator("foo-usr", "bar-psw"),
		concurrentJobs:        12,
		connectTimeout:        1 * time.Minute,
		forceAttemptHttp2:     true,
		idleConnectionTimeout: 10 * time.Minute,
		keepAlive:             12 * time.Second,
		maxIdleConnections:    32,
		responseTimeout:       31 * time.Second,
		tlsHandshakeTimeout:   13 * time.Second,
	}

	// When
	got := config.createDockerOptions()

	// Then
	skipMessage := "Docker registry API does not allow us to inspect this at this time."

	t.Run("Len() is expected value", func(t *testing.T) {
		assert.Equal(t, 3, len(got))
	})
	t.Run("WithAuthenticator(Authenticator)", func(t *testing.T) { t.Skip(skipMessage) })
	t.Run("WithJobs(int)", func(t *testing.T) { t.Skip(skipMessage) })
	t.Run("WithTransport(*RoundTripper)", func(t *testing.T) { t.Skip(skipMessage) })
}
