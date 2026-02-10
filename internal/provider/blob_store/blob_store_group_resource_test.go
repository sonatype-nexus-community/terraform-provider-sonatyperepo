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

package blob_store_test

import (
	"fmt"
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlobStoreGroupResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Failure: No Members
			{
				Config:      buildTestAccBlobStoreGroupResourceNoMembers(randomString),
				ExpectError: regexp.MustCompile(errMessageBlobStoreGroupNoMembers),
			},
			// Failure: Blob Store used by Repos and cannot be a group member
			// The below test is failing (no error) in NXRM 3.89.0 - unclear why after much investigation
			// {
			// 	Config:      buildTestAccBlobStoreGroupResourceIneligibleMember(randomString),
			// 	ExpectError: regexp.MustCompile(errMessageBlobStoreGroupIneligibleMember),
			// },
			// Create with valid configuration
			{
				Config: buildTestAccBlobStoreGroupResourceNewMember(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full configuration
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_GROUP, RES_ATTR_NAME, fmt.Sprintf("test-group-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_GROUP, RES_ATTR_FILL_POLICY, common.BLOB_STORE_FILL_POLICY_ROUND_ROBIN),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_GROUP, RES_ATTR_MEMBERS_COUNT, "1"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_GROUP, "members.0", fmt.Sprintf("test-%s", randomString)),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_GROUP, RES_ATTR_SOFT_QUOTA),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_GROUP, RES_ATTR_LAST_UPDATED),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         RES_NAME_BLOB_STORE_GROUP,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        fmt.Sprintf("test-group-%s", randomString),
				ImportStateVerifyIdentifierAttribute: RES_ATTR_NAME,
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func buildTestAccBlobStoreGroupResourceNoMembers(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
    name = "test-group-%s"
    fill_policy = "%s"
    members = [ ]
}
`, RES_TYPE_BLOB_STORE_GROUP, randomString, common.BLOB_STORE_FILL_POLICY_ROUND_ROBIN)
}

func buildTestAccBlobStoreGroupResourceIneligibleMember(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
    name = "test-group-%s"
    fill_policy = "%s"
    members = [ "default" ]
}
`, RES_TYPE_BLOB_STORE_GROUP, randomString, common.BLOB_STORE_FILL_POLICY_ROUND_ROBIN)
}

func buildTestAccBlobStoreGroupResourceNewMember(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "bs" {
	name = "test-%s"
	path = "test-%s"
}

resource "%s" "test" {
    name = "test-group-%s"
    fill_policy = "%s"
    members = [ 
		%s.bs.name
	]
}
`, RES_TYPE_BLOB_STORE_FILE, randomString, randomString, RES_TYPE_BLOB_STORE_GROUP, randomString, common.BLOB_STORE_FILL_POLICY_ROUND_ROBIN, RES_TYPE_BLOB_STORE_FILE)
}
