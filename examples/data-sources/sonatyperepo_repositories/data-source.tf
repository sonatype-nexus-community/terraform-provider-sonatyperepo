data "sonatyperepo_repositories" "all" {
}

output "repositories" {
  value = data.sonatyperepo_repositories.all.repositories
}
