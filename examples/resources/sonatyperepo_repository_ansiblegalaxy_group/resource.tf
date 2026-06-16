resource "sonatyperepo_repository_ansiblegalaxy_group" "ansiblegalaxy_group" {
  name   = "ansiblegalaxy-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_names = [
      "ansiblegalaxy-proxy-repo",
      "ansiblegalaxy-hosted-repo"
    ]
  }
}
