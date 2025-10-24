/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package system_test

import (
	"fmt"
	"os"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSysAnonymousAccess = "sonatyperepo_system_anonymous_access"
	resourceNameSysAnonymousAccess = "sonatyperepo_system_anonymous_access.aa"
	anonymousUserIDFormat          = "anonymous-%s"
)

func TestAccSystemAnonymousAccessResource(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: getSystemAnonymousAccessResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "realm_name", common.DEFAULT_REALM_NAME),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "user_id", fmt.Sprintf(anonymousUserIDFormat, randomString)),
					),
				},
				// ImportState testing
				{
					ResourceName:                         resourceNameSysAnonymousAccess,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "user_id",
					// Ignore last_updated since it will be different after import
					ImportStateVerifyIgnore: []string{"last_updated"},
					ImportStateId:           "anonymous_access", // Can be any string for this singleton resource
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func TestAccSystemAnonymousAccessResourceImport(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// First, create a resource
				{
					Config: getSystemAnonymousAccessResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "realm_name", common.DEFAULT_REALM_NAME),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "user_id", fmt.Sprintf(anonymousUserIDFormat, randomString)),
					),
				},
				// Test import with different import IDs (all should work for singleton resource)
				{
					ResourceName:                         resourceNameSysAnonymousAccess,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "user_id",
					ImportStateVerifyIgnore:              []string{"last_updated"},
					ImportStateId:                        "anonymous_access",
				},
				{
					ResourceName:                         resourceNameSysAnonymousAccess,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "user_id",
					ImportStateVerifyIgnore:              []string{"last_updated"},
					ImportStateId:                        "import",
				},
				{
					ResourceName:                         resourceNameSysAnonymousAccess,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "user_id",
					ImportStateVerifyIgnore:              []string{"last_updated"},
					ImportStateId:                        "any-string-works",
				},
			},
		})
	}
}

func TestAccSystemAnonymousAccessResourceUpdate(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		updatedRandomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: getSystemAnonymousAccessResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "realm_name", common.DEFAULT_REALM_NAME),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "user_id", fmt.Sprintf(anonymousUserIDFormat, randomString)),
					),
				},
				// Update and Read testing
				{
					Config: getSystemAnonymousAccessResourceConfig(updatedRandomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "realm_name", common.DEFAULT_REALM_NAME),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "user_id", fmt.Sprintf(anonymousUserIDFormat, updatedRandomString)),
					),
				},
				// Test import after update
				{
					ResourceName:                         resourceNameSysAnonymousAccess,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "user_id",
					ImportStateVerifyIgnore:              []string{"last_updated"},
					ImportStateId:                        "post-update-import",
				},
				// Test disabling anonymous access
				{
					Config: getSystemAnonymousAccessResourceConfigDisabled(updatedRandomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "enabled", "false"),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "realm_name", common.DEFAULT_REALM_NAME),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "user_id", fmt.Sprintf(anonymousUserIDFormat, updatedRandomString)),
					),
				},
				// Test import when disabled
				{
					ResourceName:                         resourceNameSysAnonymousAccess,
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "user_id",
					ImportStateVerifyIgnore:              []string{"last_updated"},
					ImportStateId:                        "disabled-state-import",
				},
			},
		})
	}
}

func getSystemAnonymousAccessResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "aa" {
	enabled = true
	realm_name = "%s"
	user_id = "%s"
}
`, resourceTypeSysAnonymousAccess, common.DEFAULT_REALM_NAME, fmt.Sprintf(anonymousUserIDFormat, randomString))
}

func getSystemAnonymousAccessResourceConfigDisabled(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "aa" {
	enabled = false
	realm_name = "%s"
	user_id = "%s"
}
`, resourceTypeSysAnonymousAccess, common.DEFAULT_REALM_NAME, fmt.Sprintf(anonymousUserIDFormat, randomString))
}