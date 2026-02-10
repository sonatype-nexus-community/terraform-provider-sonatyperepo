resource "sonatyperepo_repository_r_group" "r_group" {
  name   = "r-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "r-proxy-repo",
      "r-hosted-repo"
    ]
  }
}
