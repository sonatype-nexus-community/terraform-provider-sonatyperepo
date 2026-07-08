resource "sonatyperepo_repository_alpine_hosted" "alpine_hosted" {
  name   = "alpine-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }

  alpine = {
    key_pair   = "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
    passphrase = "changeit"
  }
}
