data "sonatyperepo_security_user_token" "config" {
}

output "security_user_token" {
  value = data.sonatyperepo_security_user_token.config
}
