resource "sonatyperepo_capability_firewall_audit_and_quarantine" "firewall_capability" {
  enabled = true
  notes   = "Firewall audit and quarantine configuration"
  properties = {
    repository = "maven-central"
    quarantine = true
  }
}
