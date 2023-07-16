package api

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Create a puller client for the given configuration.
func pullerFor(config RegistryClientConfig, platform *string) (*remote.Puller, error) {
	options := config.GetOptions()

	if platform != nil {
		actualPlatform, err := v1.ParsePlatform(*platform)
		if err != nil {
			return nil, err
		}

		options = append(options, remote.WithPlatform(*actualPlatform))
	}

	return remote.NewPuller(options...)
}
