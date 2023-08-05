package registryapi

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/ocicopy/set"
	"github.com/google/go-containerregistry/pkg/name"
	v1api "github.com/google/go-containerregistry/pkg/v1"
	v1remote "github.com/google/go-containerregistry/pkg/v1/remote"
)

// Registry client implementation. This is separate from the interface so we can easily mock/stub the
// internal calls elsewhere in this project.
type registryClientImpl struct {
	registry       name.Registry
	registryConfig *RegistryConfig
}

// Initialize a new registry client.
//
// If the registry name is unparsable, then an error is returned instead.
func NewRegistryClient(registryName string, registryConfig *RegistryConfig) (RegistryClient, error) {
	registry, err := name.NewRegistry(registryName)
	if err != nil {
		return nil, err
	}
	return &registryClientImpl{registry, registryConfig}, nil
}

func (sourceClient *registryClientImpl) CopyInto(
	ctx context.Context,
	sourceRepositoryName string,
	sourceDigest string,
	destinationClient RegistryClient,
	destinationRepositoryName string,
	destinationTagName string,
) (string, error) {
	destinationClientImpl, ok := destinationClient.(*registryClientImpl)
	if !ok {
		panic("Unexpected registry client impl type for destination client")
	}

	sourceRef := sourceClient.registry.Repo(sourceRepositoryName).Digest(sourceDigest)
	destinationRef := destinationClientImpl.registry.Repo(destinationRepositoryName).Tag(destinationTagName)

	sourcePuller, err := sourceClient.constructPuller(nil)
	if err != nil {
		return "", err
	}

	sourceDesc, err := sourcePuller.Get(ctx, sourceRef)
	if err != nil {
		return "", err
	}

	destinationPusher, err := destinationClientImpl.constructPusher()
	if err != nil {
		return "", err
	}

	if err = destinationPusher.Push(ctx, destinationRef, sourceDesc); err != nil {
		return "", err
	}

	destinationPuller, err := destinationClientImpl.constructPuller(nil)
	if err != nil {
		return "", err
	}

	// TODO(ascopes): could eventual consistency create an issue here for some registries?
	// Is this guaranteed to always work?
	destinationDesc, err := destinationPuller.Head(ctx, destinationRef)
	if err != nil {
		return "", err
	}

	return destinationDesc.Digest.Hex, nil
}

func (client *registryClientImpl) GetDigestForTag(ctx context.Context, repositoryName string, tagName string, platformName *string) (string, error) {
	puller, err := client.constructPuller(platformName)
	if err != nil {
		return "", err
	}

	repository := client.registry.Repo(repositoryName)
	ref := repository.Tag(tagName)
	if err != nil {
		return "", err
	}

	desc, err := puller.Head(ctx, ref)
	if err != nil {
		return "", err
	}

	return desc.Digest.Hex, nil
}

func (client *registryClientImpl) ListTags(ctx context.Context, repositoryName string, platformName *string) (set.Set[string], error) {
	puller, err := client.constructPuller(platformName)
	if err != nil {
		return nil, err
	}

	repository := client.registry.Repo(repositoryName)
	tagSlice, err := puller.List(ctx, repository)
	if err != nil {
		return nil, err
	}

	tagSet := set.NewHashSet[string]()
	for _, tag := range tagSlice {
		tagSet.Add(tag)
	}

	return tagSet, nil
}

func (client *registryClientImpl) constructPuller(platformName *string) (*v1remote.Puller, error) {
	options := client.registryConfig.createDockerOptions()

	// For some reason we have to provide the platform name ahead of time, which is mildly
	// inconvenient because it makes this a bit more messy to work with.
	if platformName != nil {
		platform, err := v1api.ParsePlatform(*platformName)
		if err != nil {
			return nil, err
		}

		options = append(options, v1remote.WithPlatform(*platform))
	}

	return v1remote.NewPuller(options...)

}

func (client *registryClientImpl) constructPusher() (*v1remote.Pusher, error) {
	return v1remote.NewPusher(client.registryConfig.createDockerOptions()...)
}
