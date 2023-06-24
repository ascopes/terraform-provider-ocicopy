package main

import (
	"context"

	"github.com/ascopes/terraform-provider-ocicopy/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	ctx := context.Background()
	opts := providerserver.ServeOpts{
		Address: "github.com/ascopes/terraform-provider-ocicopy",
	}
	providerserver.Serve(ctx, internal.NewOciCopyProvider, opts)
}
