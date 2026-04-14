# Import an existing S3 Blob Store into Terraform State.

# Example
terraform import sonatyperepo_blob_store_s3.s3_bs BLOB_STORE_NAME

# Note: API never returns the AWS Secret Access Key - so this will still show as a change requiring a `terraform apply` to be run