# Import an existing Azure Cloud Storage Blob Store into Terraform State.
#
# Not if authenticating using Account Key, this cannot be read and will show as a change after import.

# Example
terraform import sonatyperepo_blob_store_acs.blob_store BLOB_STORE_FILE_NAME