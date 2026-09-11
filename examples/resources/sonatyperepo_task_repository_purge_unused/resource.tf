resource "sonatyperepo_task_repository_purge_unused" "purge_unused" {
  name                   = "purge-unused-task"
  enabled                = true
  alert_email            = "admin@example.com"
  notification_condition = "FAILURE"

  frequency = {
    schedule        = "weekly"
    start_date      = 1777167000000
    timezone_offset = "-05:00"
    # 1 = Sunday through to 7 = Saturday
    recurring_days = [1]
  }

  properties = {
    # Use "*" to purge unused components in all repositories
    repository_name = "maven-central"
    last_used       = 30
  }
}
