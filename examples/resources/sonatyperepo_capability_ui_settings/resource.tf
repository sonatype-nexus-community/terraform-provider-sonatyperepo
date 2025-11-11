resource "sonatyperepo_capability_ui_settings" "cap" {
  enabled      = true
  last_updated = "Thursday, 06-Nov-25 14:11:02 GMT"
  notes        = "Automatically added on Wed Oct 16 12:27:43 GMT 2024"
  properties = {
    debug_allowed                 = true
    long_request_timeout          = 180
    request_timeout               = 60
    session_timeout               = 30
    status_interval_anonymous     = 60
    status_interval_authenticated = 5
    title                         = "Sonatype Nexus Repository"
  }
}