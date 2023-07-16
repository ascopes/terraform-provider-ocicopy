package api

type SourceImage struct {
	Config   RegistryClientConfig
	Image    string
	Platform *string
	Hash     string
}

type DestinationImage struct {
	Config RegistryClientConfig
	Image  string
	Tag    string
}

func TransferImage(src SourceImage, dest DestinationImage) error {
	_, err := pullerFor(src.Config, src.Platform)
	if err != nil {
		return err
	}

	_, err = pusherFor(dest.Config, src.Platform)
	if err != nil {
		return err
	}

	panic("unimplemented")
}
