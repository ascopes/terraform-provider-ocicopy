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
			registry {
				registry_url = "123456789012.eu-west-2.dkr.ecr.amazonaws.com"
				basic_auth {
					username = "AWS"
					password = "ecr-token-1234"
				}
			}
		}
	`)
)

func TestAccRepositoryResource_hello_world(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: providerConfig + dedent.Dedent(`
				resource "ocicopy_repository" "hello_world" {
					from {
						name = "docker.io/hello-world"
						tags {
							values = ["latest"]
						}
					}
					to {
						name = "123456789012.eu-west-2.dkr.ecr.amazonaws.com"
					}
				}
			`),
			//Check: func(state *terraform.State) error {
			//	panic(state.String())
			//},
			Check: resource.TestCheckResourceAttrSet("ocicopy_repository.hello_world", "from.tags.0.digests.latest"),
		}},
	})
}
