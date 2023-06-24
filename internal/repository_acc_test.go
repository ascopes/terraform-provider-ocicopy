package internal_test

import (
	"testing"

	"github.com/ascopes/terraform-provider-ocicopy/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/lithammer/dedent"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"ocicopy": providerserver.NewProtocol6WithError(internal.NewOciCopyProvider()),
	}

	providerConfig = dedent.Dedent(`
		provider "ocicopy" {
			registries {
				basic_auth {
					registry_url = "my.registry.url/myregistrypath"
					username     = "some_username"
					password     = "some_password"
				}
				basic_auth {
					registry_url = "my.registry.url/myregistrypath2"
					username     = "some_username"
					password     = "some_password"
				}
				bearer_auth {
					registry_url = "my.registry.url/myregistrypath3"
					token	     = "1a2b3c"
				}
				bearer_auth {
					registry_url = "my.registry.url/myregistrypath3"
					token	     = "1a2b3c"
				}
				ecr {
					registry_url = "https://012345678910.dkr.ecr.eu-west-2.amazonaws.com"
					token	     = "ecr-token-1234"
				}
			}
		}
	`)
)

func TestAccRepositoryResource_declare(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: providerConfig + dedent.Dedent(`
				resource "ocicopy_repository" "hello_world" {

				}
			`),
			Check: resource.ComposeAggregateTestCheckFunc(),
		}},
	})
}
