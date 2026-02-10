resource "sonatyperepo_privilege_repository_view" "view_privilege" {
  name        = "view-privilege-example"
  description = "Example repository view privilege"
  actions = [
    "BROWSE",
    "READ"
  ]
  format     = "maven2"
  repository = "maven-central"
}
