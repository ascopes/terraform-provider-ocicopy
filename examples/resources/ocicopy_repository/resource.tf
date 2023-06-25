// Copies "docker.io/hello-world:latest" from DockerHub to our local ECR registry.

resource "aws_ecr_repository" "hello_world" {
  name                 = "docker.io/hello-world"
  force_delete         = true
  image_tag_mutability = "MUTABLE"
  scan_on_push         = true
}

resource "ocicopy_repository" "hello_world" {
  from {
    name = "hello-world"
    tags {
      values = ["latest"]
    }
  }
 
  to {
    name = aws_ecr_repository.hello_world.repository_url
  }
}
