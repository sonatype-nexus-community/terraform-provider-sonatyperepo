resource "sonatyperepo_repository_pub_hosted" "pub_hosted" {
  name   = "pub-hosted-repo"
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
