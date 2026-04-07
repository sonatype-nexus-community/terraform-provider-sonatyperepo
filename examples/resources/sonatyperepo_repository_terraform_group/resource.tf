resource "sonatyperepo_repository_terraform_group" "terraform_group" {
  name   = "terraform-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = false
  }

  group = {
    member_names = ["terraform-hosted"]
  }

  terraform = {
    require_authentication = false
  }

  depends_on = [
    sonatyperepo_repository_terraform_hosted.terraform-hosted
  ]
}
