package containers

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const image = "docker.io/library/registry:2"
const primaryPort = "5000"
const primaryPortNat = nat.Port(primaryPort + "/tcp")

// Create a new object that can start a test container for a Docker registry server.
func NewRegistryTestContainer(t *testing.T) *RegistryTestContainer {
	enforceSeLinuxPermissive()

	return &RegistryTestContainer{
		name:          fmt.Sprintf("registry-%s", uuid.NewString()),
		containerImpl: nil,
		test:          t,
		LogLevel:      "info",
	}
}

// The registry test container implementation.
type RegistryTestContainer struct {
	name          string
	containerImpl testcontainers.Container
	test          *testing.T
	LogLevel      string
}

// Start the internal container.
func (registryContainer *RegistryTestContainer) Start(ctx context.Context) {
	if registryContainer.containerImpl != nil {
		return
	}

	registryContainer.test.Logf("Starting Docker Registry testcontainer (%s)", registryContainer.name)

	req := testcontainers.ContainerRequest{
		Name:         registryContainer.name,
		Image:        image,
		ExposedPorts: []string{primaryPort},
		WaitingFor:   wait.ForListeningPort(primaryPortNat),
		Env: map[string]string{
			// Docs: https://docs.docker.com/registry/configuration/

			// make HTTP auth ignore actual credentials, only check the presence of an
			// Authorization header.
			"REGISTRY_AUTH_SILLY_REALM":                          "testcontainers",
			"REGISTRY_AUTH_SILLY_SERVICE":                        registryContainer.name,
			"REGISTRY_HTTP_HTTP2_DISABLED":                       "false",
			"REGISTRY_LOG_ACCESSLOG_ENABLED":                     "true",
			"REGISTRY_LOG_LEVEL":                                 registryContainer.LogLevel,
			"REGISTRY_STORAGE_DELETE_ENABLED":                    "true",
			"REGISTRY_STORAGE_MAINTENANCE_UPLOADPURGING_ENABLED": "false",
		},
	}

	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
		Logger:           testcontainers.TestLogger(registryContainer.test),
		Reuse:            true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericReq)
	if err != nil {
		panic(err)
	}

	registryContainer.containerImpl = container

	if err = container.Start(ctx); err != nil {
		panic(err)
	}
}

// Stop the internal container.
func (registryContainer *RegistryTestContainer) Stop(ctx context.Context) {
	// Copy to prevent a race condition.
	container := registryContainer.containerImpl
	if container != nil {
		if err := registryContainer.containerImpl.Terminate(ctx); err != nil {
			panic(err)
		}
	}
}

func (registryContainer *RegistryTestContainer) HostVisibleEndpoint(ctx context.Context) string {
	port, err := registryContainer.containerImpl.MappedPort(ctx, primaryPortNat)
	if err != nil {
		panic(err.Error())
	}
	endpoint := fmt.Sprintf("localhost:%s", port.Port())
	return endpoint
}
