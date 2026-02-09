# Existing proxy Conan repository configuration can be imported as follows.
#
# NOTE: The Identifier REPOSITORY_NAME needs to match repository name in your Sonatype Nexus Repository.
#
# NOTE: This does not work when running against Sonatype Nexus Repository version 3.85 or earlier.

# Example
terraform import sonatyperepo_repository_conan_proxy.example REPOSITORY_NAME
