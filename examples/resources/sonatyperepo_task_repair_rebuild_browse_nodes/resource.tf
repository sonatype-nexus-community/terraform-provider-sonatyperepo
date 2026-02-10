resource "sonatyperepo_task_repair_rebuild_browse_nodes" "repair_browse_nodes" {
  name                   = "repair-browse-nodes-task"
  enabled                = true
  notification_condition = "FAILURE"

  frequency = {
    schedule = "manual"
  }
}
