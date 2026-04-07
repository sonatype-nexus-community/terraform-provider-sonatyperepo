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

package privilege_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceTypePrivilegeWildcard = "sonatyperepo_privilege_wildcard"
	resourceNamePrivilegeWildcard = fmt.Sprintf("%s.p", resourceTypeRepositoryContentSelector)
)

func TestAccPrivilegeWildcardResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildPrivilegeWildcardResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal configuration
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "name", fmt.Sprintf("test-priv-wildcard-%s", randomString)),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "description", "some description"),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "type", privilege_type.TypeWildcard.String()),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "pattern", "test-pattern"),
				),
			},
			// Update to full configuration
			{
				Config: buildPrivilegeWildcardResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full configuration
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "name", fmt.Sprintf("test-priv-wildcard-%s", randomString)),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "type", privilege_type.TypeWildcard.String()),
					resource.TestCheckResourceAttr(resourceTypePrivilegeWildcard, "pattern", "updated-pattern-*"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// ImportState testing
			{
				ResourceName:                         resourceTypePrivilegeWildcard,
				ImportStateId:                        fmt.Sprintf("test-priv-wildcard-%s", randomString),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func buildPrivilegeWildcardResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-wildcard-%s"
	description = "some description"
	pattern = "test-pattern"
}`, resourceTypePrivilegeWildcard, randomString)
}

func buildPrivilegeWildcardResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-wildcard-%s"
	description = "updated description"
	pattern = "updated-pattern-*"
}`, resourceTypePrivilegeWildcard, randomString)
}
