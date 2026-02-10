resource "sonatyperepo_repository_r_hosted" "r_hosted" {
  name   = "r-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
