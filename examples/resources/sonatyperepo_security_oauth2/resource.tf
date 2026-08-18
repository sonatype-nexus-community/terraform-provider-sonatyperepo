resource "sonatyperepo_security_oauth2" "oauth2_config" {
  idp_jws_algorithm     = "RS256"
  username_claim        = "preferred_username"
  first_name_claim      = "given_name"
  last_name_claim       = "family_name"
  email_claim           = "email"
  groups_claim          = "groups"
  use_trust_store       = false
  client_id             = "nexus-repository-client"
  client_secret         = "changeme"
  idp_authorization_url = "https://idp.example.com/oauth2/authorize"
  idp_token_url         = "https://idp.example.com/oauth2/token"
  idp_logout_url        = "https://idp.example.com/oauth2/logout"
  idp_jwks_url          = "https://idp.example.com/oauth2/jwks"

  # Optional, provider-specific extras
  authorization_custom_params = {
    prompt = "consent"
  }
  exact_match_claims = {
    aud = "nexus-repository-client"
  }
}
