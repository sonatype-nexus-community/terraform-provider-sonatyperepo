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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

const (
	matchersCount = "matchers.#"
)

func TestAccRoutingRuleResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "sonatyperepo_routing_rule.test"
	routingRuleName := fmt.Sprintf("test-routing-rule-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getTestAccRoutingRuleResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", routingRuleName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test routing rule"),
					resource.TestCheckResourceAttr(resourceName, "mode", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, matchersCount, "1"),
					resource.TestCheckResourceAttr(resourceName, "matchers.0", "^/com/example/.*"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        routingRuleName,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: getTestAccRoutingRuleResourceConfigUpdated(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", routingRuleName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated test routing rule"),
					resource.TestCheckResourceAttr(resourceName, "mode", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, matchersCount, "2"),
					resource.TestCheckResourceAttr(resourceName, "matchers.0", "^/com/example/.*"),
					resource.TestCheckResourceAttr(resourceName, "matchers.1", "^/org/test/.*"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRoutingRuleResourceMinimal(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "sonatyperepo_routing_rule.minimal"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with minimal configuration
			{
				Config: getTestAccRoutingRuleResourceMinimalConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("minimal-routing-rule-%s", randomString)),
					resource.TestCheckResourceAttr(resourceName, "mode", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, matchersCount, "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRoutingRuleResourceInvalidMode(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid configuration - invalid mode
			{
				Config:      getTestAccRoutingRuleResourceInvalidModeConfig(randomString),
				ExpectError: regexp.MustCompile("Attribute mode value must be one of"),
			},
		},
	})
}

func getTestAccRoutingRuleResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_routing_rule" "test" {
  name        = "test-routing-rule-%s"
  description = "Test routing rule"
  mode        = "BLOCK"
  matchers    = ["^/com/example/.*"]
}
`, randomString)
}

func getTestAccRoutingRuleResourceConfigUpdated(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_routing_rule" "test" {
  name        = "test-routing-rule-%s"
  description = "Updated test routing rule"
  mode        = "ALLOW"
  matchers    = ["^/com/example/.*", "^/org/test/.*"]
}
`, randomString)
}

func getTestAccRoutingRuleResourceMinimalConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_routing_rule" "minimal" {
  name        = "minimal-routing-rule-%s"
  description = "Minimal routing rule"
  mode        = "BLOCK"
  matchers    = ["^/com/minimal/.*"]
}
`, randomString)
}

func getTestAccRoutingRuleResourceInvalidModeConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_routing_rule" "invalid" {
  name        = "invalid-routing-rule-%s"
  description = "Invalid routing rule"
  mode        = "INVALID_MODE"
  matchers    = ["^/com/example/.*"]
}
`, randomString)
}
