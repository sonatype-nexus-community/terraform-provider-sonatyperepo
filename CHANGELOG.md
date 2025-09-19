<!-- See https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#changelog-specification -->

## X.Y.Z (UNRELEASED)

FEATURES:

* **New Resource:** `sonatyperepo_cleanup_policy` [GH-58]
* **New Resource:** `sonatyperepo_security_saml` [GH-63]

## 0.4.0 September 3, 2025

FEATURES:

* **New Resource:** `sonatyperepo_security_realms` [GH-60]

## 0.3.0 July 29, 2025

ENHANCEMENTS:

* Confirm Sonatype Nexus Repository Manager is WRITABLE [GH-54]
* Determine version of Sonatype Nexus Repository Manager for future use in this Provider [GH-54]

NOTES:
* Tested against [Sonatype Nexus Repository Manager 3.82.0](https://help.sonatype.com/en/sonatype-nexus-repository-3-82-0-release-notes.html)

## 0.2.0 July 17, 2025

FEATURES:

* **New Resource:** `sonatyperepo_repository_apt_hosted` [GH-35]
* **New Resource:** `sonatyperepo_repository_apt_proxy` [GH-35]
* **New Resource:** `sonatyperepo_repository_cargo_group` [GH-50]
* **New Resource:** `sonatyperepo_repository_cargo_hosted` [GH-50]
* **New Resource:** `sonatyperepo_repository_cargo_proxy` [GH-50]
* **New Resource:** `sonatyperepo_repository_conan_group` [GH-51]
* **New Resource:** `sonatyperepo_repository_conan_hosted` [GH-51]
* **New Resource:** `sonatyperepo_repository_conan_proxy` [GH-51]
* **New Resource:** `sonatyperepo_repository_cocoapods_proxy` [GH-39]
* **New Resource:** `sonatyperepo_repository_composer_proxy` [GH-40]
* **New Resource:** `sonatyperepo_repository_conda_proxy` [GH-45]
* **New Resource:** `sonatyperepo_repository_docker_group` [GH-36]
* **New Resource:** `sonatyperepo_repository_docker_hosted` [GH-36]
* **New Resource:** `sonatyperepo_repository_docker_proxy` [GH-36]
* **New Resource:** `sonatyperepo_repository_gitlfs_hosted` [GH-46]
* **New Resource:** `sonatyperepo_repository_go_group` [GH-41]
* **New Resource:** `sonatyperepo_repository_go_hosted` [GH-41]
* **New Resource:** `sonatyperepo_repository_helm_hosted` [GH-42]
* **New Resource:** `sonatyperepo_repository_helm_proxy` [GH-42]
* **New Resource:** `sonatyperepo_repository_huggingface_proxy` [GH-43]
* **New Resource:** `sonatyperepo_repository_nuget_group` [GH-38]
* **New Resource:** `sonatyperepo_repository_nuget_hosted` [GH-38]
* **New Resource:** `sonatyperepo_repository_nuget_proxy` [GH-38]
* **New Resource:** `sonatyperepo_repository_p2_proxy` [GH-47]
* **New Resource:** `sonatyperepo_repository_pypi_group` [GH-37]
* **New Resource:** `sonatyperepo_repository_pypi_hosted` [GH-37]
* **New Resource:** `sonatyperepo_repository_pypi_proxy` [GH-37]
* **New Resource:** `sonatyperepo_repository_r_group` [GH-49]
* **New Resource:** `sonatyperepo_repository_r_hosted` [GH-49]
* **New Resource:** `sonatyperepo_repository_r_proxy` [GH-49]
* **New Resource:** `sonatyperepo_repository_ruby_gems_group` [GH-48]
* **New Resource:** `sonatyperepo_repository_ruby_gems_hosted` [GH-48]
* **New Resource:** `sonatyperepo_repository_ruby_gems_proxy` [GH-48]
* **New Resource:** `sonatyperepo_repository_yum_group` [GH-52]
* **New Resource:** `sonatyperepo_repository_yum_hosted` [GH-52]
* **New Resource:** `sonatyperepo_repository_yum_proxy` [GH-52]

## 0.1.0 July 11, 2025

This is the first MVP release.

FEATURES:

* Resources and Datasources defined in [GH-1]
