resource "sonatyperepo_capability_ui_settings" "ui_settings" {
  enabled = true
  notes   = "UI settings configuration"
  properties = {
    title = "My Sonatype Nexus Repository"
  }
}
