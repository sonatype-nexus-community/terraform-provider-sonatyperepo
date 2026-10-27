resource "sonatyperepo_system_anonymous_access" "anonymous_access" {
  enabled    = true
  realm_name = "NexusAuthenticatingRealm"
  user_id    = "anonymous"
}
