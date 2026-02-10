resource "sonatyperepo_repository_docker_group" "docker_group" {
  name   = "docker-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  group = {
    member_repositories = [
      "docker-proxy-repo",
      "docker-hosted-repo"
    ]
  }

  docker = {
    force_basic_auth = true
  }
}
