resource "sonatyperepo_repository_nuget_group" "nuget_group" {
  name   = "nuget-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  group = {
    member_repositories = [
      "nuget-proxy-repo",
      "nuget-hosted-repo"
    ]
  }
}
