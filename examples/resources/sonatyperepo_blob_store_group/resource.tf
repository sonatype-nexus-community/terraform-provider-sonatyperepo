// Simplest example
resource "sonatyperepo_blob_store_group" "test_group" {
  name        = "test-group"
  fill_policy = "roundRobin"
  members     = ["test1", "test2"]
}

// Depend on Blob Store resources
resource "sonatyperepo_blob_store_group" "test_group" {
  name        = "test-group"
  fill_policy = "writeToFirst"
  members = [
    sonatyperepo_blob_store.test1.name,
    sonatyperepo_blob_store.test2.name,
  ]
}
