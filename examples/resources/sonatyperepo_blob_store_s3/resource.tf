resource "sonatyperepo_blob_store_s3" "s3_blob_store" {
  name = "s3-blob-store"

  bucket_configuration = {
    bucket = "my-nexus-bucket"
    prefix = "nexus"
    region = "us-east-1"
  }

  bucket_security = {
    access_key_id     = "AKIAIOSFODNN7EXAMPLE"
    secret_access_key = "[REDACTED:secret]"
  }

  encryption = {
    encryption_type = "s3"
  }

  advanced_bucket_connection = {
    endpoint          = "https://s3.us-east-1.amazonaws.com"
    signature_version = "AWSSIGV4"
    path_style_access = false
    chunked_encoding  = false
  }
}
