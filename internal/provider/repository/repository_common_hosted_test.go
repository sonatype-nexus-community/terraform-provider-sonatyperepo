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
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ------------------------------------------------------------
// Test Data Scenarios
// ------------------------------------------------------------
var hostedTestData = []repositoryHostedTestData{
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_APT_DISTRIBUTION, "bionic"),
			}
		},
		RepoFormat: common.REPO_FORMAT_APT,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: false, // TODO: Document this does not work based on NXRM 3.88
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_CARGO,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_CONAN,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_V1_ENABLED, "false"),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_PATH_ENABLED, "false"),
			}
		},
		RepoFormat: common.REPO_FORMAT_DOCKER,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_GIT_LFS,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_HELM,
		SchemaFunc: repositoryHostedResourceConfig,
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
		RepoFormat: common.REPO_FORMAT_MAVEN,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_NPM,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_NUGET,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_PYPI,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_R,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_RAW,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_RUBY_GEMS,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_YUM,
		SchemaFunc: repositoryHostedResourceConfig,
		TestImport: true,
	},
}

// ------------------------------------------------------------
// HOSTED REPO TESTING (GENERIC)
// ------------------------------------------------------------
func TestAccRepositoryGenericHostedByFormat(t *testing.T) {
	for _, td := range hostedTestData {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeHostedFString, strings.ToLower(td.RepoFormat))
		resourceName := fmt.Sprintf(repoNameFString, resourceType)
		repoName := fmt.Sprintf(hostedNameFString, td.RepoFormat, randomString)

		var steps []resource.TestStep
		// 1. Create with minimal configuration relying on defaults
		steps = append(steps, resource.TestStep{
			Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, randomString, false),
			Check: resource.ComposeAggregateTestCheckFunc(
				append(
					// Test Case Specific Checks
					td.CheckFunc(resourceName),

					// Generic Checks
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceName, RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "false"),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
					resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"),
				)...,
			),
		})
		// 2. Update to use full config
		steps = append(steps, resource.TestStep{
			Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, randomString, true),
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
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
					resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_ROUTING_RULE_NAME),
					resource.TestCheckResourceAttr(resourceName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "true"),
				)...,
			),
		})

		// 3. Import and verify no changes
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
			Steps:                    steps,
		})
	}
}

func TestAccRepositoryGenericHostedInvalidBlobStore(t *testing.T) {
	for _, repoFormat := range common.AllHostedFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeHostedFString, strings.ToLower(repoFormat))
		repoName := fmt.Sprintf(hostedNameFString, repoFormat, randomString)

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
	write_policy = "ALLOW"
  }
  %s
 }
`, resourceType, repoName, formatSpecificHostedDefaultConfig(repoFormat)),
					ExpectError: regexp.MustCompile(errorMessageBlobStoreNotFound),
				},
			},
		})
	}
}

const (
	configBlockHostedDefaultApt      string = "apt = { distribution = \"bionic\" }\napt_signing = { key_pair = \"something\" }"
	configBlockHostedDefaultCargo    string = ""
	configBlockHostedDefaultConan    string = ""
	configBlockHostedDefaultDocker   string = "docker = { force_basic_auth = false\nv1_enabled = false\npath_enabled = false }"
	configBlockHostedDefaultGitLfs   string = ""
	configBlockHostedDefaultHelm     string = ""
	configBlockHostedDefaultMaven    string = "maven = { layout_policy = \"PERMISSIVE\"\nversion_policy = \"RELEASE\" }"
	configBlockHostedDefaultNpm      string = ""
	configBlockHostedDefaultNuget    string = ""
	configBlockHostedDefaultPypi     string = ""
	configBlockHostedDefaultR        string = ""
	configBlockHostedDefaultRaw      string = "raw = { content_disposition = \"ATTACHMENT\" }"
	configBlockHostedDefaultRubyGems string = ""
	configBlockHostedDefaultYum      string = "yum = { repo_data_depth = 1 }"
)

func formatSpecificHostedDefaultConfig(repoFormat string) string {
	switch repoFormat {
	case common.REPO_FORMAT_APT:
		return configBlockHostedDefaultApt
	case common.REPO_FORMAT_CARGO:
		return configBlockHostedDefaultCargo
	case common.REPO_FORMAT_CONAN:
		return configBlockHostedDefaultConan
	case common.REPO_FORMAT_DOCKER:
		return configBlockHostedDefaultDocker
	case common.REPO_FORMAT_GIT_LFS:
		return configBlockHostedDefaultGitLfs
	case common.REPO_FORMAT_HELM:
		return configBlockHostedDefaultHelm
	case common.REPO_FORMAT_MAVEN:
		return configBlockHostedDefaultMaven
	case common.REPO_FORMAT_NPM:
		return configBlockHostedDefaultNpm
	case common.REPO_FORMAT_NUGET:
		return configBlockHostedDefaultNuget
	case common.REPO_FORMAT_PYPI:
		return configBlockHostedDefaultPypi
	case common.REPO_FORMAT_R:
		return configBlockHostedDefaultR
	case common.REPO_FORMAT_RAW:
		return configBlockHostedDefaultRaw
	case common.REPO_FORMAT_RUBY_GEMS:
		return configBlockHostedDefaultRubyGems
	case common.REPO_FORMAT_YUM:
		return configBlockHostedDefaultYum

	default:
		return ""
	}
}

func repositoryHostedResourceFullConfig(resourceType, repoName, formatSpecificConfig string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  component = {
    proprietary_components = true
  }
  %s
 }
`, resourceType, repoName, formatSpecificConfig)
}

func repositoryHostedResourceMinimalConfigWithDefaults(resourceType, repoName, formatSpecificConfig string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = false
	write_policy = "ALLOW"
  }
  %s
 }
`, resourceType, repoName, formatSpecificConfig)
}

func repositoryHostedResourceConfig(resourceType, repoName, repoFormat, randomString string, completeData bool) string {
	configBlock := formatSpecificHostedDefaultConfig(repoFormat)
	if completeData {
		return repositoryHostedResourceFullConfig(
			resourceType, repoName, configBlock,
		)
	} else {
		return repositoryHostedResourceMinimalConfigWithDefaults(
			resourceType, repoName, configBlock,
		)
	}
}
