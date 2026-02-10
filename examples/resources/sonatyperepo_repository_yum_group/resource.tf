resource "sonatyperepo_repository_yum_group" "yum_group" {
  name   = "yum-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "yum-proxy-repo",
      "yum-hosted-repo"
    ]
  }
}
