resource "sonatyperepo_capability_webhook_global" "webhook_global" {
  enabled = true
  notes   = "Global webhook configuration"
  properties = {
    names  = ["repository"]
    url    = "https://webhook.example.com/receive"
    secret = "[REDACTED:secret]"
  }
}
