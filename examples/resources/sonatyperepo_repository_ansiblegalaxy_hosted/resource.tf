resource "sonatyperepo_repository_ansiblegalaxy_hosted" "ansiblegalaxy_hosted" {
  name   = "ansiblegalaxy-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
