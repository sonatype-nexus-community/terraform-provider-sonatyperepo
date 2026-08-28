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
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	configBlockProxyDefaultAlpine    string = "alpine = { key_pair = \"something\" }"
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
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:            TEST_DATA_ALPINE_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_ALPINE,
		SchemaFunc:           repositoryProxyResourceConfig,
		FormatSpecificConfig: configBlockProxyDefaultAlpine,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.93.0 or later
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 0,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 92,
					Patch: 99,
				})
			}
		},
		// Import is broken for Alpine Proxy as alpineSigning is never returned by API
		// See: https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/290
		TestImport: false,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RemoteUrl:  TEST_DATA_ANSIBLE_GALAXY_REMOTE_URL,
		RepoFormat: common.REPO_FORMAT_ANSIBLE_GALAXY,
		SchemaFunc: repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// Only works on NXRM 3.93.0 or later
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 0,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 92,
					Patch: 99,
				})
			}
		},
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_APT_DISTRIBUTION, "bionic"),
			}
		},
		FormatSpecificConfig: configBlockProxyDefaultApt,
		RemoteUrl:            TEST_DATA_APT_PROXY_REMOTE_URL,
		RepoFormat:           common.REPO_FORMAT_APT,
		SchemaFunc:           repositoryProxyResourceConfig,
		TestPreCheck: func(t *testing.T) func() {
			return func() {
				// NXRM 3.82.0 has a bug where proxy repositories may end up in STOPPED state
				// preventing updates - skip tests for this version
				testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
					Major: 3,
					Minor: 82,
					Patch: 0,
				}, &common.SystemVersion{
					Major: 3,
					Minor: 82,
					Patch: 99,
				})
			}
		},
		TestImport: true,
	},
	// NEXUS-48088 prevented this working prior to NXRM 3.88.0 (cargo.requireAuthentication was always returned as false)
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CARGO_REQUIRE_AUTHENTICATION, "false"),
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
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CARGO_REQUIRE_AUTHENTICATION, "true"),
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
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CONAN_PROXY_CONAN_VERSION, "V2"),
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
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_DOCKER_V1_ENABLED, "false"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_DOCKER_PROXY_CACHE_FOREIGN_LAYERS, "false"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_DOCKER_PROXY_INDEX_TYPE, "REGISTRY"),
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
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_MAVEN_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_MAVEN_LAYOUT_POLICY, common.MAVEN_LAYOUT_PERMISSIVE),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_MAVEN_VERSION_POLICY, common.MAVEN_VERSION_POLICY_RELEASE),
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
	// Live repository_firewall coverage (against a real, connected Sonatype IQ Server) is in
	// TestAccRepositoryGenericProxyFirewallToggle below, not this table - see
	// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285.
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NUGET_PROXY_NUGET_VERSION, common.NUGET_PROTOCOL_V3),
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
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(resourceName, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						// resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_CONTENT_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_CONTENT_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_METADATA_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_METADATA_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NEGATIVE_CACHE_ENABLED, fmt.Sprintf("%t", common.DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, fmt.Sprintf("%d", common.DEFAULT_PROXY_NEGATIVE_CACHE_TTL)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_AUTO_BLOCK)),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE)),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_ROUTING_RULE_NAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
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
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(resourceName, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						// resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CLEANUP_POLICY_COUNT, "1"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_CONTENT_MAX_AGE, TEST_DATA_TIMEOUT),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_METADATA_MAX_AGE, TEST_DATA_TIMEOUT),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NEGATIVE_CACHE_ENABLED, "false"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, TEST_DATA_TIMEOUT),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, "false"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, common.HTTP_AUTH_TYPE_USERNAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PREMPTIVE, "false"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_USERNAME, "user"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_PASSWORD, "pass"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, "true"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, "true"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, "2"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, "59"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, "true"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX, "custom-suffix"),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_ROUTING_RULE_NAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
						// resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPLICATION_ASSET_PATH_REGEX, ".*"),
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
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_ONLINE, "true"),
						resource.TestCheckResourceAttrSet(resourceName, repotest.RES_ATTR_URL),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
						// resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_REMOTE_URL, td.RemoteUrl),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_CONTENT_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_CONTENT_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_PROXY_METADATA_MAX_AGE, fmt.Sprintf("%d", common.DEFAULT_PROXY_METADATA_MAX_AGE)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NEGATIVE_CACHE_ENABLED, fmt.Sprintf("%t", common.DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NEGATIVE_CACHE_TIME_TO_LIVE, fmt.Sprintf("%d", common.DEFAULT_PROXY_NEGATIVE_CACHE_TTL)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_BLOCKED, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_BLOCKED)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTO_BLOCK, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_AUTO_BLOCK)),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_CIRCULAR_REDIRECTS)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_ENABLE_COOKIES)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_RETRIES, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_TIMEOUT, fmt.Sprintf("%d", common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE, fmt.Sprintf("%t", common.DEFAULT_HTTP_CLIENT_CONNECTION_USE_TRUST_STORE)),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_HTTP_CLIENT_CONNECTION_USER_AGENT_SUFFIX),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_ROUTING_RULE_NAME),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPLICATION_PRE_EMPTIVE_PULL_ENABLED, "false"),
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
		// Skip formats not supported on older NXRM versions
		if repoFormat == common.REPO_FORMAT_ALPINE {
			// Alpine proxy repositories were added in NXRM 3.93.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 92,
				Patch: 99,
			})
		}
		if repoFormat == common.REPO_FORMAT_ANSIBLE_GALAXY {
			// Ansible Galaxy hosted repositories were added in NXRM 3.93.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 92,
				Patch: 99,
			})
		}

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
	case common.REPO_FORMAT_ALPINE:
		return configBlockProxyDefaultAlpine
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

