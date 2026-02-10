resource "sonatyperepo_task_repository_docker_upload_purge" "docker_upload_purge" {
  name                   = "docker-upload-purge-task"
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
