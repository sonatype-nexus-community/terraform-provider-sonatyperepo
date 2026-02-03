# THIS IS TO BE FORMATTED

## Breaking Changes in 1.0.0

- FW config for proxy repositories

- Resource `sonatyperepo_repository_apt_hosted`: 
  - The `apt_signing` block is now required
  - The `apt_signing.passphrase` field is now optional
- Resource `sonatyperepo_repository_maven_group` has been renamed to `sonatyperepo_repository_maven2_group`
- Resource `sonatyperepo_repository_maven_hosted` has been renamed to `sonatyperepo_repository_maven2_hosted`
- Resource `sonatyperepo_repository_maven_proxy` has been renamed to `sonatyperepo_repository_maven2_proxy`

## Improvements / Behaviour Changes in 1.0.0

### Default Values Introduced
- Resources `sonatyperepo_repository_*_proxy`: 
  - `.proxy.content_max_age` now has a default value (`1440`) and does not need to be supplied unless you wish to use a different value
  - `.proxy.metadata_max_age` now has a default value (`1440`) and does not need to be supplied unless you wish to use a different value
  - `.negative_cache.enabled` now has a default value (`true`) and does not need to be supplied unless you wish to use a different value
  - `.negative_cache.time_to_live` now has a default value (`1440`) and does not need to be supplied unless you wish to use a different value
