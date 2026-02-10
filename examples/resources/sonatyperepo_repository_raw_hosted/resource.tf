resource "sonatyperepo_repository_raw_hosted" "raw_hosted" {
  name   = "raw-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
