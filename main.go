package main

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	ctx := context.Background()
	opts := providerserver.ServeOpts{Address: "github.com/ascopes/terraform-provider-ocicopy"}
	providerserver.Serve(ctx, provider.NewProvider, opts)
}

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate
//go:vet      go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate
