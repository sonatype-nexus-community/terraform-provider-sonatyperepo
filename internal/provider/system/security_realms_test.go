/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package system_test

import (
	"fmt"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSecurityRealms = "sonatyperepo_security_realms"
	resourceNameSecurityRealms = "sonatyperepo_security_realms.test"

	// Test attribute paths
	activeCount  = "active.#"
	activeIndex0 = "active.0"
	activeIndex1 = "active.1"
	activeIndex2 = "active.2"
	activeIndex3 = "active.3"
	activeIndex4 = "active.4"
	idAttr       = "id"
)

func TestAccSecurityRealmsResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing - single realm
			{
				Config: getSecurityRealmsResourceConfig(randomString, []string{"NexusAuthenticatingRealm"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify single realm configuration
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeCount, "1"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex0, "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, idAttr),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "security_realms",
			},
			// Update and Read testing - multiple realms
			{
				Config: getSecurityRealmsResourceConfig(randomString, []string{"DockerToken", "NexusAuthenticatingRealm"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify multiple realm configuration
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeCount, "2"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex0, "DockerToken"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex1, "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, idAttr),
				),
			},
			// ImportState testing after update
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "security_realms",
			},
			// Update and Read testing - reorder realms
			{
				Config: getSecurityRealmsResourceConfig(randomString, []string{"NexusAuthenticatingRealm", "DockerToken"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify reordered realm configuration
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeCount, "2"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex0, "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex1, "DockerToken"),
					resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, idAttr),
				),
			},
			// ImportState testing with alternative ID (tests silent override)
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "any_id_here",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSecurityRealmsResourceMinimalConfig(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing - minimal configuration
			{
				Config: getSecurityRealmsMinimalResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify minimal realm configuration
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeCount, "1"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex0, "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, idAttr),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "security_realms",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSecurityRealmsResourceCommonRealms(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing - common enterprise realms
			{
				Config: getSecurityRealmsResourceConfig(randomString, []string{
					"NexusAuthenticatingRealm",
					"LdapRealm",
					"DockerToken",
					"NpmToken",
					"NuGetApiKey",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify common enterprise realm configuration
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeCount, "5"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex0, "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex1, "LdapRealm"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex2, "DockerToken"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex3, "NpmToken"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex4, "NuGetApiKey"),
					resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, idAttr),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "security_realms",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccSecurityRealmsResourceImportOnly tests import functionality in isolation
func TestAccSecurityRealmsResourceImportOnly(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create initial state to import from
			{
				Config: getSecurityRealmsResourceConfig(randomString, []string{"NexusAuthenticatingRealm", "DockerToken"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeCount, "2"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex0, "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttr(resourceNameSecurityRealms, activeIndex1, "DockerToken"),
				),
			},
			// Test import with standard ID
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "security_realms",
			},
			// Test import with empty string ID
			{
				ResourceName:      resourceNameSecurityRealms,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "",
			},
		},
	})
}

func getSecurityRealmsResourceConfig(randomString string, activeRealms []string) string {
	realmsConfig := ""
	for _, realm := range activeRealms {
		realmsConfig += fmt.Sprintf(`    "%s",`+"\n", realm)
	}
	// Remove trailing comma and newline
	if len(realmsConfig) > 0 {
		realmsConfig = realmsConfig[:len(realmsConfig)-2] + "\n"
	}

	return fmt.Sprintf(utils_test.ProviderConfig+`
		resource "%s" "test" {
		active = [
		%s  ]
		}
		`, resourceTypeSecurityRealms, realmsConfig)
}

func getSecurityRealmsMinimalResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
		resource "%s" "test" {
		active = ["NexusAuthenticatingRealm"]
		}
		`, resourceTypeSecurityRealms)
}
