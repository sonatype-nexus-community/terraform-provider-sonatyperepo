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
			// Create and Read testing
			{
				Config: getTestAccBlobStoreFileResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr("sonatyperepo_blob_store_file.bsf", "name", fmt.Sprintf("test-%s", randomString)),
					resource.TestCheckResourceAttr("sonatyperepo_blob_store_file.bsf", "path", fmt.Sprintf("path-%s", randomString)),
					resource.TestCheckResourceAttrSet("sonatyperepo_blob_store_file.bsf", "last_updated"),
				),
			},
			// Update and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_file" "bsf" {
  name = "test-%s"
  path = "path-updated-%s"
}
`, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated values
					resource.TestCheckResourceAttr("sonatyperepo_blob_store_file.bsf", "name", fmt.Sprintf("test-%s", randomString)),
					resource.TestCheckResourceAttr("sonatyperepo_blob_store_file.bsf", "path", fmt.Sprintf("path-updated-%s", randomString)),
					resource.TestCheckResourceAttrSet("sonatyperepo_blob_store_file.bsf", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getTestAccBlobStoreFileResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_blob_store_file" "bsf" {
  name = "test-%s"
  path = "path-%s"
}
`, randomString, randomString)
}
