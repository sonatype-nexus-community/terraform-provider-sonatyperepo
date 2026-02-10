resource "sonatyperepo_capability_ui_branding" "ui_branding" {
  enabled = true
  notes   = "UI branding configuration"
  properties = {
    header_enabled = true
    header_html    = "<h1>Welcome to My Nexus Repository</h1>"
    footer_enabled = false
  }
}
