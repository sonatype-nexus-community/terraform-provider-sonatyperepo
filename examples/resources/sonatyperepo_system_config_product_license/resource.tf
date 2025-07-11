resource "sonatyperepo_system_config_product_license" "license" {
  license_data = filebase64("/path/to/sonatype.lic")
}