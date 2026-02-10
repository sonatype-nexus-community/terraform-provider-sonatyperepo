resource "sonatyperepo_repository_pypi_hosted" "pypi_hosted" {
  name   = "pypi-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
