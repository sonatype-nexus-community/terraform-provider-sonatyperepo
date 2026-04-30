# Import an existing 'license.expiration.notification' Task into Terraform State.

# Example
terraform import sonatyperepo_task_license_expiration_notification.task TASK_ID

# Note: the public REST API does not return full `frequency` for a Task, so the
# next `terraform plan` will show that field as a diff against your
# configuration. The first `terraform apply` after import re-asserts it in
# Nexus to match your HCL.
