resource "sonatyperepo_repository_npm_hosted" "npm_hosted" {
  name   = "npm-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
