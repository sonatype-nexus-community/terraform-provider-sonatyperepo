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
package system_test

import (
	"fmt"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSecurityUserToken = "sonatyperepo_security_user_tokens"
	resourceNameSecurityUserToken = "sonatyperepo_security_user_tokens.test"
)

func TestAccSecurityUserTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getSecurityUserTokenResourceConfig(true, 30, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "30"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "protect_content", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceNameSecurityUserToken,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				// Ignore last_updated since it will be different after import
				ImportStateVerifyIgnore: []string{"last_updated"},
				ImportStateId:           "user-tokens", // Can be any string for this singleton resource
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSecurityUserTokenResourceUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with initial configuration
			{
				Config: getSecurityUserTokenResourceConfig(true, 30, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "30"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "protect_content", "false"),
				),
			},
			// Update expiration days
			{
				Config: getSecurityUserTokenResourceConfig(true, 90, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "90"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "protect_content", "false"),
				),
			},
			// Enable protect_content
			{
				Config: getSecurityUserTokenResourceConfig(true, 90, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "90"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "protect_content", "true"),
				),
			},
			// Disable expiration
			{
				Config: getSecurityUserTokenResourceConfig(true, 90, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "90"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "protect_content", "true"),
				),
			},
			// Disable the feature entirely
			{
				Config: getSecurityUserTokenResourceConfig(false, 90, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccSecurityUserTokenResourceImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First, create a resource
			{
				Config: getSecurityUserTokenResourceConfig(true, 30, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "30"),
				),
			},
			// Test import with different import IDs (all should work for singleton resource)
			{
				ResourceName:                         resourceNameSecurityUserToken,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "user-tokens",
			},
			{
				ResourceName:                         resourceNameSecurityUserToken,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "import",
			},
			{
				ResourceName:                         resourceNameSecurityUserToken,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "any-string-works",
			},
		},
	})
}

func getSecurityUserTokenResourceConfig(enabled bool, expirationDays int, expirationEnabled bool, protectContent bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
	enabled = %t
	expiration_days = %d
	expiration_enabled = %t
	protect_content = %t
}
`, resourceTypeSecurityUserToken, enabled, expirationDays, expirationEnabled, protectContent)
}
