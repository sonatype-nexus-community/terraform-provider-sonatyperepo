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

// ------------------------------------------------------------
// Test Data Scenarios
// ------------------------------------------------------------
var hostedTestData = []repositoryHostedTestData{
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_ANSIBLE_GALAXY,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: false,
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
		RepoFormat:                    common.REPO_FORMAT_APT,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		// Import is broken for APT Hosted as aptSigning is never returned by API
		// See: https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/290
		TestImport: false,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_CARGO,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_CONAN,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_DOCKER_V1_ENABLED, "false"),
			}
		},
		RepoFormat:                    common.REPO_FORMAT_DOCKER,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_GIT_LFS,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_HELM,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_MAVEN_CONTENT_DISPOSITION, common.CONTENT_DISPOSITION_ATTACHMENT),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_MAVEN_LAYOUT_POLICY, common.MAVEN_LAYOUT_PERMISSIVE),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_MAVEN_VERSION_POLICY, common.MAVEN_VERSION_POLICY_RELEASE),
			}
		},
		RepoFormat:                    common.REPO_FORMAT_MAVEN,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_NPM,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_NUGET,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_PYPI,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_R,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_RAW,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_RUBY_GEMS,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat: common.REPO_FORMAT_TERRAFORM,
		SchemaFunc: repositoryHostedResourceConfig,
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
		SupportsProprietaryComponents: false,
		TestImport:                    true,
	},
	{
		CheckFunc: func(resourceName string) []resource.TestCheckFunc {
			return []resource.TestCheckFunc{}
		},
		RepoFormat:                    common.REPO_FORMAT_YUM,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: true,
		TestImport:                    true,
	},
}

// ------------------------------------------------------------
// HOSTED REPO TESTING (GENERIC)
// ------------------------------------------------------------
func TestAccRepositoryGenericHostedByFormat(t *testing.T) {
	for _, td := range hostedTestData {
		t.Run(td.RepoFormat, func(t *testing.T) {
			randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
			resourceType := fmt.Sprintf(resourceTypeHostedFString, strings.ToLower(td.RepoFormat))
			resourceName := fmt.Sprintf(repoNameFString, resourceType)
			repoName := strings.ToLower(fmt.Sprintf(hostedNameFString, td.RepoFormat, randomString))

			var steps []resource.TestStep

			// 1. Create with minimal configuration relying on defaults
			step1Checks := append(
				// Test Case Specific Checks
				td.CheckFunc(resourceName),

				// Generic Checks
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_ONLINE, "true"),
				resource.TestCheckResourceAttrSet(resourceName, repotest.RES_ATTR_URL),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "false"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
				resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_ROUTING_RULE_NAME),
			)
			if td.SupportsProprietaryComponents {
				step1Checks = append(step1Checks, resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, "false"))
			} else {
				step1Checks = append(step1Checks, resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS))
			}
			steps = append(steps, resource.TestStep{
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, randomString, false, td.SupportsProprietaryComponents),
				Check:  resource.ComposeAggregateTestCheckFunc(step1Checks...),
			})

			// 2. Update to use full config
			step2Checks := append(
				// Test Case Specific Checks
				td.CheckFunc(resourceName),

				// Generic Checks
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_ONLINE, "true"),
				resource.TestCheckResourceAttrSet(resourceName, repotest.RES_ATTR_URL),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_STORAGE_WRITE_POLICY, common.WRITE_POLICY_ALLOW_ONCE),
				resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_CLEANUP_POLICY_COUNT, "0"),
				resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_ROUTING_RULE_NAME),
			)
			if td.SupportsProprietaryComponents {
				step2Checks = append(step2Checks, resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, fmt.Sprintf("%t", td.SupportsProprietaryComponents)))
			} else {
				step2Checks = append(step2Checks, resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS))
			}

			steps = append(steps, resource.TestStep{
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, randomString, true, td.SupportsProprietaryComponents),
				Check:  resource.ComposeAggregateTestCheckFunc(step2Checks...),
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

func TestAccRepositoryGenericHostedInvalidBlobStore(t *testing.T) {
	for _, repoFormat := range common.AllHostedFormats() {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resourceType := fmt.Sprintf(resourceTypeHostedFString, strings.ToLower(repoFormat))
		repoName := strings.ToLower(fmt.Sprintf(hostedNameFString, repoFormat, randomString))

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
	configBlockHostedDefault         string = ""
	configBlockHostedDefaultApt      string = "apt = { distribution = \"bionic\" }\napt_signing = { key_pair = \"something\" }"
	configBlockHostedDefaultCargo    string = configBlockHostedDefault
	configBlockHostedDefaultConan    string = configBlockHostedDefault
	configBlockHostedDefaultDocker   string = "docker = { force_basic_auth = false\nv1_enabled = false}"
	configBlockHostedDefaultGitLfs   string = configBlockHostedDefault
	configBlockHostedDefaultHelm     string = configBlockHostedDefault
	configBlockHostedDefaultMaven    string = "maven = { layout_policy = \"PERMISSIVE\"\nversion_policy = \"RELEASE\" }"
	configBlockHostedDefaultNpm      string = configBlockHostedDefault
	configBlockHostedDefaultNuget    string = configBlockHostedDefault
	configBlockHostedDefaultPypi     string = configBlockHostedDefault
	configBlockHostedDefaultR        string = configBlockHostedDefault
	configBlockHostedDefaultRaw      string = "raw = { content_disposition = \"ATTACHMENT\" }"
	configBlockHostedDefaultRubyGems string = configBlockHostedDefault
	configBlockHostedDefaultYum      string = "yum = { repo_data_depth = 1 }"
)

func formatSpecificHostedDefaultConfig(repoFormat string) string {
	switch repoFormat {
	case common.REPO_FORMAT_ANSIBLE_GALAXY:
		return configBlockHostedDefault
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
	case common.REPO_FORMAT_TERRAFORM:
		return repotest.ConfigBlockHostedDefaultTerraform
	case common.REPO_FORMAT_YUM:
		return configBlockHostedDefaultYum

	default:
		return ""
	}
}

func repositoryHostedResourceFullConfig(resourceType, repoName, formatSpecificConfig string, supportsProprietaryComponents bool) string {
	if supportsProprietaryComponents {
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
    proprietary_components = %t
  }
  %s
 }
`, resourceType, repoName, supportsProprietaryComponents, formatSpecificConfig)
	} else {
		return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  %s
 }
`, resourceType, repoName, formatSpecificConfig)
	}

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

func repositoryHostedResourceConfig(resourceType, repoName, repoFormat, randomString string, completeData, supportsProprietaryComponents bool) string {
	configBlock := formatSpecificHostedDefaultConfig(repoFormat)
	if completeData {
		return repositoryHostedResourceFullConfig(
			resourceType, repoName, configBlock, supportsProprietaryComponents,
		)
	} else {
		return repositoryHostedResourceMinimalConfigWithDefaults(
			resourceType, repoName, configBlock,
		)
	}
}
