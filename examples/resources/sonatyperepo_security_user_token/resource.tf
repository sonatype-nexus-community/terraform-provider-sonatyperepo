resource "sonatyperepo_security_user_token" "user_token" {
  user_id = "admin"
  name    = "api-token"
}
