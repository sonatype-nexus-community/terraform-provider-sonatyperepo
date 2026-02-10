resource "sonatyperepo_task_repository_maven_remove_snapshots" "maven_remove_snapshots" {
  name                   = "maven-remove-snapshots-task"
  enabled                = true
  alert_email            = "admin@example.com"
  notification_condition = "FAILURE"

  frequency = {
    schedule       = "weekly"
    recurring_days = [1] # Monday
  }

  properties = {
    repository_name  = "maven-hosted-repo"
    min_days_to_keep = 7
    regex            = ".*"
  }
}
