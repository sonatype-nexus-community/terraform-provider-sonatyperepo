resource "sonatyperepo_blob_store_file" "blob_store" {
  name = "file-blob-store"
  path = "/mnt/nexus-storage"
}
