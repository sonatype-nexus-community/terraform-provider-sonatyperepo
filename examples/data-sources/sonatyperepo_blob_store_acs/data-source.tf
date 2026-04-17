data "sonatyperepo_blob_store_acs" "example" {
  name = "acs-blob-store"
}

output "acs_blob_store" {
  value = data.sonatyperepo_blob_store_acs.example
}
