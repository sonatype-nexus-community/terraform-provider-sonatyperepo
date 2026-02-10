resource "sonatyperepo_privilege_repository_admin" "repo_admin_privilege" {
  name        = "repo-admin-privilege-example"
  description = "Example repository admin privilege"
  actions = [
    "BROWSE",
    "READ"
  ]
  format     = "maven2"
  repository = "maven-central"
}
