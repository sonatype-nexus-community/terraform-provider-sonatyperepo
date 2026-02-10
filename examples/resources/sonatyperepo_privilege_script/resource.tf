resource "sonatyperepo_privilege_script" "script_privilege" {
  name        = "script-privilege-example"
  description = "Example script privilege"
  actions = [
    "BROWSE",
    "RUN"
  ]
  script_name = "example-script"
}
