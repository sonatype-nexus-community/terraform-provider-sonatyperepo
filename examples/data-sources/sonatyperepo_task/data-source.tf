data "sonatyperepo_task" "example" {
  id = "a1b2c3d4e5f6"
}

output "task" {
  value = data.sonatyperepo_task.example
}
