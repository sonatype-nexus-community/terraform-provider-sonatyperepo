# Existing OAuth2/OIDC configuration can be imported as follows.
#
# NOTE: The Identifier (OAUTH2) in below example has no meaning and is just to comply with Terraform syntax.
# NOTE: client_secret is never returned by the API, so it will not be populated after import -
#       set it explicitly in configuration afterwards.

# Example
terraform import sonatyperepo_security_oauth2.oauth2_config OAUTH2
