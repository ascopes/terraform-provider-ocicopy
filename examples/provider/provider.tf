terraform {
  required_providers {
    ocicopy = {
      source  = "<tbc>/ocicopy"
      version = "~> <tbc>"
    }
  }

  required_version = ">= 1.4.0, < 2.0.0"
}

provider "ocicopy" {
  registry {
    basic_auth = {
      username = data.aws_ecr_authorization_token.token.username
      password = data.aws_ecr_authorization_token.token.password
    }
    url = data.aws_ecr_authorization_token.token.proxy_endpoint
  }
}

data "aws_ecr_authorization_token" "token" {
}