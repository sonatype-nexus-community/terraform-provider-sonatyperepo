data "sonatyperepo_content_selectors" "all" {
}

output "content_selectors" {
  value = data.sonatyperepo_content_selectors.all.content_selectors
}
