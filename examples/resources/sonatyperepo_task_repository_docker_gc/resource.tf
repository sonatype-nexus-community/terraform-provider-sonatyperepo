resource "sonatyperepo_task_repository_docker_gc" "docker_gc" {
  name                   = "docker-gc-task"
  enabled                = true
  alert_email            = "admin@example.com"
  notification_condition = "FAILURE"

  frequency = {
    schedule = "daily"
  }

  properties = {
    repository_name = "docker-hosted-repo"
  }
}
