resource "sonatyperepo_repository_pypi_group" "pypi_group" {
  name   = "pypi-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "pypi-proxy-repo",
      "pypi-hosted-repo"
    ]
  }
}
