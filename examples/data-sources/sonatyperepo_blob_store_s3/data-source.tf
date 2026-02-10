data "sonatyperepo_blob_store_s3" "example" {
  name = "s3-blob-store"
}

output "s3_blob_store" {
  value = data.sonatyperepo_blob_store_s3.example
}
