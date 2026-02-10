resource "sonatyperepo_repository_conan_group" "conan_group" {
  name   = "conan-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "conan-proxy-repo",
      "conan-hosted-repo"
    ]
  }
}
