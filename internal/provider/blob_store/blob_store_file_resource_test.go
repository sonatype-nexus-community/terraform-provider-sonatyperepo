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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlobStoreFileResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: buildTestAccBlobStoreFileResourceMinimal(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal configuration
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_NAME, fmt.Sprintf("test-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_PATH, fmt.Sprintf("path-%s", randomString)),
					resource.TestCheckNoResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_SOFT_QUOTA),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_FILE, RES_ATTR_LAST_UPDATED),
				),
			},
			// Create with full configuration
			{
				Config: buildTestAccBlobStoreFileResourceComplete(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify full configuration
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_NAME, fmt.Sprintf("test-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_PATH, fmt.Sprintf("path-complete-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_SOFT_QUOTA_TYPE, "spaceRemainingQuota"),
					resource.TestCheckResourceAttr(RES_NAME_BLOB_STORE_FILE, RES_ATTR_SOFT_QUOTA_LIMIT, "1099511627776"),
					resource.TestCheckResourceAttrSet(RES_NAME_BLOB_STORE_FILE, RES_ATTR_LAST_UPDATED),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func buildTestAccBlobStoreFileResourceMinimal(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-%s"
  path = "path-%s"
}
`, RES_TYPE_BLOB_STORE_FILE, randomString, randomString)
}

func buildTestAccBlobStoreFileResourceComplete(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name = "test-%s"
  path = "path-complete-%s"
  soft_quota = {
    type  = "spaceRemainingQuota"
    limit = 1099511627776
  }
}
`, RES_TYPE_BLOB_STORE_FILE, randomString, randomString)
}
