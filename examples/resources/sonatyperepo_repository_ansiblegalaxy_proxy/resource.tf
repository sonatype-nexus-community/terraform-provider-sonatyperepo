resource "sonatyperepo_repository_ansiblegalaxy_proxy" "ansiblegalaxy_proxy" {
  name   = "ansiblegalaxy-proxy-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  proxy = {
    remote_url       = "https://galaxy.ansible.com"
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache = {
    enabled      = true
    time_to_live = 1440
  }

  http_client = {
    blocked    = false
    auto_block = false
    connection = {
      enable_circular_redirects = true
      enable_cookies            = true
      retries                   = 3
      timeout                   = 60
      use_trust_store           = false
    }
  }
}
