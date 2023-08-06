GENERATED_DIRS 		 := docs/ site/
GENERATED_CACHE_DIRS := .venv/

.PHONY: help
help:
	@-echo "terraform-provider-ocicopy"
	@-echo "=========================="
	@-echo ""
	@-echo "Usage: make [<target> ...]"
	@-echo "Run a build step."
	@-echo ""
	@-echo "Supported targets:"
	@-echo ""
	@-echo "  build     - Compile the provider and any tests."
	@-echo "  clean     - Delete cached files from builds."
	@-echo "  deepclean - Same as clean, but clears internal Go caches as well."
	@-echo "  docs      - Produce HTML documentation from generated Terraform Provider "
	@-echo "              Markdown documentation (via mkdocs)."
	@-echo "  generate  - Generate any sources needed for builds.
	@-echo "  get       - Download any build dependencies and cache them."
	@-echo "  install   - Install the provider into \$$GOBIN so it can be used by Terraform."
	@-echo "  test      - Run unit, integration, and acceptance tests."
	@-echo "  vet       - Run any linters."

.PHONY: build
build: get vet generate
	go fmt ./...
	go mod tidy
	go build ./...

.PHONY: clean
clean:
	go clean -i
	$(RM) -R $(GENERATED_DIRS)

.PHONY: deepclean
deepclean: clean
	go clean -r -testcache
	$(RM) -R $(GENERATED_CACHE_DIRS)

.PHONY: docs
docs: .venv/bin/mkdocs generate
	.venv/bin/mkdocs build

.PHONY: generate
generate: get
	if [[ ! -d docs/ ]]; then mkdir docs/; fi
	go generate ./...

.PHONY: get
get:
	go get -x

.PHONY: install
install: build
	go install -x internal/...

.PHONY: test
test: build
	TF_ACC=1 go test ./...

.PHONY: vet
vet: get
	go vet ./...

.venv/bin/mkdocs:
	python3 -m venv .venv
	.venv/bin/pip install mkdocs-material
