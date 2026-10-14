resource "sonatyperepo_security_realms" "realms" {
  active = [
    "NexusAuthenticatingRealm",
    "DefaultRole"
  ]
}
