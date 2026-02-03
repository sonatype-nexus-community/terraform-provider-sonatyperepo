/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository_test

import (
	"fmt"
	"strings"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	repoNameFString string = "%s.repo"

	hostedNameFString         string = "test-%s-hosted-repo-%s"
	proxyNameFString          string = "test-%s-proxy-repo-%s"
	resourceTypeHostedFString string = "sonatyperepo_repository_%s_hosted"
	resourceTypeProxyFString  string = "sonatyperepo_repository_%s_proxy"

	RES_ATTR_NAME                                             string = "name"
	RES_ATTR_ONLINE                                           string = "online"
	RES_ATTR_CLEANUP                                          string = "cleanup"
	RES_ATTR_CLEANUP_POLICY_COUNT                             string = "cleanup.policy_names.#"
	RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS                 string = "component.proprietary_components"
	RES_ATTR_DOCKER_FORCE_BASIC_AUTH                          string = "docker.force_basic_auth"
	RES_ATTR_DOCKER_PATH_ENABLED                              string = "docker.path_enabled"
	RES_ATTR_DOCKER_V1_ENABLED                                string = "docker.v1_enabled"
	RES_ATTR_DOCKER_PROXY_CACHE_FOREIGN_LAYERS                string = "docker_proxy.cache_foreign_layers"
	RES_ATTR_DOCKER_PROXY_INDEX_TYPE                          string = "docker_proxy.index_type"
	RES_ATTR_RAW_CONTENT_DISPOSITION                          string = "raw.content_disposition"
	RES_ATTR_STORAGE_BLOB_STORE_NAME                          string = "storage.blob_store_name"
	RES_ATTR_STORAGE_LATEST_POLICY                            string = "storage.latest_policy"
	RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION           string = "storage.strict_content_type_validation"
	RES_ATTR_STORAGE_WRITE_POLICY                             string = "storage.write_policy"
	RES_ATTR_URL                                              string = "url"
	RES_ATTR_PROXY_REMOTE_URL                                 string = "proxy.remote_url"
	RES_ATTR_PROXY_CONTENT_MAX_AGE                            string = "proxy.content_max_age"
	RES_ATTR_PROXY_METADATA_MAX_AGE                           string = "proxy.metadata_max_age"
	RES_ATTR_NEGATIVE_CACHE_ENABLED                           string = "negative_cache.enabled"
	RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE                      string = "negative_cache.time_to_live"
	RES_ATTR_HTTP_CLIENT_BLOCKED                              string = "http_client.blocked"
	RES_ATTR_HTTP_CLIENT_AUTO_BLOCK                           string = "http_client.auto_block"
	RES_ATTR_HTTP_CLIENT_AUTHENTICATION                       string = "http_client.authentication"
	RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS string = "http_client.connection.enable_circular_redirects"
	RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES            string = "http_client.connection.enable_cookies"
	RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE           string = "http_client.connection.use_trust_store"
	RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES                   string = "http_client.connection.retries"
	RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT                   string = "http_client.connection.timeout"
	RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX         string = "http_client.connection.user_agent_suffix"
	RES_ATTR_GROUP_MEMBER_NAMES                               string = "group.member_names.#"
	RES_ATTR_REPLICATION                                      string = "replication"
	RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED             string = "replication.preemptive_pull_enabled"
	RES_ATTR_REPLICATION_ASSET_PATH_REGEX                     string = "replication.asset_path_regex"
	RES_ATTR_ROUTING_RULE_NAME                                string = "routing_rule"
	RES_ATTR_REPOSITORY_FIREWALL_ENABLED                      string = "repository_firewall.enabled"
	RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE                   string = "repository_firewall.quarantine"
	RES_ATTR_APT_DISTRIBUTION                                 string = "apt.distribution"
	RES_ATTR_CARGO_REQUIRE_AUTHENTICATION                     string = "cargo.require_authentication"
	RES_ATTR_CONAN_PROXY_CONAN_VERSION                        string = "conan.conan_version"
	RES_ATTR_MAVEN_CONTENT_DISPOSITION                        string = "maven.content_disposition"
	RES_ATTR_MAVEN_LAYOUT_POLICY                              string = "maven.layout_policy"
	RES_ATTR_MAVEN_VERSION_POLICY                             string = "maven.version_policy"
	RES_ATTR_NUGET_PROXY_NUGET_VERSION                        string = "nuget_proxy.nuget_version"

	TEST_DATA_APT_PROXY_REMOTE_URL          string = "https://archive.ubuntu.com/ubuntu/"
	TEST_DATA_CARGO_PROXY_REMOTE_URL        string = "https://index.crates.io/"
	TEST_DATA_COCOAPODS_PROXY_REMOTE_URL    string = "https://cdn.cocoapods.org/"
	TEST_DATA_COMPOSER_PROXY_REMOTE_URL     string = "https://packagist.org/"
	TEST_DATA_CONAN_PROXY_REMOTE_URL        string = "https://center2.conan.io"
	TEST_DATA_CONDA_PROXY_REMOTE_URL        string = "https://repo.anaconda.com/pkgs/"
	TEST_DATA_DOCKER_PROXY_REMOTE_URL       string = "https://registry-1.docker.io"
	TEST_DATA_GO_PROXY_REMOTE_URL           string = "https://proxy.golang.org"
	TEST_DATA_HELM_PROXY_REMOTE_URL         string = "https://charts.helm.sh/stable"
	TEST_DATA_HUGGING_FACE_PROXY_REMOTE_URL string = "https://huggingface.co"
	TEST_DATA_MAVEN_PROXY_REMOTE_URL        string = "https://repo.maven.apache.org/maven2"
	TEST_DATA_NPM_PROXY_REMOTE_URL          string = "https://registry.npmjs.org"
	TEST_DATA_NUGET_PROXY_REMOTE_URL        string = "https://api.nuget.org/v3/index.json"
	TEST_DATA_P2_PROXY_REMOTE_URL           string = "https://download.eclipse.org/releases/2025-06"
	TEST_DATA_TIMEOUT                       string = "1439"
)

