resource "sonatyperepo_blob_store_gcs" "blob_store" {
  name = "gcs-blob-store"

  bucket_configuration = {
    authentication = {
      account_key           = "ACCOUNT_KEY"
      authentication_method = "METHOD"
    }

    bucket = {
      name = "bucket_name"
    }
  }
}
