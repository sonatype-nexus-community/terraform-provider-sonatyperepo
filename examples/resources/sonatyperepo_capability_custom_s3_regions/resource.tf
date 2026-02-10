resource "sonatyperepo_capability_custom_s3_regions" "custom_regions" {
  enabled = true
  notes   = "Custom S3 regions configuration"
  properties = {
    regions = [
      "us-west-1",
      "eu-west-1",
      "ap-southeast-1"
    ]
  }
}