var (
	errorMessageBlobStoreNotFound                = "Blob store.*not found"
	errorMessageGroupMemberNamesEmpty            = "Attribute group.member_names list must contain at least 1 elements"
	errorMessageInvalidRemoteUrl                 = "Attribute proxy.remote_url must be a valid HTTP URL"
	errorMessageHttpClientConnectionRetriesValue = fmt.Sprintf(
		"Attribute http_client.connection.retries value must be between %d and %d",
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MIN,
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MAX,
	)
	errorMessageHttpClientConnectionTimeoutValue = fmt.Sprintf(
		"Attribute http_client.connection.timeout value must be between %d and %d",
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MIN,
		common.REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MAX,
	)
	errorMessageNegativeCacheTimeoutValue = "Attribute negative_cache.time_to_live value must be at least 0"
	errorMessageStorageRequired           = "The argument \"storage\" is required, but no definition was found."
)

type repositoryHostedTestData struct {
	CheckFunc  func(resourceName string) []resource.TestCheckFunc
	RepoFormat string
	SchemaFunc func(resourceType, repoName, repoFormat, randomString string, completeData bool) string
	TestImport bool
}

type repositoryProxyTestData struct {
	CheckFunc  func(resourceName string) []resource.TestCheckFunc
	RemoteUrl  string
	RepoFormat string
	SchemaFunc func(resourceType, repoName, repoFormat, remoteUrl, randomString string, completeData bool) string
}

// func allRepositoryFormatsHostedGeneric() []string {
// 	return []string{
// 		// strings.ToLower(common.REPO_FORMAT_APT),		// Requires specific apt configuration
// 		strings.ToLower(common.REPO_FORMAT_CARGO),
// 		strings.ToLower(common.REPO_FORMAT_CONAN),
// 		// strings.ToLower(common.REPO_FORMAT_DOCKER),	// Requires specific docker configuration
// 		strings.ToLower(common.REPO_FORMAT_GIT_LFS),
// 		strings.ToLower(common.REPO_FORMAT_HELM),
// 		// strings.ToLower(common.REPO_FORMAT_MAVEN),	// Requires specific maven configuration
// 		strings.ToLower(common.REPO_FORMAT_NPM),
// 		strings.ToLower(common.REPO_FORMAT_NUGET),
// 		strings.ToLower(common.REPO_FORMAT_PYPI),
// 		strings.ToLower(common.REPO_FORMAT_R),
// 		// strings.ToLower(common.REPO_FORMAT_RAW),		// Requires specific raw configuration
// 		strings.ToLower(common.REPO_FORMAT_RUBY_GEMS),
// 		// strings.ToLower(common.REPO_FORMAT_YUM),		// Requires specific yum configuration
// 	}
// }

