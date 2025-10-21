resource "sonatyperepo_task_blobstore_compact" "task_bc" {
  name                   = "test-blobstore-compact"
  enabled                = true
  alert_email            = ""
  notification_condition = "FAILURE"
  frequency = {
    schedule = "manual"
  }
  properties = {
    blob_store_name  = "default"
    blobs_older_than = 99
  }
}