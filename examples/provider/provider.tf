# Simplest Configuration
provider "sonatyperepo" {
  url      = "https://my-sonatype-nexus-repository.tld:port"
  username = "username"
  password = "password"
}

# Using environment variables for credentials (useful for CI/CD)
# Set NXRM_SERVER_URL, NXRM_SERVER_USERNAME, and NXRM_SERVER_PASSWORD
provider "sonatyperepo" {
  # Credentials provided via environment variables
}

# Mix environment variables with explicit configuration
provider "sonatyperepo" {
  username = "terraform-user"
  password = "terraform-password"
  # URL provided via NXRM_SERVER_URL environment variable
}

# If you run with a base path, you can add it:
provider "sonatyperepo" {
  url           = "https://my-sonatype-nexus-repository.tld:port"
  username      = "username"
  password      = "password"
  api_base_path = "/my-custom-base/service/rest"
}

# If you access via a Load Balancer or service that strips the `Server` header
# you can provide a hint as to the version of Sonatype Nexus Repository:
provider "sonatyperepo" {
  url          = "https://my-sonatype-nexus-repository.tld:port"
  username     = "username"
  password     = "password"
  version_hint = "3.89.1-01 (PRO)"
}
