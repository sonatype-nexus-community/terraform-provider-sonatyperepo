data "sonatyperepo_roles" "all" {
}

output "roles" {
  value = data.sonatyperepo_roles.all.roles
}
