resource "sonatyperepo_repository_raw_group" "raw_group" {
  name   = "raw-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "raw-proxy-repo",
      "raw-hosted-repo"
    ]
  }
}
