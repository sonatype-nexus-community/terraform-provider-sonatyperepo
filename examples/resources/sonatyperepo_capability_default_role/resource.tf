resource "sonatyperepo_capability_default_role" "default_role" {
  enabled = true
  notes   = "Default role configuration"
  properties = {
    role = "nx-anonymous"
  }
}
