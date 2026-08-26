resource "sonatyperepo_repository_oci_group" "oci_group" {
  name   = "oci-group-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  group = {
    member_names = [
      "oci-proxy-repo",
      "oci-hosted-repo"
    ]
  }

  oci = {
    force_basic_auth = true
    v1_enabled       = false
  }
}
