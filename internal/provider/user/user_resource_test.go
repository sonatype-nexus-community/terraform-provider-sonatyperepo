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

func TestAccRoleResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getUserResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameUser, "user_id", fmt.Sprintf("acc-test-user-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, "first_name", fmt.Sprintf("Acc Test %s", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, "last_name", "User"),
					resource.TestCheckResourceAttr(resourceNameUser, "email_address", fmt.Sprintf("acc-test-%s@local", randomString)),
					resource.TestCheckResourceAttr(resourceNameUser, "status", "active"),
					resource.TestCheckResourceAttr(resourceNameUser, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceNameUser, "source", common.DEFAULT_USER_SOURCE),
					resource.TestCheckResourceAttr(resourceNameUser, "roles.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func getUserResourceConfig(randomString string) string {
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