// func allRepositoryFormatsProxyGeneric() []string {
// 	return []string{
// 		// strings.ToLower(common.REPO_FORMAT_APT),		// Requires specific apt configuration
// 		strings.ToLower(common.REPO_FORMAT_CARGO),
// 		strings.ToLower(common.REPO_FORMAT_CONAN),
// 		// strings.ToLower(common.REPO_FORMAT_DOCKER),	// Requires specific docker configuration
// 		strings.ToLower(common.REPO_FORMAT_GIT_LFS),
// 		strings.ToLower(common.REPO_FORMAT_HELM),
// 		// strings.ToLower(common.REPO_FORMAT_MAVEN),	// Requires specific maven configuration
// 		strings.ToLower(common.REPO_FORMAT_NPM),
// 		strings.ToLower(common.REPO_FORMAT_NUGET),
// 		strings.ToLower(common.REPO_FORMAT_PYPI),
// 		strings.ToLower(common.REPO_FORMAT_R),
// 		// strings.ToLower(common.REPO_FORMAT_RAW),		// Requires specific raw configuration
// 		strings.ToLower(common.REPO_FORMAT_RUBY_GEMS),
// 		// strings.ToLower(common.REPO_FORMAT_YUM),		// Requires specific yum configuration
// 	}
// }

func allRepositoryFormatsGroupGeneric() []string {
	return []string{
		strings.ToLower(common.REPO_FORMAT_CONAN),
		// strings.ToLower(common.REPO_FORMAT_GO),		// No hosted repo
		// strings.ToLower(common.REPO_FORMAT_MAVEN),	// Hosted repo requires specific config
		strings.ToLower(common.REPO_FORMAT_NPM),
		strings.ToLower(common.REPO_FORMAT_NUGET),
		strings.ToLower(common.REPO_FORMAT_PYPI),
		strings.ToLower(common.REPO_FORMAT_R),
		strings.ToLower(common.REPO_FORMAT_RUBY_GEMS),
		// strings.ToLower(common.REPO_FORMAT_YUM),		// Hosted repo requires specific config
	}
}

// ------------------------------------------------------------
// HOSTED REPO TESTING (GENERIC)
// ------------------------------------------------------------
// func TestAccRepositoryGenericHostedResources(t *testing.T) {
// 	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

// 	for _, repoFormat := range allRepositoryFormatsHostedGeneric() {
// 		resourceType := fmt.Sprintf("sonatyperepo_repository_%s_hosted", repoFormat)
// 		reosurceName := fmt.Sprintf(repoNameFString, resourceType)
// 		repoName := fmt.Sprintf("%s-hosted-repo-%s", repoFormat, randomString)

// 		resource.Test(t, resource.TestCase{
// 			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
// 			Steps: []resource.TestStep{
// 				// Create and Read testing
// 				{
// 					Config: repositoryHostedResourceNoCleanupNoProprietaryConfig(resourceType, repoName, common.WRITE_POLICY_ALLOW_ONCE, true),
// 					Check: resource.ComposeAggregateTestCheckFunc(
// 						// Verify Hosted
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_NAME, repoName),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_ONLINE, "true"),
// 						resource.TestCheckResourceAttrSet(reosurceName, RES_ATTR_URL),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
// 						resource.TestCheckNoResourceAttr(reosurceName, RES_ATTR_CLEANUP),
// 					),
// 				},
// 				// Update 1 - Set Offline
// 				{
// 					Config: repositoryHostedResourceNoCleanupWithProprietaryConfig(resourceType, repoName, common.WRITE_POLICY_ALLOW_ONCE, false, false),
// 					Check: resource.ComposeAggregateTestCheckFunc(
// 						// Verify Hosted
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_NAME, repoName),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_ONLINE, "false"),
// 						resource.TestCheckResourceAttrSet(reosurceName, RES_ATTR_URL),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
// 						resource.TestCheckNoResourceAttr(reosurceName, RES_ATTR_CLEANUP),
// 					),
// 				},
// 				// Update 2 - Set Proprietary Components
// 				{
// 					Config: repositoryHostedResourceNoCleanupWithProprietaryConfig(resourceType, repoName, common.WRITE_POLICY_ALLOW_ONCE, true, true),
// 					Check: resource.ComposeAggregateTestCheckFunc(
// 						// Verify Hosted
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_NAME, repoName),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_ONLINE, "true"),
// 						resource.TestCheckResourceAttrSet(reosurceName, RES_ATTR_URL),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
// 						resource.TestCheckResourceAttr(reosurceName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "true"),
// 						resource.TestCheckNoResourceAttr(reosurceName, RES_ATTR_CLEANUP),
// 					),
// 				},
// 			},
// 		})
// 	}
// }

