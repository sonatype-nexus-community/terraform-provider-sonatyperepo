resource "sonatyperepo_blob_store_acs" "blob_store" {
  name = "my-acs-blobs-tore"
  bucket_configuration = {
    account_name   = "AZURE-STORAGE-ACCOUNT-NAME"
    container_name = "AZURE-CONTAINER-NAME" # Will be created if doesn't exist
    authentication = {
      authentication_method = "ACCOUNTKEY"
      account_key           = "YOUR-SECRET-ACCOUNT-KEY-HERE"
    }
  }
}