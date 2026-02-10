# Import Local user
terraform import sonatyperepo_user.saml_user my-uid,DEFAULT

# Import LDAP user
terraform import sonatyperepo_user.saml_user my-ldap-uid,LDAP

# Import SAML user
terraform import sonatyperepo_user.saml_admin my-user-id,SAML