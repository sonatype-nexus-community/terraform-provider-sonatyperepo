resource "sonatyperepo_privilege_wildcard" "wildcard_privilege" {
  name        = "wildcard-privilege-example"
  description = "Example wildcard privilege"
  pattern     = "nx-repository-view-*-*-browse"
}
