resource "sonatyperepo_blob_store_google_cloud" "gcs_blob_store" {
  name = "gcs-blob-store"

  bucket_configuration = {
    bucket = "my-nexus-gcs-bucket"
    prefix = "nexus"
  }

  bucket_authentication = {
    authentication_type = "service_account"
    json_key_file       = file("/path/to/gcs-service-account-key.json")
  }

  encryption = {
    encryption_type = "gcs"
  }
}
