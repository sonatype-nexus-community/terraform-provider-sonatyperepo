# Terraform Provider for Sonatype Nexus Repository Manager

[![shield_tfr-version]][link_tfr]
![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dm/572010)
[![shield_gh-workflow-test]][link_gh-workflow-test]
[![shield_license]][license_file] 
![shield_tf_version]
[![shield_nxrm_version]][link_nxrm_release]


---

This Terraform Provider allows you to use Configuration-as-Code (CasC) practises for managing the configuration of Sonatype Nexus Repository Manager.

This provider does not provide functionality for actually deploying Sonatype Nexus Repository Manager (i.e. Infrastructure or Application installation). For deployment and installation, see  the [official Help Documentation](https://help.sonatype.com/en/sonatype-nexus-repository.html).

## Version Support

We test this Provider against a range of Terraform versions and Sonatype Nexus Repository Manager versions as noted below.

### Sonatype Nexus Repository Manager

We test on `N - 5` releases (where possible, but no earlier than `3.79.1`). See [here](https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/blob/main/.github/workflows/test.yml) for the current list.

### Terraform Version support

We test on the latest patch release of each the earliest and latest version of Terraform i.e. `1.0.x` and  `1.12.x` - i.e. we aim to support all Terraform versions since `1.0.0`.

## Usage

See our [documentation](./docs/index.md) and the [examples directory](./examples/).

## Development

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## The Fine Print

Remember:

It is worth noting that this is **NOT SUPPORTED** by Sonatype, and is a contribution of ours to the open source community (read: you!)

* Use this contribution at the risk tolerance that you have
* Do NOT file Sonatype support tickets related to `terraform-provider-sonatyperepo`
* DO file issues here on GitHub, so that the community can pitch in

Phew, that was easier than I thought. Last but not least of all - have fun!


[shield_gh-workflow-test]: https://img.shields.io/github/actions/workflow/status/sonatype-nexus-community/terraform-provider-sonatyperepo/test.yml?branch=main&logo=GitHub&logoColor=white "build"
[shield_tfr-version]: https://img.shields.io/badge/Terraform%20Registry-0.2.0-8A2BE2
[shield_license]: https://img.shields.io/github/license/sonatype-nexus-community/terraform-provider-sonatyperepo?logo=open%20source%20initiative&logoColor=white "license"
[shield_tf_version]: https://img.shields.io/badge/Terraform-1.0.0+-blue
[shield_nxrm_version]: https://img.shields.io/badge/Sonatype_Nexus_Repository-3.79.1&nbsp;&ndash;&nbsp;3.82.0-blue

[link_tfr]: https://registry.terraform.io/providers/sonatype-nexus-community/sonatyperepo/latest
[link_gh-workflow-test]: https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/actions/workflows/test.yml?query=branch%3Amain
[license_file]: https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/blob/main/LICENSE
[link_nxrm_release]: https://help.sonatype.com/en/sonatype-nexus-repository-3-81-0-release-notes.html