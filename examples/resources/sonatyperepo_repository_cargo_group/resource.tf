resource "sonatyperepo_repository_cargo_group" "cargo_group" {
  name   = "cargo-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_repositories = [
      "cargo-proxy-repo",
      "cargo-hosted-repo"
    ]
  }
}
