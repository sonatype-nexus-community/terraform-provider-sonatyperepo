resource "sonatyperepo_repository_rubygems_group" "rubygems_group" {
  name   = "rubygems-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "rubygems-proxy-repo",
      "rubygems-hosted-repo"
    ]
  }
}
