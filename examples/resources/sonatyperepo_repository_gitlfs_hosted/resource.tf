resource "sonatyperepo_repository_gitlfs_hosted" "gitlfs_hosted" {
  name   = "gitlfs-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
