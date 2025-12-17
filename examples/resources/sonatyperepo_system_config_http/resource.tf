resource "sonatyperepo_system_config_http" "http" {
  http_proxy = {
    enabled = false
  }
  https_proxy = {
    enabled = true
    host    = "my.proxy.tld"
    port    = 8080
    authentication = {
      enabled  = true
      username = "proxy-user"
      password = "proxy-password"
    }
  }
}