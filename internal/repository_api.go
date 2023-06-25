package internal

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func determineDigestsForTags(ctx context.Context, repositoryName string, tags []string, registries map[string]registryModel) (map[string]string, []apiError) {
	repository, apiErrors := parseRepository(repositoryName)
	if apiErrors != nil {
		return nil, apiErrors
	}

	puller, apiErrors := createPullerFor(repository.Registry.Name(), registries)
	if apiErrors != nil {
		return nil, apiErrors
	}

	return fetchRawDigestsForTags(ctx, puller, *repository, tags)
}

func parseRepository(repositoryName string) (*name.Repository, []apiError) {
	repository, err := name.NewRepository(repositoryName)
	if err != nil {
		return nil, singleApiError(err, "Failed to parse repository '%s'", repositoryName)
	}
	return &repository, nil
}

func createPullerFor(registry string, registries map[string]registryModel) (*remote.Puller, []apiError) {
	options := []remote.Option{
		remote.WithJobs(15),
	}

	if configuredRegistry, ok := registries[registry]; ok {
		options = append(
			options,
			remote.WithAuth(configuredRegistry.getAuthenticator()),
		)
	}

	puller, err := remote.NewPuller(options...)
	if err != nil {
		return nil, singleApiError(err, "Failed to create internal API client for registry '%s'", registry)
	}

	return puller, nil
}

func fetchRawDigestsForTags(
	ctx context.Context,
	puller *remote.Puller,
	repository name.Repository,
	tags []string,
) (map[string]string, []apiError) {
	digests := make(map[string]string)
	apiErrors := make([]apiError, 0)
	repositoryName := repository.Name()

	existingTagsSlice, err := puller.List(ctx, repository)

	if err != nil {
		return nil, singleApiError(err, "Failed to list tags for repository '%s'", repositoryName)
	}

	existingTags := sliceToSet(existingTagsSlice)

	for _, tag := range tags {
		if _, ok := existingTags[tag]; !ok {
			apiErrors = append(apiErrors, newApiError(nil, "Tag '%s' does not exist for repository '%s'", tag, repositoryName))
			continue
		}

		tagRef := repository.Tag(tag)

		tagDescriptor, err := puller.Get(ctx, tagRef)
		if err != nil {
			apiErrors = append(apiErrors, newApiError(err, "Failed to fetch descriptors for ref '%s:%s'", repositoryName, tag))
			continue
		}

		digests[tag] = tagDescriptor.Digest.String()
	}

	return digests, apiErrors
}
