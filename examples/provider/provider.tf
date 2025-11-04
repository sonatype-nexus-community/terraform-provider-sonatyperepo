# Simplest Configuration
provider "sonatyperepo" {
  host     = "https://my-sonatype-nexus-repository.tld:port"
  username = "username"
  password = "password"
}

# If you run with a base path, you can add it:
provider "sonatyperepo" {
  host          = "https://my-sonatype-nexus-repository.tld:port"
  username      = "username"
  password      = "password"
  api_base_path = "/my-custom-base/service/rest"
}

# If you access via a Load Balancer or service that strips the `Server` header
# you can provide a hint as to the version of Sonatype Nexus Repository:
provider "sonatyperepo" {
  host         = "https://my-sonatype-nexus-repository.tld:port"
  username     = "username"
  password     = "password"
  version_hint = "3.85.0-03 (PRO)"
}
