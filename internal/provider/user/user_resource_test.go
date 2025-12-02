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

package user_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceTypeUser = "sonatyperepo_user"
	resourceNameUser = fmt.Sprintf("%s.u", resourceTypeUser)
)

const (
	attrUserID       = "user_id"
	attrFirstName    = "first_name"
	attrLastName     = "last_name"
	attrEmailAddress = "email_address"
	attrStatus       = "status"
	attrReadOnly     = "read_only"
	attrSource       = "source"
	attrRolesCount   = "roles.#"
)

func TestAccUserResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration (no password)
			{
				Config: buildUserResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal configuration
					resource.TestCheckResourceAttr(resourceNameUser, attrUserID, fmt.Sprintf("acc-test-user-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrFirstName, fmt.Sprintf("Acc Test %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrLastName, "User"),
					resource.TestCheckResourceAttr(resourceNameUser, attrEmailAddress, fmt.Sprintf("acc-test-%s@local", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrStatus, "active"),
					resource.TestCheckResourceAttr(resourceNameUser, attrReadOnly, "false"),
					resource.TestCheckResourceAttr(resourceNameUser, attrSource, common.DEFAULT_USER_SOURCE),
					resource.TestCheckResourceAttr(resourceNameUser, attrRolesCount, "1"),
				),
			},
			// Create with full configuration
			{
				Config: buildUserResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full configuration
					resource.TestCheckResourceAttr(resourceNameUser, attrUserID, fmt.Sprintf("acc-test-user-complete-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrFirstName, fmt.Sprintf("Complete %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrLastName, "TestUser"),
					resource.TestCheckResourceAttr(resourceNameUser, attrEmailAddress, fmt.Sprintf("complete-%s@local", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrStatus, "active"),
					resource.TestCheckResourceAttr(resourceNameUser, attrReadOnly, "false"),
					resource.TestCheckResourceAttr(resourceNameUser, attrSource, common.DEFAULT_USER_SOURCE),
					resource.TestCheckResourceAttr(resourceNameUser, attrRolesCount, "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func TestAccUserResourceUpdate(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildUserResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUser, attrUserID, fmt.Sprintf("acc-test-user-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrFirstName, fmt.Sprintf("Acc Test %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrLastName, "User"),
					resource.TestCheckResourceAttr(resourceNameUser, attrStatus, "active"),
					resource.TestCheckResourceAttr(resourceNameUser, attrRolesCount, "1"),
				),
			},
			// Update to full configuration
			{
				Config: buildUserResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameUser, attrUserID, fmt.Sprintf("acc-test-user-complete-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrFirstName, fmt.Sprintf("Complete %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrLastName, "TestUser"),
					resource.TestCheckResourceAttr(resourceNameUser, attrEmailAddress, fmt.Sprintf("complete-%s@local", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, attrStatus, "active"),
					resource.TestCheckResourceAttr(resourceNameUser, attrRolesCount, "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func buildUserResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "u" {
  user_id = "acc-test-user-%s"
  first_name = "Acc Test %s"
  last_name = "User"
  email_address = "acc-test-%s@local"
  password = "Password"
  status = "active"
  roles = [
    "nx-anonymous"
  ]
}
`, resourceTypeUser, randomString, randomString, randomString)
}

func buildUserResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "u" {
  user_id = "acc-test-user-complete-%s"
  first_name = "Complete %s"
  last_name = "TestUser"
  email_address = "complete-%s@local"
  password = "CompletePassword"
  status = "active"
  roles = [
    "nx-anonymous",
    "nx-admin"
  ]
}
`, resourceTypeUser, randomString, randomString, randomString)
}
