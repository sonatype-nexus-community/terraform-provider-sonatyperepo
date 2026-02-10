data "sonatyperepo_blob_store_file" "example" {
  name = "default"
}

output "file_blob_store" {
  value = data.sonatyperepo_blob_store_file.example
}
