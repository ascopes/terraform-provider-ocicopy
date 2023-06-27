locals {
  local_registry = "localhost:5000"
}

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    ocicopy = {
      source = "hashicorp.com/ascopes/ocicopy"
    }
  }
}

provider "ocicopy" {
  registry {
    url = local.local_registry
  }
}

resource "ocicopy_repository" "hello_world" {
  from {
    image = "hello-world"
    tag   = "latest"
  }
  to {
    image = "${local.local_registry}/docker.io/hello-world"
  }
}