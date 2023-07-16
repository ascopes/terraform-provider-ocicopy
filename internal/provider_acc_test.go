package internal_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ascopes/terraform-provider-ocicopy/containers/registrycontainer"
	"github.com/ascopes/terraform-provider-ocicopy/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type config struct {
	content string
}

func terraformConfig(lines ...string) config {
	return config{
		content: strings.Join(lines, "\n"),
	}
}

func (cfg config) format(args ...any) string {
	return fmt.Sprintf(cfg.content, args...)
}

func createRegistryContainer(ctx context.Context, t *testing.T) *registrycontainer.RegistryTestContainer {
	container := registrycontainer.NewRegistryTestContainer()
	container.Start(ctx)
	return container
}

func TestAccProvider_CopyImageAcross(t *testing.T) {
	ctx := context.Background()
	container := createRegistryContainer(ctx, t)
	defer container.Stop(ctx)

	config := terraformConfig(
		"provider \"ocicopy\" {}",
		"",
		"resource \"ocicopy_image\" \"hello_world\" {",
		"  source = {",
		"    registry_url = \"%[1]s\"",
		"    image        = \"hello-world\"",
		"    tag          = \"latest\"",
		"  }",
		"  target = {",
		"    registry_url = \"%[2]s\"",
		"  }",
		"}",
	).format(
		"docker.io",
		container.HostVisibleEndpoint(ctx),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"ocicopy": providerserver.NewProtocol6WithError(internal.NewOciCopyProvider()),
		},
		Steps: []resource.TestStep{
			{Config: config},
		},
	})
}
