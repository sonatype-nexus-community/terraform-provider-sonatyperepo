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
	"terraform-provider-sonatyperepo/internal/provider/testutil"
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
		RepoFormat:                    common.REPO_FORMAT_APT,
		SchemaFunc:                    repositoryHostedResourceConfig,
		SupportsProprietaryComponents: false,
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
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_DOCKER_V1_ENABLED, "false"),
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
				resource.TestCheckNoResourceAttr(resourceName, RES_ATTR_MAVEN_CONTENT_DISPOSITION),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_MAVEN_LAYOUT_POLICY, common.MAVEN_LAYOUT_PERMISSIVE),
				resource.TestCheckResourceAttr(resourceName, RES_ATTR_MAVEN_VERSION_POLICY, common.MAVEN_VERSION_POLICY_RELEASE),
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
			steps = append(steps, resource.TestStep{
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, randomString, false, td.SupportsProprietaryComponents),
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
				Config: td.SchemaFunc(resourceType, repoName, td.RepoFormat, randomString, true, td.SupportsProprietaryComponents),
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
						resource.TestCheckResourceAttr(resourceName, RES_ATTR_COMPONENT_PROPRIETARY_COMPONENTS, fmt.Sprintf("%t", td.SupportsProprietaryComponents)),
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
	configBlockHostedDefaultApt       string = "apt = { distribution = \"bionic\" }\napt_signing = { key_pair = \"something\" }"
	configBlockHostedDefaultCargo     string = ""
	configBlockHostedDefaultConan     string = ""
	configBlockHostedDefaultDocker    string = "docker = { force_basic_auth = false\nv1_enabled = false}"
	configBlockHostedDefaultGitLfs    string = ""
	configBlockHostedDefaultHelm      string = ""
	configBlockHostedDefaultMaven     string = "maven = { layout_policy = \"PERMISSIVE\"\nversion_policy = \"RELEASE\" }"
	configBlockHostedDefaultNpm       string = ""
	configBlockHostedDefaultNuget     string = ""
	configBlockHostedDefaultPypi      string = ""
	configBlockHostedDefaultR         string = ""
	configBlockHostedDefaultRaw       string = "raw = { content_disposition = \"ATTACHMENT\" }"
	configBlockHostedDefaultRubyGems  string = ""
	configBlockHostedDefaultTerraform string = `terraform_signing = { signing_key = <<-EOT
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQHYBGmLEQEBBADBLrTiM/XmBoTSBTdGRSMFgqM12vVi0+3K2vMk9Zd+HUN3O0zY
ho0Q16SQU9hY7eWRXp/XiyL59u7HQhtrBq36dthvZTCPh23G3ldCtlruPhQWtHI/
xO3phio8skST4MDRfS3csoyRc/rnY7Rc00P7J8HP7dx+sRqv+SnIBeyOLQARAQAB
AAP9E3Q4Z4IrjGlSJVM8pIEwXGzyMil1Ziko7HF9pFZuFddtFJv+alysZoqMyjMD
WbtFT80bZCmhEVKWa68C01WWHfK2CqPOsEFiWG/fxUbnUG7RlehMKrI6KF+2wWBv
o452loV/Bzua64uR1kP+l43BH69LzJE6uWHl5KNJyX1uoskCAMvp9kkzc1Pe2/hT
Vc72s/CkMlw6GMSI4Lk6+YuvajGlr/HxsFhBjM9ADLkWIDoywxCQ1kKSxtF/FG4a
zZG3GxkCAPKHAh6ByWSc+dfg1acQx1/LHaGdmLACaJYK7OAy4+ra+VrX6c3th6ye
T8SzJG2sq3aBztDBwdtjdWf+8BazwjUCALcVncOFj5a/N1vZZ6chuo27wEVw6Bpb
iN2rb+SxuT1iTaCE3/RfSywlqf0aVMkh3Dygz5/CwOwEmffNVGhV9gOpM7RpU29u
YXR5cGUgQ29tbXVuaXR5IFRlc3QgRGF0YSBLZXkgKFRFU1QgREFUQSBPTkxZIC0g
RE8gTk9UIFVTRSBGT1IgU0lHTklORykgPGNvbW11bml0eS1ncm91cEBzb25hdHlw
ZS5jb20+iM4EEwEIADgWIQTHqk34+snIQVhkIMl2y863mbo3rgUCaYsRAQIbAwUL
CQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRB2y863mbo3rqk8A/97anQKQ5ZUc/Oz
FSUpRI7sKCwot2C9dP7wAAifUtfX7vn32H4fz3T9BDB8CJMtVurMYVNskLvy5rKe
n//joo0cp+XX+KDVmuEPtxbZD5+Py+JGUwOuOkK8bO5N/xGe0N/CWfW2DXvpMvut
1x/uN9aslm9GBhezNr1V5totT7Tx5p0B2ARpixEBAQQAzmA6ajRv7U3tQV9aJav1
/y/+byUrm4hta2pGuo0qTeP8i2SEao3S7DkoENcA3MhRtxJk4fzMGmH68saQyKK5
66se0mZLTQjkPKGXje0pjAT4hTaffi3PDycR5MBEe80rbu08ouqSkKQ5xUrTL3FE
MC9BBbFMacl8EAbeIiIyvLEAEQEAAQAD+QGysjOuIRA2jpMrGj2NEXlHMJXYZuLu
PkoRR09TrU9pDB/skf1DXm2OmkCFDVsJBqjt2hYaLPF9YNnF0JqHAhENSJqShPgd
tyRprrObKuTYTH+cv4QI1cr52Oxr6BkAqP3VqPJxqrqXWLnscveryoxEPMlNvXbk
5ATexXyThwRBAgDlorsxh4YywRJdrQCSSnQiiYlqt/L2cKliTedbEN3ffFOD3OH/
zZbaoXruev75FrIZtgtfSgprLELw/fwsTxHBAgDmEd01R3S8R3tpkzMGvdCwZFL2
6uGaVmgTZ415XpXiNrDWSO0QeD15F6mnMwM8PsEEarRwrnc23pPZSLw2eIbxAf9D
o/bDNpno4GBpCd6P8sUhFRCw8UweU4EHVfz7OfnBkid8tvn4y85U2HJUi9jXj4/v
+yDRM+uhsch4VBac5xhVpwOItgQYAQgAIBYhBMeqTfj6ychBWGQgyXbLzreZujeu
BQJpixEBAhsMAAoJEHbLzreZujeu0o0D/iHzDEXpkHE/sbd82JwPR48YR8cmBzq+
CMnhPvAykyWgXvRoEQmXj+rzH9nlTsD9TNIVrnReTT5PBDGWVTk5DXkpb9ZfaajU
USlNrkrzgatlosKJskQSSrKSbmqju1/R885DZMtTb4ryjzHVqvwALzCFTXyyEpMl
OVIAhWNMezYE
=RbgK
-----END PGP PRIVATE KEY BLOCK-----
EOT
}`
	configBlockHostedDefaultYum string = "yum = { repo_data_depth = 1 }"
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
	case common.REPO_FORMAT_TERRAFORM:
		return configBlockHostedDefaultTerraform
	case common.REPO_FORMAT_YUM:
		return configBlockHostedDefaultYum

	default:
		return ""
	}
}

func repositoryHostedResourceFullConfig(resourceType, repoName, formatSpecificConfig string, supportsProprietaryComponents bool) string {
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
