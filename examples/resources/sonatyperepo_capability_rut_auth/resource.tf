resource "sonatyperepo_capability_rut_auth" "cap_base_url" {
  notes   = "These are notes from Terraform"
  enabled = true
  properties = {
    http_header = "X-Your-Custom-Header-Name"
  }
}