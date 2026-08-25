resource "sonatyperepo_repository_pub_group" "pub_group" {
  name   = "pub-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_names = [
      "pub-proxy-repo"
    ]
  }
}
