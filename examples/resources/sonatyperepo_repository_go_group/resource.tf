resource "sonatyperepo_repository_go_group" "go_group" {
  name   = "go-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "go-proxy-repo"
    ]
  }
}
