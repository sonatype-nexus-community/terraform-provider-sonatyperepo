resource "sonatyperepo_capability_custom_s3_regions" "cap" {
  enabled = true
  notes   = "These are notes from Terraform"
  properties = {
    regions = ["somewhere-1", "somewhere-2"]
  }
}