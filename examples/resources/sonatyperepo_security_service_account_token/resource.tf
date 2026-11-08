# The role that decides what the token may do
resource "sonatyperepo_role" "ci" {
  id          = "ci-deployer"
  name        = "ci-deployer"
  description = "Deploys artifacts from CI"
  privileges  = ["nx-repository-view-*-*-read"]
}

# A token that never expires
resource "sonatyperepo_security_service_account_token" "ci" {
  name        = "ci-pipeline"
  role_id     = sonatyperepo_role.ci.id
  description = "Used by the CI pipeline"
}

# A token that expires after 90 days. Changing any argument replaces the token
# with a new one - Nexus Repository cannot update a token in place.
resource "sonatyperepo_security_service_account_token" "short_lived" {
  name            = "ci-pipeline-short-lived"
  role_id         = sonatyperepo_role.ci.id
  expiration_days = 90
}

# The token value is sensitive and is only available because Terraform created
# the token - it cannot be read back from Nexus Repository later.
output "ci_token" {
  value     = sonatyperepo_security_service_account_token.ci.token
  sensitive = true
}

# Tokens only authenticate while the ServiceAccountRealm security realm is
# active - include it in sonatyperepo_security_realms, alongside the realms you
# already use:
#
# resource "sonatyperepo_security_realms" "realms" {
#   active = [
#     "NexusAuthenticatingRealm",
#     "ServiceAccountRealm"
#   ]
# }
