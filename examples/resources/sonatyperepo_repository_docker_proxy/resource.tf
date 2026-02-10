resource "sonatyperepo_repository_docker_proxy" "docker_proxy" {
  name   = "docker-proxy-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy = {
    remote_url       = "https://registry-1.docker.io"
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache = {
    enabled      = true
    time_to_live = 1440
  }

  http_client = {
    blocked                   = false
    auto_block                = false
    connection                = "0"
    enable_circular_redirects = true
    enable_cookies            = true
    retries                   = 0
    timeout                   = 60
    use_trust_store           = false
  }

  docker = {
    force_basic_auth = false
    v1_enabled       = false
  }

  docker_proxy = {
    index_type   = "REGISTRY"
    index_url    = "https://index.docker.io"
    use_nexus_ci = false
  }
}
