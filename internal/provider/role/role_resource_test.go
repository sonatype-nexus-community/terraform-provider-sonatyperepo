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

package role_test

import (
	"fmt"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeRole    = "sonatyperepo_role"
	resourceNameRole    = "sonatyperepo_role.rl"
	attrID              = "id"
	attrName            = "name"
	attrDescription     = "description"
	attrPrivilegesCount = "privileges.#"
	attrRolesCount      = "roles.#"
)

func TestAccRoleResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildRoleResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameRole, attrID, fmt.Sprintf("my-test-role-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrName, fmt.Sprintf("My Test Role %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrDescription, "This is a test role"),
					resource.TestCheckResourceAttr(resourceNameRole, attrPrivilegesCount, "1"),
					resource.TestCheckResourceAttr(resourceNameRole, attrRolesCount, "1"),
				),
			},
			// Update to full configuration
			{
				Config: buildRoleResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated values
					resource.TestCheckResourceAttr(resourceNameRole, attrID, fmt.Sprintf("my-test-role-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrName, fmt.Sprintf("My Updated Role %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrDescription, "This is an updated test role"),
					resource.TestCheckResourceAttr(resourceNameRole, attrPrivilegesCount, "2"),
					resource.TestCheckResourceAttr(resourceNameRole, attrRolesCount, "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            resourceNameRole,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
				ImportStateId:           fmt.Sprintf("my-test-role-%s", randomString),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRoleResourceOnlyPrivileges(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildRoleResourceOnlyPrivileges(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameRole, attrID, fmt.Sprintf("my-test-role-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrName, fmt.Sprintf("My Test Role %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrDescription, "This is a test role"),
					resource.TestCheckResourceAttr(resourceNameRole, attrPrivilegesCount, "2"),
					resource.TestCheckResourceAttr(resourceNameRole, attrRolesCount, "0"),
				),
			},
		},
	})
}

func TestAccRoleResourceOnlyRoles(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildRoleResourceOnlyRoles(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameRole, attrID, fmt.Sprintf("my-test-role-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrName, fmt.Sprintf("My Test Role %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameRole, attrDescription, "This is a test role"),
					resource.TestCheckResourceAttr(resourceNameRole, attrPrivilegesCount, "0"),
					resource.TestCheckResourceAttr(resourceNameRole, attrRolesCount, "2"),
				),
			},
		},
	})
}

func buildRoleResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "rl" {
  id = "my-test-role-%s"
  name = "My Test Role %s"
  description = "This is a test role"
  privileges = [
    "nx-healthcheck-read"
  ]
  roles = [
    "nx-anonymous"
  ]
}
`, resourceTypeRole, randomString, randomString)
}

func buildRoleResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "rl" {
  id = "my-test-role-%s"
  name = "My Updated Role %s"
  description = "This is an updated test role"
  privileges = [
    "nx-healthcheck-read",
    "nx-healthcheck-summary-read"
  ]
  roles = [
    "nx-anonymous",
    "nx-admin"
  ]
}
`, resourceTypeRole, randomString, randomString)
}

func buildRoleResourceOnlyPrivileges(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "rl" {
  id = "my-test-role-%s"
  name = "My Test Role %s"
  description = "This is a test role"
  privileges = [
    "nx-healthcheck-read",
    "nx-healthcheck-summary-read"
  ]
}
`, resourceTypeRole, randomString, randomString)
}

func buildRoleResourceOnlyRoles(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "rl" {
  id = "my-test-role-%s"
  name = "My Test Role %s"
  description = "This is a test role"
  roles = [
    "nx-anonymous"
  ]
}
`, resourceTypeRole, randomString, randomString)
}
