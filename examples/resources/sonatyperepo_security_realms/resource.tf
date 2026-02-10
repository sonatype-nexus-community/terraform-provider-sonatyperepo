resource "sonatyperepo_security_realms" "realms" {
  active = [
    "NexusAuthorizingRealm",
    "DefaultRole"
  ]
}
