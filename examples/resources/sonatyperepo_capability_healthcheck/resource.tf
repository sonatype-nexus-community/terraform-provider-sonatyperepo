resource "sonatyperepo_capability_healthcheck" "cap" {
  notes   = "These are notes from Terraform"
  enabled = true
  properties = {
    configured_for_all_proxies = true
    use_nexus_truststore       = false
  }
}