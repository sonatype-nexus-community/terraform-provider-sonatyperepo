resource "sonatyperepo_repository_alpine_group" "alpine_group" {
  name   = "alpine-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_names = [
      "alpine-proxy-repo",
      "alpine-hosted-repo"
    ]
  }

  alpine = {
    key_pair   = "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
    passphrase = "changeit"
  }
}
