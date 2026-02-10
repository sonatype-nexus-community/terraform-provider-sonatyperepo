data "sonatyperepo_content_selector" "example" {
  name = "my-content-selector"
}

output "content_selector" {
  value = data.sonatyperepo_content_selector.example
}
