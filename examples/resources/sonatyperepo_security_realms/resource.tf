resource "sonatyperepo_security_realms" "realms" {
  active = [
    "NexusAuthorizingRealm",
    "DefaultRole"
  ]
}

# Other realms, such as "ServiceAccountRealm" (needed for service account tokens, see
# sonatyperepo_security_service_account_token), can be added to the list. Realms are
# consulted in the order listed, and this list replaces the active realms on the server.