func repositoryProxyResourceFullConfig(resourceType, repoName, remoteUrl, formatSpecificConfig, formatType string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_cleanup_policy" "my_policy" {
  name = "cleanup-policy-%s"
  format = "%s"
  criteria = {
    asset_regex = "something-here"
  }
}

resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  cleanup = {
	policy_names = [ "cleanup-policy-%s" ]
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
  depends_on = [ sonatyperepo_cleanup_policy.my_policy ]
 }
`, formatType, formatType, resourceType, repoName, formatType, remoteUrl, formatSpecificConfig)
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

func repositoryProxyResourceConfig(resourceType, repoName, repoFormat, remoteUrl, randomString, formatSpecificConfig string, completeData bool) string {
	if completeData {
		return repositoryProxyResourceFullConfig(
			resourceType, repoName, remoteUrl, formatSpecificConfig, strings.ToLower(repoFormat),
		)
	} else {
		return repositoryProxyResourceMinimalConfigWithDefaults(
			resourceType, repoName, remoteUrl, formatSpecificConfig,
		)
	}
}

// ------------------------------------------------------------
// Repository Firewall (Sonatype IQ Server) - live acceptance tests
// ------------------------------------------------------------
// See https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285.
// These tests require a real, licensed Sonatype IQ Server connected to the NXRM instance under
// test - set TF_ACC_IQ_SERVER=1 to enable (see testutil.SkipIfNoIqServer). NXRM rejects
// repository_firewall.enabled = true with "Unable to configure Repository Firewall as not
// connected to Sonatype IQ" unless such a connection exists; internal/provider/provider_test.go's
// TestMain bootstraps and verifies one when TF_ACC_IQ_SERVER=1 is set.

// firewallProxyTestData describes one repository format's live repository_firewall coverage.
type firewallProxyTestData struct {
	RepoFormat           string
	RemoteUrl            string
	FormatSpecificConfig string
	SupportsPccs         bool
}

// firewallProxyTestData enumerates every inline-firewall-capable proxy format (NXRM 3.94+).
// Alpine and Pub are deliberately excluded: NXRM's own server-side validation rejects them
// with "Firewall does not support repository format '<x>'." - confirmed against a real,
// connected IQ Server, not merely a schema-level limitation in this provider. Composer is
// also excluded: repository_firewall applies but doesn't survive a refresh - see
// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/471.
var firewallProxyTestDataTable = []firewallProxyTestData{
	{RepoFormat: common.REPO_FORMAT_CARGO, RemoteUrl: TEST_DATA_CARGO_PROXY_REMOTE_URL, FormatSpecificConfig: configBlockProxyDefaultCargo},
	{RepoFormat: common.REPO_FORMAT_COCOAPODS, RemoteUrl: TEST_DATA_COCOAPODS_PROXY_REMOTE_URL},
	{RepoFormat: common.REPO_FORMAT_CONAN, RemoteUrl: TEST_DATA_CONAN_PROXY_REMOTE_URL, FormatSpecificConfig: configBlockProxyDefaultConan},
	{RepoFormat: common.REPO_FORMAT_CONDA, RemoteUrl: TEST_DATA_CONDA_PROXY_REMOTE_URL},
	{RepoFormat: common.REPO_FORMAT_DOCKER, RemoteUrl: TEST_DATA_DOCKER_PROXY_REMOTE_URL, FormatSpecificConfig: configBlockProxyDefaultDocker},
	{RepoFormat: common.REPO_FORMAT_GO, RemoteUrl: TEST_DATA_GO_PROXY_REMOTE_URL},
	{RepoFormat: common.REPO_FORMAT_HUGGING_FACE, RemoteUrl: TEST_DATA_HUGGING_FACE_PROXY_REMOTE_URL},
	{RepoFormat: common.REPO_FORMAT_MAVEN, RemoteUrl: TEST_DATA_MAVEN_PROXY_REMOTE_URL, FormatSpecificConfig: configBlockProxyDefaultMaven},
	{RepoFormat: common.REPO_FORMAT_NPM, RemoteUrl: TEST_DATA_NPM_PROXY_REMOTE_URL, SupportsPccs: true},
	{RepoFormat: common.REPO_FORMAT_NUGET, RemoteUrl: TEST_DATA_NUGET_PROXY_REMOTE_URL, FormatSpecificConfig: configBlockProxyDefaultNuget},
	{RepoFormat: common.REPO_FORMAT_PYPI, RemoteUrl: TEST_DATA_PYPI_PROXY_REMOTE_URL, SupportsPccs: true},
	{RepoFormat: common.REPO_FORMAT_R, RemoteUrl: TEST_DATA_R_PROXY_REMOTE_URL},
	{RepoFormat: common.REPO_FORMAT_RAW, RemoteUrl: TEST_DATA_RAW_PROXY_REMOTE_URL, FormatSpecificConfig: configBlockProxyDefaultRaw},
	{RepoFormat: common.REPO_FORMAT_RUBY_GEMS, RemoteUrl: TEST_DATA_RUBY_GEMS_PROXY_REMOTE_URL},
	{RepoFormat: common.REPO_FORMAT_YUM, RemoteUrl: TEST_DATA_YUM_PROXY_REMOTE_URL},
}

// repositoryFirewallBlockHcl renders the repository_firewall HCL block, or "" when disabled -
// omitting the block entirely (not just setting enabled = false) reproduces the exact #469/#412
// scenario, where NXRM's response also omits the field entirely once firewall is disabled.
func repositoryFirewallBlockHcl(enabled, quarantine, pccsEnabled bool) string {
	if !enabled {
		return ""
	}
	if pccsEnabled {
		return fmt.Sprintf(`repository_firewall = {
    enabled = true
    quarantine = %t
    pccs_enabled = true
  }`, quarantine)
	}
	return fmt.Sprintf(`repository_firewall = {
    enabled = true
    quarantine = %t
  }`, quarantine)
}

// repositoryFirewallBlockHclExplicit renders an explicit repository_firewall block even when
// disabled (enabled = false), unlike repositoryFirewallBlockHcl above which omits the block
// entirely in that case. This is what reproduces the LITERAL GH-469 report: the block was
// present in the practitioner's HCL with enabled = false, not omitted - repository_firewall is
// Optional-but-not-Computed, so Terraform's plan is a non-null object in that case, whereas
// omitting the block plans null. Those are different scenarios for the provider's state
// reconciliation even though both disable the firewall server-side.
func repositoryFirewallBlockHclExplicit(enabled, quarantine, pccsEnabled bool) string {
	if pccsEnabled {
		return fmt.Sprintf(`repository_firewall = {
    enabled = %t
    quarantine = %t
    pccs_enabled = true
  }`, enabled, quarantine)
	}
	return fmt.Sprintf(`repository_firewall = {
    enabled = %t
    quarantine = %t
  }`, enabled, quarantine)
}

// repositoryProxyResourceConfigWithFirewall builds a minimal proxy repository config with the
// given (possibly empty) repository_firewall block appended.
func repositoryProxyResourceConfigWithFirewall(resourceType, repoName, remoteUrl, formatSpecificConfig, firewallBlock string) string {
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
  %s
 }
`, resourceType, repoName, remoteUrl, formatSpecificConfig, firewallBlock)
}

