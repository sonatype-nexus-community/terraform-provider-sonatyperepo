resource "sonatyperepo_repository_oci_hosted" "oci_hosted" {
  name   = "oci-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }

  oci = {
    force_basic_auth = true
    v1_enabled       = false
    http_port        = 8082
  }

  cosign = {
    enforcement = "NONE"
  }
}
