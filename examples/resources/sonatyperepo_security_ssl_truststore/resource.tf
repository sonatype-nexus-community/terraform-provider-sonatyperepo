resource "sonatyperepo_security_ssl_truststore" "my_ca" {
  pem = file("${path.module}/certs/MY_CA.pem")
}
