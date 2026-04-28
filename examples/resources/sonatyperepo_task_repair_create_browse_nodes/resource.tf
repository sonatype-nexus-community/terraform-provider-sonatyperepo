resource "sonatyperepo_task_repair_create_browse_nodes" "repair_browse_nodes" {
  name                   = "repair-browse-nodes-task"
  enabled                = true
  notification_condition = "FAILURE"

  frequency = {
    schedule = "manual"
  }

  properties = {
    repository_name = "*"
  }
}
