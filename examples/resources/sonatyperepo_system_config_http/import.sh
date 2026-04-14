# Existing System HTTP Configuration can be imported to state.
#
# NOTE: The Identifier (SECURITY_REALMS) in below example has no meaning and is just to comply with Terraform syntax.
#
# HTTP Proxy Passwords are not accessible via the REST API and will require a terrafom apply after terraform import to complete state.

# Example
terraform import sonatyperepo_system_config_http.config anything