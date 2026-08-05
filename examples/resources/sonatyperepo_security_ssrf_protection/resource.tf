resource "sonatyperepo_security_ssrf_protection" "ssrf" {
  enabled = true
  allowed_domains = [
    "internal.example.com"
  ]
  allowed_ips = [
    "10.0.0.5"
  ]
}
