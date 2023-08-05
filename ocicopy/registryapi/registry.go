package registryapi

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/ocicopy/set"
)

// Base interface for a registry client. The implementation is kept separate to allow easy
// mocking and stubbing without physical API calls being performed.
type RegistryClient interface {
	// Transfer an image into another registry.
	//
	// The other registry client should be the same implementation type as this registry client.
	//
	// Returns the digest of the copied image, or an error if something goes wrong.
	CopyInto(
		ctx context.Context,
		sourceRepositoryName string,
		sourceDigest string,
		destinationClient RegistryClient,
		destinationRepositoryName string,
		destinationTagName string,
	) (string, error)

	// Fetch the digest for the given tag in the given repository.
	//
	// Returns the string digest, or an error if something goes wrong.
	//
	// Optionally, provide a platform to target a different platform to the default
	// (linux/amd64).
	GetDigestForTag(ctx context.Context, repositoryName string, tagName string, platformName *string) (string, error)

	// List tags for the given repository in the registry. The platform name can be nil
	// if the default should be used (linux/amd64).
	//
	// Returns a set of tags, or an error if something goes wrong.
	ListTags(ctx context.Context, repositoryName string, platformName *string) (set.Set[string], error)
}
