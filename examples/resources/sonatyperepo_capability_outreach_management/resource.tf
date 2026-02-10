resource "sonatyperepo_capability_outreach_management" "outreach_capability" {
  enabled = true
  notes   = "Outreach management configuration"
  properties = {
    always_remote = false
    override_url  = "https://outreach.example.com"
  }
}
