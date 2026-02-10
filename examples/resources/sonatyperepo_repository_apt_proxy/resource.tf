resource "sonatyperepo_repository_apt_proxy" "apt_proxy" {
  name   = "apt-proxy-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy = {
    remote_url       = "https://archive.ubuntu.com/ubuntu/"
    content_max_age  = 1439
    metadata_max_age = 1439
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

  apt = {
    distribution = "bionic"
  }
}
