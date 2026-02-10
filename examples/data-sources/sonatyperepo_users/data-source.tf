data "sonatyperepo_users" "all" {
}

output "users" {
  value = data.sonatyperepo_users.all.users
}
