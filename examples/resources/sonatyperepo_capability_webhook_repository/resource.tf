resource "sonatyperepo_capability_webhook_repository" "webhook" {
  enabled = true
  notes   = "These are notes from Terraform"
  properties = {
    names      = ["component"]
    url        = "https://test.tld"
    secret     = "testing"
    repository = "maven-central"
  }
}
