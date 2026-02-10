data "sonatyperepo_blob_stores" "all" {
}

output "blob_stores" {
  value = data.sonatyperepo_blob_stores.all.blob_stores
}
