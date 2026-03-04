// Create a role with privileges
resource "sonatyperepo_role" "example" {
  description = "Administrator role"
  id          = "nx-admin"
  name        = "nx-admin"
  privileges  = ["nx-all"]
}

// Create a role with roles
resource "sonatyperepo_role" "example" {
  description = "Docker Roles"
  id          = "docker-all"
  name        = "docker-all"
  roles       = ["some-docker-role", "some-other-docker-role"]
}
