.PHONY: build
build:
	go build ./...

.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

.PHONY: docs
docs:
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	go generate

.PHONY: docs-html
docs-html: docs
	python3 -m venv .venv
	.venv/bin/pip install mkdocs-material
	.venv/bin/mkdocs build

.PHONY: download
download:
	go mod download

.PHONY: clean
clean:
	go clean -i -r -x
	$(RM) -Rv build/ docs/ site/ .venv/

test: build
	TF_ACC=1 TF_LOG=INFO go test -v ./...

.PHONY: rebuild
rebuild: clean build