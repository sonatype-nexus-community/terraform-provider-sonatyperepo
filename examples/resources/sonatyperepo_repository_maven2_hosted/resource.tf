resource "sonatyperepo_repository_maven2_hosted" "maven_hosted" {
  name   = "maven-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }

  maven = {
    version_policy = "RELEASE"
    layout_policy  = "STRICT"
  }
}
