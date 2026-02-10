resource "sonatyperepo_repository_conan_hosted" "conan_hosted" {
  name   = "conan-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
