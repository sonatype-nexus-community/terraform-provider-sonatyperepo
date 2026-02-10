data "sonatyperepo_blob_store_group" "example" {
  name = "group-blob-store"
}

output "group_blob_store" {
  value = data.sonatyperepo_blob_store_group.example
}
