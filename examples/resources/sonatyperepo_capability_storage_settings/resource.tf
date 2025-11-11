resource "sonatyperepo_capability_storage_settings" "cap_storage_settings" {
  notes   = "These are notes from Terraform"
  enabled = true
  properties = {
    last_downloaded_interval = 24
  }
}