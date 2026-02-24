# Creating a Local User
resource "sonatyperepo_user" "test_user" {
  user_id       = "test-local-user"
  first_name    = "Testing"
  last_name     = "User"
  email_address = "test@local.user"
  # If you wish to manage a user without it's password, set password = null
  password = "somethingSecurer"
  status   = "active"
  roles = [
    "nx-anonymous"
  ]
}

# Defining a LDAP/CROWD/SAML User to bring into Terraform State
#
# See Import below also.
resource "sonatyperepo_user" "saml_admin" {
  # Local Roles can be assigned and managed for non-local users
  roles = ["nx-test-role"]

  # These fields cannot be managed by Terraform - they are readonly for non-local Users
  user_id       = "admin"
  email_address = "admin@somewhere.tld"
  first_name    = "SAML"
  last_name     = "Admin"
  status        = "active"
}
