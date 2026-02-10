resource "sonatyperepo_repository_helm_hosted" "helm_hosted" {
  name   = "helm-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
