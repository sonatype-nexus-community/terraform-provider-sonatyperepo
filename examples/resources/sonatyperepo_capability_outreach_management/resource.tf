resource "sonatyperepo_capability_outreach_management" "cap" {
  enabled = true
  notes   = ""
  properties = {
    always_remote = false
    override_url  = "https://some.url.tld"
  }
}