data "sonatyperepo_tasks" "all" {
}

output "tasks" {
  value = data.sonatyperepo_tasks.all.tasks
}
