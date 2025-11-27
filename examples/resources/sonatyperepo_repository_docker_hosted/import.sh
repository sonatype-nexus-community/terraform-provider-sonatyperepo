# Existing hosted docker repository configuration can be imported as follows.
#
# NOTE: The Identifier REPOSITORY_NAME needs to match repository name in your sonatype nexus repository instance.

#
# WARNING: `storage.latest_policy` will be set to `false` during import as it cannot be read from the API.
#

# Example
terraform import sonatyperepo_repository_docker_hosted.docker REPOSITORY_NAME