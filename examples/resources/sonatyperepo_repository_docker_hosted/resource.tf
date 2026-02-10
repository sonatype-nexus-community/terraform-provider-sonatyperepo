resource "sonatyperepo_repository_docker_hosted" "docker_hosted" {
  name   = "docker-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }

  docker = {
    force_basic_auth = true
    v1_enabled       = false
    http_port        = 8082
  }
}
