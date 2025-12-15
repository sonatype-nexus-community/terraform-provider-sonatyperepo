resource "sonatyperepo_task_license_expiration_notification" "task" {
  name                   = "my-notifciation"
  enabled                = true
  alert_email            = "notifications@somewhere.tld"
  notification_condition = "FAILURE"
  frequency = {
    schedule = "manual"
  }
}