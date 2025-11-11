resource "sonatyperepo_capability_webhook_global" "webhook" {
  enabled = true
  notes   = "These are notes from Terraform"
  properties = {
    names = [
      "repository"
    ]
    url    = "https://test.tld"
    secret = "super-secret-key"
  }
}