func repositoryHostedResourceNoCleanupNoProprietaryConfig(resourceType, repoName, writePolicy string, repoOnline bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
name = "%s"
online = %t
storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "%s"
}
}`, resourceType, repoName, repoOnline, writePolicy)
}

func repositoryHostedResourceNoCleanupWithProprietaryConfig(resourceType, repoName, writePolicy string, proprietaryComponents, repoOnline bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
name = "%s"
online = %t
storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "%s"
}
component = {
	proprietary_components = %t
}
}`, resourceType, repoName, repoOnline, writePolicy, proprietaryComponents)
}

// ------------------------------------------------------------
// GROUP REPO TESTING (GENERIC)
// ------------------------------------------------------------
func TestAccRepositoryGenericGroupResources(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	for _, repoFormat := range allRepositoryFormatsGroupGeneric() {
		resourceTypeGroup := fmt.Sprintf("sonatyperepo_repository_%s_group", repoFormat)
		resourceTypeHosted := fmt.Sprintf("sonatyperepo_repository_%s_hosted", repoFormat)
		reosurceNameGroup := fmt.Sprintf(repoNameFString, resourceTypeGroup)
		reosurceNameHosted := fmt.Sprintf(repoNameFString, resourceTypeHosted)
		repoNameGroup := fmt.Sprintf("%s-group-repo-%s", repoFormat, randomString)
		repoNameHosted := fmt.Sprintf("%s-hosted-repo-%s", repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: repositoryGroupResourceNoCleanupNoProprietaryConfig(
						resourceTypeGroup,
						resourceTypeHosted,
						repoNameGroup,
						repoNameHosted,
						common.WRITE_POLICY_ALLOW_ONCE,
						false,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify Hosted
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_NAME, repoNameHosted),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_ONLINE, "false"),
						resource.TestCheckResourceAttrSet(reosurceNameHosted, RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
						resource.TestCheckNoResourceAttr(reosurceNameHosted, RES_ATTR_CLEANUP),

						// Verify Group
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_NAME, repoNameGroup),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_ONLINE, "false"),
						resource.TestCheckResourceAttrSet(reosurceNameGroup, RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					),
				},
				// Update 1 - Set Online
				{
					Config: repositoryGroupResourceNoCleanupNoProprietaryConfig(
						resourceTypeGroup,
						resourceTypeHosted,
						repoNameGroup,
						repoNameHosted,
						common.WRITE_POLICY_ALLOW_ONCE,
						true,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify Hosted
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_NAME, repoNameHosted),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(reosurceNameHosted, RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
						resource.TestCheckResourceAttr(reosurceNameHosted, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
						resource.TestCheckNoResourceAttr(reosurceNameHosted, RES_ATTR_CLEANUP),

						// Verify Group
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_NAME, repoNameGroup),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(reosurceNameGroup, RES_ATTR_URL),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						resource.TestCheckResourceAttr(reosurceNameGroup, RES_ATTR_GROUP_MEMBER_NAMES, "1"),
					),
				},
			},
		})
	}
}

func repositoryGroupResourceNoCleanupNoProprietaryConfig(resourceTypeGroup, resourceTypeHosted, repoNameGroup, repoNameHosted, writePolicy string, repoOnline bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
	name = "%s"
	online = %t
	storage = {
		blob_store_name = "default"
		strict_content_type_validation = true
		write_policy = "%s"
	}
}

resource "%s" "repo" {
  name = "%s"
  online = %t
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["%s"]
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeHosted, repoNameHosted, repoOnline, writePolicy, resourceTypeGroup, repoNameGroup, repoOnline, repoNameHosted, resourceTypeHosted)
}

// ------------------------------------------------------------
// PROXY REPO TESTING (GENERIC)
// ------------------------------------------------------------
