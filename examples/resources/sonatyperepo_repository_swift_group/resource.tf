resource "sonatyperepo_repository_swift_group" "swift_group" {
  name   = "swift-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_names = [
      "swift-proxy-repo"
    ]
  }
}
