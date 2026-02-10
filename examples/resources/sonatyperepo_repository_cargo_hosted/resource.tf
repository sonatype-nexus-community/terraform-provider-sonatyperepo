resource "sonatyperepo_repository_cargo_hosted" "cargo_hosted" {
  name   = "cargo-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
