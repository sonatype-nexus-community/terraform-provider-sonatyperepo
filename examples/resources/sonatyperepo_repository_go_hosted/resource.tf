resource "sonatyperepo_repository_go_hosted" "go_hosted" {
  name   = "go-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }

  component = {
    proprietary_components = false
  }
}
