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

package content_selector_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

const (
	dataSourceContentSelectors = "data.sonatyperepo_content_selectors.cses"
)

func TestAccContentSelectorsDataSource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	contentSelectorName := fmt.Sprintf("tf-test-cs-list-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Create content selector and verify it appears in list
			{
				Config: testAccContentSelectorsDataSourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify list is populated
					resource.TestCheckResourceAttrSet(dataSourceContentSelectors, "content_selectors.#"),
					// Verify created content selector appears with expected attributes
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceContentSelectors,
						"content_selectors.*",
						map[string]string{
							"name": contentSelectorName,
						},
					),
				),
			},
		},
	})
}

func testAccContentSelectorsDataSourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_content_selector" "test" {
	name        = "tf-test-cs-list-%s"
	description = "Test content selector for list data source"
	expression  = "format == \"raw\""
}

data "sonatyperepo_content_selectors" "cses" {
	depends_on = [sonatyperepo_content_selector.test]
}
`, randomString)
}
