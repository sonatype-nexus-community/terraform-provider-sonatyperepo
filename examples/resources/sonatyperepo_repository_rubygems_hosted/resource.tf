resource "sonatyperepo_repository_rubygems_hosted" "rubygems_hosted" {
  name   = "rubygems-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
