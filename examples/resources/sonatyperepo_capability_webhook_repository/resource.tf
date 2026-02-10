resource "sonatyperepo_capability_webhook_repository" "webhook_repo" {
  enabled = true
  notes   = "Repository webhook configuration"
  properties = {
    names      = ["component"]
    url        = "https://webhook.example.com/receive"
    secret     = "[REDACTED:secret]"
    repository = "maven-central"
  }
}
