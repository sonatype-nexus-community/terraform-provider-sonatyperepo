resource "sonatyperepo_repository_yum_hosted" "yum_hosted" {
  name   = "yum-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
