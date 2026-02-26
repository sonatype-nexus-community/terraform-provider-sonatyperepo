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
	"regexp"
	"strings"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	configBlockProxyDefaultApt       string = "apt = { distribution = \"bionic\" }"
	configBlockProxyDefaultCargo     string = "cargo = { require_authentication = false }"
	configBlockProxyDefaultConan     string = "conan = { conan_version = \"V2\" }"
	configBlockProxyDefaultDocker    string = "docker = { force_basic_auth = false\nv1_enabled = false }\ndocker_proxy = { }"
	configBlockProxyDefaultMaven     string = "maven = { layout_policy = \"PERMISSIVE\"\nversion_policy = \"RELEASE\" }"
	configBlockProxyDefaultNuget     string = "nuget_proxy = { nuget_version = \"V3\" }"
	configBlockProxyDefaultRaw       string = "raw = { content_disposition = \"ATTACHMENT\" }"
	configBlockProxyDefaultSwift     string = "swift = { }"
	configBlockProxyDefaultTerraform string = "terraform = { }"
)

// ------------------------------------------------------------
// Test Data Scenarios
// ------------------------------------------------------------
var proxyTestData = []repositoryProxyTestData{
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_APT_DISTRIBUTION, "bionic"),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultApt,
		RemoteUrl:            TEST_DATA_APT_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_APT,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestImport:           true,
	},
	// NEXUS-48088 prevented this working prior to NXRM 3.88.0 (cargo.requireAuthentication was always returned as false)
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_CARGO_REQUIRE_AUTHENTICATION, "false"),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultCargo,
		RemoteUrl:            TEST_DATA_CARGO_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_CARGO,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 88,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 4,
					Minor: 0,
					Patch: 0,
				})
			}
		},
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_CARGO_REQUIRE_AUTHENTICATION, "true"),
			}
		},
		FormatSpecificConfig: "cargo = { require_authentication = true }",
		RemoteUrl:            TEST_DATA_CARGO_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_CARGO,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.88.0 or later
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 0,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 87,
					Patch: 99,
				})
			}
		},
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_COCOAPODS_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_COCOAPODS,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_COMPOSER_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_COMPOSER,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		// Prior to NXRM 3.86 - conanProxy was not returned by API
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		FormatSpecificConfig: configBlockProxyDefaultConan,
		RemoteUrl:            TEST_DATA_CONAN_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_CONAN,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.86.0 or earlier
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 86,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 4,
					Minor: 0,
					Patch: 0,
				})
			}
		},
		TestImport: false,
	},
	{
		// Required NXRM 3.86 or later to work (NEXUS-49755, NEXUS-47906)
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_CONAN_PROXY_CONAN_VERSION, "V2"),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultConan,
		RemoteUrl:            TEST_DATA_CONAN_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_CONAN,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.86.0 or earlier
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 0,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 85,
					Patch: 99,
				})
			}
		},
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_CONDA_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_CONDA,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_V1_ENABLED, "false"),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_PROXY_CACHE_FOREIGN_LAYERS, "false"),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_PROXY_INDEX_TYPE, "REGISTRY"),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultDocker,
		RemoteUrl:            TEST_DATA_DOCKER_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_DOCKER,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestImport:           true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_GO_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_GO,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_HELM_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_HELM,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_HUGGING_FACE_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_HUGGING_FACE,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_MAVEN_CONTENT_DISPOSITION),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_MAVEN_LAYOUT_POLICY, common.MAVEN_LAYOUT_PERMISSIVE),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_MAVEN_VERSION_POLICY, common.MAVEN_VERSION_POLICY_RELEASE),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultMaven,
		RemoteUrl:            TEST_DATA_MAVEN_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_MAVEN,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestImport:           true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_NPM_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_NPM,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	// Will fail without IQ Server Connection: "Unable to configure Repository Firewall as not connected to Sonatype IQ"
	// {
	// 	CheckFunc: func(resourceName string) []resource.TestCheckFunc {
	// 		return []resource.TestCheckFunc{
	// 			resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPOSITORY_FIREWALL_ENABLED, "true"),
	// 			resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE, "true"),
	// 		}
	// 	},
	// 	RemoteUrl:  TEST_DATA_NPM_PROXY_REMOTE_URL,
	// 	RepoFormat: common.REPO_FORMAT_NPM,
	// 	SchemaFunc: repositoryProxyResourceConfigWithFirewall,
	// },
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_NUGET_PROXY_NUGET_VERSION, common.NUGET_PROTOCOL_V3),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultNuget,
		RemoteUrl:            TEST_DATA_NUGET_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_NUGET,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestImport:           true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_P2_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_P2,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_PYPI_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_PYPI,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_R_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_R,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		FormatSpecificConfig: configBlockProxyDefaultRaw,
		RemoteUrl:            TEST_DATA_RAW_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_RAW,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestImport:           true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_RUBY_GEMS_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_RUBY_GEMS,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		FormatSpecificConfig: configBlockProxyDefaultSwift,
		RemoteUrl:            TEST_DATA_SWIFT_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_SWIFT,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.89.0 or later
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 0,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 88,
					Patch: 99,
				})
			}
		},
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		FormatSpecificConfig: configBlockProxyDefaultTerraform,
		RemoteUrl:            TEST_DATA_TERRAFORM_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_TERRAFORM,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.88.0 or later
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 0,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 87,
					Patch: 99,
				})
			}
		},
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_YUM_PROXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_YUM,
		SchemaFunc: repositoryProxyResourceConfig,
		TestImport: true,
	},
}

