package registrycontainer

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/opencontainers/selinux/go-selinux"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Create a new object that can start a test container for a Docker registry server.
func NewRegistryTestContainer() *RegistryTestContainer {
	if runtime.GOOS == "linux" && selinux.EnforceMode() == selinux.Enforcing {
		panic(
			"SELinux is enabled and set to 'Enforcing'. testcontainers-ryuk cannot start under this setting. " +
				"Please run 'sudo setenforce permissive' to lower the SELinux enforcement level to allow " +
				"Ryuk to start.",
		)
	}

	return &RegistryTestContainer{
		Debug:         false,
		name:          fmt.Sprintf("registry-%s", uuid.NewString()),
		containerImpl: nil,
	}
}

// The registry test container implementation.
type RegistryTestContainer struct {
	Debug         bool
	name          string
	containerImpl testcontainers.Container
}

// Start the internal container.
func (registryContainer *RegistryTestContainer) Start(ctx context.Context) {
	if registryContainer.containerImpl != nil {
		return
	}

	var logLevel string
	if registryContainer.Debug {
		logLevel = "debug"
	} else {
		logLevel = "info"
	}

	req := testcontainers.ContainerRequest{
		Name:         registryContainer.name,
		Image:        "docker.io/library/registry:2",
		ExposedPorts: []string{"5000"},
		WaitingFor:   wait.ForListeningPort("5000/tcp"),
		Env: map[string]string{
			// Docs: https://docs.docker.com/registry/configuration/
			
			// make HTTP auth ignore actual credentials, only check the presence of an
			// Authorization header.
			"REGISTRY_AUTH_SILLY_REALM": "testcontainers",
			"REGISTRY_AUTH_SILLY_SERVICE": registryContainer.name,
			"REGISTRY_HTTP_HTTP2_DISABLED": "false",
			"REGISTRY_LOG_ACCESSLOG_ENABLED": "true",
			"REGISTRY_LOG_LEVEL": logLevel,
			"REGISTRY_STORAGE_DELETE_ENABLED": "true",
			"REGISTRY_STORAGE_MAINTENANCE_UPLOADPURGING_ENABLED": "false",
		},
	}

	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericReq)
	if err != nil {
		panic(err)
	}

	registryContainer.containerImpl = container
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
	port, err := registryContainer.containerImpl.MappedPort(ctx, "5000/tcp")
	if err != nil {
		panic(err.Error())
	}
	endpoint := fmt.Sprintf("localhost:%s", port.Port())
	return endpoint
}
