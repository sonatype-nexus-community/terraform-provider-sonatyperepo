resource "sonatyperepo_capability_firewall_audit_and_quarantine" "cap" {
  notes   = "These are notes from Terraform"
  enabled = true
  properties = {
    repository = "maven-central"
    quarantine = true
  }
}