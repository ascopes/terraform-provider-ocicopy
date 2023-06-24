package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewRepositoryResource(provider *OciCopyProvider) resource.Resource {
	return &RepositoryResource{Provider: provider}
}

var _ resource.Resource = &RepositoryResource{}

type RepositoryResource struct {
	Provider *OciCopyProvider
}

func (*RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	plan := &RepositoryModel{}

	diags := req.Plan.Get(ctx, plan)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	plan.Id = types.StringValue("1234")

	diags = res.State.Set(ctx, plan)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

}

func (*RepositoryResource) Delete(_ context.Context, _ resource.DeleteRequest, res *resource.DeleteResponse) {
}

func (*RepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_repository"
}

func (*RepositoryResource) Read(_ context.Context, _ resource.ReadRequest, res *resource.ReadResponse) {
}

func (*RepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Declare a repository resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (*RepositoryResource) Update(_ context.Context, _ resource.UpdateRequest, res *resource.UpdateResponse) {
}

type RepositoryModel struct {
	Id types.String `tfsdk:"id"`
}
