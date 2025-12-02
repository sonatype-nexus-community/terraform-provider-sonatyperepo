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

package repository_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-sonatyperepo/internal/provider/repository"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

func TestAccRoutingRulesDataSource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	dataSourceName := "data.sonatyperepo_routing_rules.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create multiple routing rules and list them
			{
				Config: testAccRoutingRulesDataSourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that we have at least our created routing rules
					resource.TestCheckResourceAttrSet(dataSourceName, "routing_rules.#"),
					// Verify the list contains routing rules with expected attributes
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName,
						"routing_rules.*",
						map[string]string{
							"name": fmt.Sprintf("test-routing-rule-ds1-%s", randomString),
							"mode": repository.RoutingRuleModeBlock,
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName,
						"routing_rules.*",
						map[string]string{
							"name": fmt.Sprintf("test-routing-rule-ds2-%s", randomString),
							"mode": repository.RoutingRuleModeAllow,
						},
					),
				),
			},
		},
	})
}

func testAccRoutingRulesDataSourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_routing_rule" "test1" {
  name        = "test-routing-rule-ds1-%s"
  description = "Test routing rule 1"
  mode        = "%s"
  matchers    = ["^/com/test1/.*"]
}

resource "sonatyperepo_routing_rule" "test2" {
  name        = "test-routing-rule-ds2-%s"
  description = "Test routing rule 2"
  mode        = "%s"
  matchers    = ["^/com/test2/.*"]
}

data "sonatyperepo_routing_rules" "test" {
  depends_on = [
    sonatyperepo_routing_rule.test1,
    sonatyperepo_routing_rule.test2
  ]
}
`, randomString, repository.RoutingRuleModeBlock, randomString, repository.RoutingRuleModeAllow)
}
