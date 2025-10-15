resource "sonatyperepo_content_selector" "cs1" {
  name        = "test-content-selector"
  description = "This is a test content selector"
  expression  = "format == \"maven2\" and path =^ \"/org/sonatype/sub\""
}