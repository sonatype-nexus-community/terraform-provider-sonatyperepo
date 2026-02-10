resource "sonatyperepo_repository_apt_hosted" "apt_hosted" {
  name   = "apt-hosted-repo"
  online = true

  storage = {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }

  apt = {
    distribution = "bionic"
  }

  apt_signing = {
    key_pair   = "-----BEGIN PGP PRIVATE KEY BLOCK-----\n...\n-----END PGP PRIVATE KEY BLOCK-----"
    passphrase = "[REDACTED:passphrase]"
  }
}
