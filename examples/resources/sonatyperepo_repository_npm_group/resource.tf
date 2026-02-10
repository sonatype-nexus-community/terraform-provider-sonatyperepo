resource "sonatyperepo_repository_npm_group" "npm_group" {
  name   = "npm-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  group = {
    member_repositories = [
      "npm-proxy-repo",
      "npm-hosted-repo"
    ]
  }
}
