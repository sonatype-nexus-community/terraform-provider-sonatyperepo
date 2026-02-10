resource "sonatyperepo_repository_maven2_group" "maven_group" {
  name   = "maven-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  group = {
    member_repositories = [
      "maven-proxy-repo",
      "maven-hosted-repo"
    ]
  }

  maven = {
    version_policy = "RELEASE"
    layout_policy  = "STRICT"
  }
}
