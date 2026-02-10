resource "sonatyperepo_system_anonymous_access" "anonymous_access" {
  enabled    = true
  realm_name = "NexusAuthorizingRealm"
  user_id    = "anonymous"
}
