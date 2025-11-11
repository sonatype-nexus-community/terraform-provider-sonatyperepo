resource "sonatyperepo_capability_base_url" "cap_base_url" {
  notes   = "These are notes from Terraform"
  enabled = true
  properties = {
    url = "https://my.fqdn.tld"
  }
}