resource "sonatyperepo_capability_ui_branding" "cap" {
  enabled = true
  notes   = ""
  properties = {
    header_enabled = true
    header_html    = "TESTING 1 2 3"
  }
}