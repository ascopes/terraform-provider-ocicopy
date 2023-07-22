package api

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Description of an image and how to access it across a network.
type Image struct {
	config     *RegistryClientConfig
	registry   string
	repository string
	platform   *string
}

// Create a new image descriptor.
func NewImage(config *RegistryClientConfig, registry string, repository string) *Image {
	return &Image{config, registry, repository, nil}
}

// Get the remote descriptor for the image with the given digest.
func (image *Image) GetImageDescriptorFromDigest(ctx context.Context, digestValue string) (v1.Image, error) {
	repositoryObj, err := image.repositoryObj()
	if err != nil {
		return nil, err
	}

	digest := repositoryObj.Digest(digestValue)
	if err != nil {
		return nil, err
	}

	return image.getImageDescriptor(ctx, digest)
}

// Get the remote descriptor for the image with the given tag.
func (image *Image) GetImageDescriptorFromTag(ctx context.Context, tagName string) (v1.Image, error) {
	repositoryObj, err := image.repositoryObj()
	if err != nil {
		return nil, err
	}

	tag := repositoryObj.Tag(tagName)
	if err != nil {
		return nil, err
	}

	return image.getImageDescriptor(ctx, tag)
}

// Get all existing tags for the object, mapped to their digests.
func (image *Image) GetExistingTags(ctx context.Context) ([]string, error) {
	repositoryObj, err := image.repositoryObj()
	if err != nil {
		return nil, err
	}

	puller, err := image.puller()
	if err != nil {
		return nil, err
	}

	return puller.List(ctx, repositoryObj)
}

// Upload an image and assign it the given tag. Return the uploaded digest.
func (dest *Image) UploadImageAsTag(ctx context.Context, tagName string, src v1.Image) (string, error) {
	pusher, err := dest.pusher()
	if err != nil {
		return "", err
	}

	repo, err := dest.repositoryObj()
	if err != nil {
		return "", err
	}

	tag := repo.Tag(tagName)
	err = pusher.Push(ctx, tag, src)
	if err != nil {
		return "", err
	}

	puller, err := dest.puller()
	if err != nil {
		return "", err
	}

	desc, err := puller.Head(ctx, tag)
	if err != nil {
		return "", err
	}

	return desc.Digest.String(), nil
}

// Delete the given image from the remote.
func (dest *Image) DeleteImage(ctx context.Context, digestValue string) error {
	pusher, err := dest.pusher()
	if err != nil {
		return err
	}

	repo, err := dest.repositoryObj()
	if err != nil {
		return err
	}

	tag := repo.Digest(digestValue)
	return pusher.Delete(ctx, tag)
}

func (image *Image) getImageDescriptor(ctx context.Context, ref name.Reference) (v1.Image, error) {

	puller, err := image.puller()
	if err != nil {
		return nil, err
	}

	desc, err := puller.Get(ctx, ref)
	if err != nil {
		return nil, err
	}

	return desc.Image()
}

// Create an image puller.
func (image *Image) puller() (*remote.Puller, error) {
	options, err := image.options()
	if err != nil {
		return nil, err
	}

	// TODO: cache result.
	return remote.NewPuller(options...)
}

// Create an image pusher or return it if it already exists.
func (image *Image) pusher() (*remote.Pusher, error) {
	options, err := image.options()
	if err != nil {
		return nil, err
	}

	// TODO: cache result.
	return remote.NewPusher(options...)
}

func (image *Image) options() ([]remote.Option, error) {
	options := image.config.GetOptions()

	if image.platform != nil {
		actualPlatform, err := v1.ParsePlatform(*image.platform)
		if err != nil {
			return nil, err
		}

		options = append(options, remote.WithPlatform(*actualPlatform))
	}

	return options, nil
}

func (image *Image) repositoryObj() (name.Repository, error) {
	registryObj, err := name.NewRegistry(image.registry)

	if err != nil {
		return name.Repository{}, err
	}

	return registryObj.Repo(image.repository), nil
}
