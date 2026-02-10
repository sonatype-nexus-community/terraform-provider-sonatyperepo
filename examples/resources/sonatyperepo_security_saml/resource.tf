resource "sonatyperepo_security_saml" "saml_config" {
  idp_metadata                 = file("/path/to/idp-metadata.xml")
  username_attribute           = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier"
  email_attribute              = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"
  first_name_attribute         = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"
  last_name_attribute          = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"
  groups_attribute             = "http://schemas.xmlsoap.org/claims/Group"
  validate_assertion_signature = true
  validate_response_signature  = true
}
