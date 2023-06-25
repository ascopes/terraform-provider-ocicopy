package internal

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewRepositoryResource(provider *OciCopyProvider) resource.Resource {
	return &RepositoryResource{Provider: provider}
}

var _ resource.Resource = &RepositoryResource{}

type RepositoryResource struct {
	Provider *OciCopyProvider
}

func (resource *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	plan := &repositoryModel{}
	diags := req.Plan.Get(ctx, plan)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Generate the new ID.
	plan.populateTags(ctx, resource.Provider.Registries, &res.Diagnostics)
	plan.populateId()

	diags = res.State.Set(ctx, plan)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (resource *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	state := &repositoryModel{}
	diags := req.State.Get(ctx, state)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	state.populateId()

	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (*RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	plan := &repositoryModel{}
	diags := req.Plan.Get(ctx, plan)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	plan.populateId()

	diags = res.State.Set(ctx, plan)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (*RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	// Just remove it from the state. We don't have anything to physically delete here as we
	// keep the images.
	res.State.RemoveResource(ctx)
}

func (*RepositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_repository"
}

func (*RepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Declares a repository to copy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An internal identifier for this resource",
				PlanModifiers: []planmodifier.String{
					// This ID will change depending on the tag hashes, so ignore it for now.
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"from": schema.SingleNestedBlock{
				Description: "The repository to pull images from. This block is required",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The full name of the repository to copy",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"tags": schema.SetNestedBlock{
						Description: "Provide a set of tag constraints to copy. You must provide at least one of these blocks",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"values": schema.SetAttribute{
									Description: "Set tags to transfer.",
									ElementType: types.StringType,
									Required:    true,
									Validators: []validator.Set{
										setvalidator.IsRequired(),
										setvalidator.SizeAtLeast(1),
									},
								},
								"digests": schema.MapAttribute{
									Computed:    true,
									Description: "Mapping of the tags to their expected digest values.",
									ElementType: types.StringType,
								},
							},
						},
						Validators: []validator.Set{
							setvalidator.IsRequired(),
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
			"to": schema.SingleNestedBlock{
				Description: "Define a repository to transfer images to",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The full name of the registry to copy to",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		},
	}
}

type repositoryModel struct {
	From repositoryFromModel `tfsdk:"from" json:"from"`
	To   repositoryToModel   `tfsdk:"to"   json:"to"`

	// Generated attributes.
	Id types.String `tfsdk:"id"`
}

type repositoryFromModel struct {
	Name types.String          `tfsdk:"name" json:"name"`
	Tags []repositoryTagsModel `tfsdk:"tags" json:"tags"`
}

type repositoryToModel struct {
	Name types.String `tfsdk:"name" json:"name"`
}

type repositoryTagsModel struct {
	Values []types.String `tfsdk:"values"`

	// In the future, I may want to add other constraints such as semantic versioning
	// ranges, etc.

	// Generated attributes.
	Digests types.Map `tfsdk:"digests" json:"digests"`
}

func (repository *repositoryModel) populateTags(ctx context.Context, configuredRegistries registriesModel, diagnostics *diag.Diagnostics) {
	for i, tags := range repository.From.Tags {
		tagStrings := make([]string, len(tags.Values))
		for j, tag := range tags.Values {
			tagStrings[j] = tag.ValueString()
		}

		tagDigestMapping, apiErrors := determineDigestsForTags(
			ctx,
			repository.From.Name.ValueString(),
			tagStrings,
			configuredRegistries,
		)

		if processApiErrors(diagnostics, "Failed to fetch digests for tags", apiErrors...) {
			// Fail later.
			continue
		}

		serializableTagDigestMapping, diags := types.MapValueFrom(ctx, types.StringType, tagDigestMapping)
		diagnostics.Append(diags...)
		repository.From.Tags[i].Digests = serializableTagDigestMapping
	}
}

// Create a sentinel hash to use for the ID that is unique enough to used per resource
// in Terraform.
func (repository *repositoryModel) populateId() {
	// Unreachable error, should always be successful.
	jsonRepr, _ := json.Marshal(&repository)

	hash := sha512.New()
	hash.Write(jsonRepr)
	encodedHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	repository.Id = types.StringValue(encodedHash)
}
