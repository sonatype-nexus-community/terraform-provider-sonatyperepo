# Import an existing 'repository.maven.remove-snapshots' Task into Terraform State.

# Example
terraform import sonatyperepo_task_repository_maven_remove_snapshots.task TASK_ID

# Note: the public REST API does not return `properties` or full `frequency` for
# a Task, so the next `terraform plan` will show those fields as a diff against
# your configuration. The first `terraform apply` after import re-asserts them
# in Nexus to match your HCL.