// TestAccRepositoryGenericProxyFirewallToggle exercises repository_firewall against a real,
// connected Sonatype IQ Server across every inline-firewall-capable proxy format, covering both
// ways a practitioner disables it after it was enabled:
//   - explicitly setting enabled = false while keeping the repository_firewall block in HCL -
//     the LITERAL scenario reported in #469: repository_firewall is Optional-but-not-Computed,
//     so Terraform's plan is a non-null object in this case, and the provider previously nulled
//     the whole attribute out regardless, producing "inconsistent result after apply".
//   - removing the repository_firewall block from HCL entirely - the scenario #412 and this
//     test's step 3 exercise. NXRM omits the field entirely from its response either way once
//     disabled, but the two HCL shapes are different inputs to the provider's state
//     reconciliation and both must converge cleanly.
func TestAccRepositoryGenericProxyFirewallToggle(t *testing.T) {
	for _, td := range firewallProxyTestDataTable {
		t.Run(td.RepoFormat, func(t *testing.T) {
			randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
			resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(td.RepoFormat))
			resourceName := fmt.Sprintf(repoNameFString, resourceType)
			repoName := strings.ToLower(fmt.Sprintf(proxyNameFString, td.RepoFormat, randomString))

			// capability_id is a legacy field from the pre-3.94 Capability-based firewall API -
			// the inline firewall path (exercised here) always leaves it null, by design (see
			// internal/provider/repository/format/proxy_inline_firewall_test.go).
			quarantineChecks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_ENABLED, "true"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE, "true"),
				resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_CAPABILITY_ID),
			}

			steps := []resource.TestStep{
				// 1. No firewall configuration
				{
					Config: repositoryProxyResourceConfigWithFirewall(resourceType, repoName, td.RemoteUrl, td.FormatSpecificConfig, repositoryFirewallBlockHcl(false, false, false)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
						resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL),
					),
				},
				// 2. Firewall enabled with quarantine
				{
					Config: repositoryProxyResourceConfigWithFirewall(resourceType, repoName, td.RemoteUrl, td.FormatSpecificConfig, repositoryFirewallBlockHcl(true, true, false)),
					Check:  resource.ComposeAggregateTestCheckFunc(quarantineChecks...),
				},
			}

			// PCCS is its own FirewallMode (common.FirewallModePccs), mutually exclusive with
			// Quarantine - not a modifier on top of it - so npm/pypi get a dedicated step rather
			// than combining pccs_enabled=true with quarantine=true above.
			if td.SupportsPccs {
				steps = append(steps, resource.TestStep{
					Config: repositoryProxyResourceConfigWithFirewall(resourceType, repoName, td.RemoteUrl, td.FormatSpecificConfig, repositoryFirewallBlockHcl(true, false, true)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_ENABLED, "true"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE, "false"),
						resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_PCCS_ENABLED, "true"),
					),
				})
			}

			// 3. Explicitly disabled - repository_firewall block still present in HCL, with
			// enabled = false. This is the LITERAL #469 regression scenario: the plan is a
			// non-null object (the block is configured), so the provider must keep it in state
			// rather than nulling the whole attribute the way it previously did just because
			// NXRM's response omits `firewall` once disabled.
			steps = append(steps, resource.TestStep{
				Config: repositoryProxyResourceConfigWithFirewall(resourceType, repoName, td.RemoteUrl, td.FormatSpecificConfig, repositoryFirewallBlockHclExplicit(false, false, false)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_ENABLED, "false"),
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE, "false"),
					resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_CAPABILITY_ID),
				),
			})

			// 4. Back to no firewall configuration at all (block removed from HCL, not just
			// enabled = false) - the #412 scenario, and the opposite-direction transition
			// (configured -> unconfigured) that this fix's Update()-path reconciliation must
			// also get right.
			steps = append(steps, resource.TestStep{
				Config: repositoryProxyResourceConfigWithFirewall(resourceType, repoName, td.RemoteUrl, td.FormatSpecificConfig, repositoryFirewallBlockHcl(false, false, false)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL),
				),
			})

			// Import verification for the two formats named in the original bugs' PCCS-capable
			// class (npm, pypi) - the two formats where this coverage matters most.
			if td.SupportsPccs {
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
					testutil.SkipIfNoIqServer(t)
					// repository_firewall requires NXRM 3.94+ (inline firewall)
					testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
						Major: 3,
						Minor: 0,
						Patch: 0,
					}, &common.SystemVersion{
						Major: 3,
						Minor: 93,
						Patch: 99,
					})
				},
				Steps: steps,
			})
		})
	}
}
