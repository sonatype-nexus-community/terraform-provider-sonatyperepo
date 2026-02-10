resource "sonatyperepo_system_config_mail" "mail_config" {
  enabled                           = true
  host                              = "smtp.example.com"
  port                              = 587
  from_address                      = "nexus@example.com"
  ssl_on_connect_enabled            = false
  ssl_server_identity_check_enabled = true
  start_tls_enabled                 = true
  start_tls_required                = false
  nexus_trust_store_enabled         = false
  username                          = "nexus-user"
  password                          = "[REDACTED:password]"
  subject_prefix                    = "[Nexus]"
}
