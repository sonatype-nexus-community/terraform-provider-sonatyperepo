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
	"regexp"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataSourceContentSelector = "data.sonatyperepo_content_selector.cs"
)

func TestAccContentSelectorDataSource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Missing required argument
			{
				Config:      utils_test.ProviderConfig + `data "sonatyperepo_content_selector" "cs" {}`,
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
			// Test 2: Non-existent selector returns empty
			{
				Config: configContentSelectorDoesNotExist(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceContentSelector, "name"),
				),
			},
			// Test 3: Happy path - create selector resource and read via data source
			{
				Config: configContentSelectorResource(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceContentSelector, "name"),
					resource.TestCheckResourceAttr(dataSourceContentSelector, "name", fmt.Sprintf("tf-test-cs-%s", randomString)),
					resource.TestCheckResourceAttrSet(dataSourceContentSelector, "description"),
					resource.TestCheckResourceAttrSet(dataSourceContentSelector, "expression"),
				),
			},
		},
	})
}

func configContentSelectorDoesNotExist(suffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`data "sonatyperepo_content_selector" "cs" {
	name = "non-existent-content-selector-%s"
}`, suffix)
}

// Helper to create a content selector resource and read it via data source
func configContentSelectorResource(suffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_content_selector" "r" {
	name        = "tf-test-cs-%s"
	description = "Test content selector for acceptance test"
	expression  = "path =~ \".*\""
}

data "sonatyperepo_content_selector" "cs" {
	name       = sonatyperepo_content_selector.r.name
	depends_on = [sonatyperepo_content_selector.r]
}
`, suffix)
}