// ------------------------------------------------------------
// PROXY REPO TESTING (GENERIC)
// ------------------------------------------------------------
func TestAccRepositoryGenericProxyByFormat(t *testing.T) {
	for _, td := range proxyTestData {
		t.Run(td.RepoFormat, func(t *testing.T) {
			randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
			resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(td.RepoFormat))
			resourceName := fmt.Sprintf(repoNameFString, resourceType)
			repoName := strings.ToLower(fmt.Sprintf(proxyNameFString, td.RepoFormat, randomString))

			var steps []resource.TestStep

			// 1. Create with minimal configuration relying on defaults
			steps = append(steps, resource.TestStep{
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, td.RemoteUrl, randomString, td.FormatSpecificConfig, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						// Test Case Specific Checks
						td.CheckFunc(resourceName),

						// Generic Checks
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(resourceName, RES_ATTR_URL),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						// resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_CONTENT_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_CONTENT_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_METADATA_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_METADATA_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_ENABLED, fmt.Sprintf("%t", common.DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, fmt.Sprintf("%d", common.DEFAULT_PROXY_NEGATIVE_CACHE_TTL)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_AUTO_BLOCK)),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE)),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
					)...,
				),
			})

			// 2. Update to use full config
			steps = append(steps, resource.TestStep{
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, td.RemoteUrl, randomString, td.FormatSpecificConfig, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						// Test Case Specific Checks
						td.CheckFunc(resourceName),

						// Generic Checks
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(resourceName, RES_ATTR_URL),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						// resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_CONTENT_MAX_AGE, TEST_DATA_TIMEOUT),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_METADATA_MAX_AGE, TEST_DATA_TIMEOUT),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_ENABLED, "false"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, TEST_DATA_TIMEOUT),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "false"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, common.HTTP_AUTH_TYPE_USERNAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PREMPTIVE, "false"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME, "user"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD, "pass"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "true"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "2"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "59"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "custom-suffix"),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
						// resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_ASSET_PATH_REGEX, ".*"),
					)...,
				),
			})

			// 3. Revert back to Simple Config
			steps = append(steps, resource.TestStep{
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, td.RemoteUrl, randomString, td.FormatSpecificConfig, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						// Test Case Specific Checks
						td.CheckFunc(resourceName),

						// Generic Checks
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(resourceName, RES_ATTR_URL),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						// resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_CONTENT_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_CONTENT_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_PROXY_METADATA_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_METADATA_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_ENABLED, fmt.Sprintf("%t", common.DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, fmt.Sprintf("%d", common.DEFAULT_PROXY_NEGATIVE_CACHE_TTL)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_AUTO_BLOCK)),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE)),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX),
						resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
					)...,
				),
			})

			// 4. Import and verify no changes
			if td.TestImport {
				steps = append(steps, resource.TestStep{
					ResourceName:                         resourceName,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        repoName,
					ImportStateVerifyIdentifierAttribute: "name",
					ImportStateVerifyIgnore:              []string{"last_updated"},
				})
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
				PreCheck: func() {
					if td.TestPreCheck != nil {
						td.TestPreCheck(t)()
					}
				},
				Steps: steps,
			})
		})
	}
}

