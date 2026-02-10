resource "sonatyperepo_privilege_repository_content_selector" "cs_privilege" {
  name             = "cs-privilege-example"
  description      = "Example content selector privilege"
  actions          = ["BROWSE"]
  content_selector = "test-content-selector"
  format           = "maven2"
  repository       = "maven-central"
}
