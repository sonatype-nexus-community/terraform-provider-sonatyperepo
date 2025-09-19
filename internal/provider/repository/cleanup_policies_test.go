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
	criteriaLastBlobUpdated = "criteria.last_blob_updated"
	criteriaAssetRegex      = "criteria.asset_regex"
)

func TestAccCleanupPolicyResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "sonatyperepo_cleanup_policy.test"
	cleanupPolicyName := fmt.Sprintf("test-cleanup-policy-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getTestAccCleanupPolicyResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", cleanupPolicyName),
					resource.TestCheckResourceAttr(resourceName, "format", "maven2"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Test cleanup policy"),
					resource.TestCheckResourceAttr(resourceName, criteriaLastBlobUpdated, "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        cleanupPolicyName,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: getTestAccCleanupPolicyResourceConfigUpdated(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", cleanupPolicyName),
					resource.TestCheckResourceAttr(resourceName, "format", "maven2"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Updated test cleanup policy"),
					resource.TestCheckResourceAttr(resourceName, criteriaLastBlobUpdated, "60"),
					resource.TestCheckResourceAttr(resourceName, criteriaAssetRegex, ".*\\.war$"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCleanupPolicyResourceMinimal(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "sonatyperepo_cleanup_policy.minimal"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with minimal configuration
			{
				Config: getTestAccCleanupPolicyResourceMinimalConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("minimal-cleanup-policy-%s", randomString)),
					resource.TestCheckResourceAttr(resourceName, "format", "maven2"),
					resource.TestCheckResourceAttr(resourceName, criteriaLastBlobUpdated, "30"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCleanupPolicyResourceInvalidFormat(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid configuration - no valid criteria
			{
				Config:      getTestAccCleanupPolicyResourceInvalidConfig(randomString),
				ExpectError: regexp.MustCompile("Invalid cleanup policy configuration"),
			},
		},
	})
}

func getTestAccCleanupPolicyResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_cleanup_policy" "test" {
  name   = "test-cleanup-policy-%s"
  format = "maven2"
  notes  = "Test cleanup policy"

  criteria = {
    last_blob_updated = 30
  }
}
`, randomString)
}

func getTestAccCleanupPolicyResourceConfigUpdated(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_cleanup_policy" "test" {
  name   = "test-cleanup-policy-%s"
  format = "maven2"
  notes  = "Updated test cleanup policy"

  criteria = {
    last_blob_updated = 60
    asset_regex       = ".*\\.war$"
  }
}
`, randomString)
}

func getTestAccCleanupPolicyResourceMinimalConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_cleanup_policy" "minimal" {
  name   = "minimal-cleanup-policy-%s"
  format = "maven2"
  
  criteria = {
    last_blob_updated = 30
  }
}
`, randomString)
}

func getTestAccCleanupPolicyResourceInvalidConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_cleanup_policy" "invalid" {
  name   = "invalid-cleanup-policy-%s"
  format = "maven2"
  
  criteria = {
  }
}
`, randomString)
}