func TestAccRepositoryGenericProxyInvalidRemoteUrl(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageInvalidRemoteUrl),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidBlobStore(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := strings.ToLower(fmt.Sprintf(proxyNameFString, repoFormat, randomString))

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://some.source.url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidHttpConnectionRetries(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// HTTP Connection Timeout to large
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      retries = 11
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
				},
				// HTTP Connection Timeout to small
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      retries = -1
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionRetriesValue),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidHttpConnectionTimeout(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// HTTP Connection Timeout to large
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      timeout = 3601
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
				},
				// HTTP Connection Timeout to small
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-remote-url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      timeout = 0
    }
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageHttpClientConnectionTimeoutValue),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyInvalidNegativeCacheTtl(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://some.source.url"
  }
  negative_cache = {
    enabled = true
    time_to_live = -1
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageNegativeCacheTimeoutValue),
				},
			},
		})
	}
}

func TestAccRepositoryGenericProxyMissingStorage(t *testing.T) {
	for _, repoFormat := range common.AllProxyFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(proxyNameFString, repoFormat, randomString)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  proxy = {
    remote_url = "https://some.source.url"
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
    connection = {
      timeout = 0
    }
  }
  %s
 }`, resourceType, repoName, formatSpecificProxyDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageStorageRequired),
				},
			},
		})
	}
}

func formatSpecificProxyDefaultConfig(repoFormat string) string {
	switch repoFormat {
	case common.REPO_FORMAT_APT:
		return configBlockProxyDefaultApt
	case common.REPO_FORMAT_CARGO:
		return configBlockProxyDefaultCargo
	case common.REPO_FORMAT_CONAN:
		return configBlockProxyDefaultConan
	case common.REPO_FORMAT_DOCKER:
		return configBlockProxyDefaultDocker
	case common.REPO_FORMAT_MAVEN:
		return configBlockProxyDefaultMaven
	case common.REPO_FORMAT_NUGET:
		return configBlockProxyDefaultNuget
	case common.REPO_FORMAT_RAW:
		return configBlockProxyDefaultRaw
	case common.REPO_FORMAT_TERRAFORM:
		return configBlockProxyDefaultTerraform

	default:
		return ""
	}
}

func repositoryProxyResourceFullConfig(resourceType, repoName, remoteUrl, formatSpecificConfig string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "%s"
    content_max_age = 1439
    metadata_max_age = 1439
  }
  negative_cache = {
    enabled = false
    time_to_live = 1439
  }
  http_client = {
    blocked = false
    auto_block = false
	authentication = {
	  type = "username"
      preemptive = false
      password    = "pass"
      username    = "user"
	}
    connection = {
	  enable_circular_redirects = true
	  enable_cookies = true
	  retries = 2
	  timeout = 59
	  use_trust_store = true
	  user_agent_suffix = "custom-suffix"
	}
  }
  replication = {
    preemptive_pull_enabled = false
  }
  %s
 }
`, resourceType, repoName, remoteUrl, formatSpecificConfig)
}

func repositoryProxyResourceMinimalConfigWithDefaults(resourceType, repoName, remoteUrl, formatSpecificConfig string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "%s"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  %s
 }
`, resourceType, repoName, remoteUrl, formatSpecificConfig)
}

// See https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285
// func repositoryProxyResourceMinimalConfigWithFirewallEnabledNoPccs(resourceType, repoName, remoteUrl, formatSpecificConfig string) string {
// 	return fmt.Sprintf(utils_test.ProviderConfig+`
// resource "%s" "repo" {
//   name = "%s"
//   online = true
//   storage = {
//     blob_store_name = "default"
//     strict_content_type_validation = true
//   }
//   proxy = {
//     remote_url = "%s"
//     content_max_age = 1440
//     metadata_max_age = 1440
//   }
//   negative_cache = {
//     enabled = true
//     time_to_live = 1440
//   }
//   http_client = {
//     blocked = false
//     auto_block = true
//   }
//   repository_firewall = {
//     enabled = true
//     quarantine = true
//   }
//   %s
//  }
// `, resourceType, repoName, remoteUrl, formatSpecificConfig)
// }

func repositoryProxyResourceConfig(resourceType, repoName, repoFormat, remoteUrl, randomString, formatSpecificConfig string, completeData bool) string {
	if completeData {
		return repositoryProxyResourceFullConfig(
			resourceType, repoName, remoteUrl, formatSpecificConfig,
		)
	} else {
		return repositoryProxyResourceMinimalConfigWithDefaults(
			resourceType, repoName, remoteUrl, formatSpecificConfig,
		)
	}
}

// See https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285
// func repositoryProxyResourceConfigWithFirewall(resourceType, repoName, repoFormat, remoteUrl, randomString string, completeData bool) string {
// 	return repositoryProxyResourceMinimalConfigWithFirewallEnabledNoPccs(
// 		resourceType, repoName, remoteUrl, formatSpecificProxyDefaultConfig(repoFormat),
// 	)
// }
