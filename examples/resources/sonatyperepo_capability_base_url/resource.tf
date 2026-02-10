resource "sonatyperepo_capability_base_url" "base_url_capability" {
  enabled = true
  notes   = "Configured via Terraform"
  properties = {
    url = "https://repo.example.com"
  }
}
