resource "sonatyperepo_privilege_application" "app_privilege" {
  name        = "app-privilege-example"
  description = "Example application privilege"
  domain      = "capabilities"
  actions = [
    "BROWSE",
    "READ",
    "EDIT"
  ]
}
