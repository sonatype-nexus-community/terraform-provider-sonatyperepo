resource "sonatyperepo_repository_nuget_hosted" "nuget_hosted" {
  name   = "nuget-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
