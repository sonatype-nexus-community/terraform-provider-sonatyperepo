data "sonatyperepo_privileges" "all" {
}

output "privileges" {
  value = data.sonatyperepo_privileges.all.privileges
}
