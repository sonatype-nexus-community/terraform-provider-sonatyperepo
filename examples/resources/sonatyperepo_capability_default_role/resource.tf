resource "sonatyperepo_capability_default_role" "cap" {
  enabled = true
  notes   = "These are notes from Terraform"
  properties = {
    role = "nx-anonymous"
  }
}