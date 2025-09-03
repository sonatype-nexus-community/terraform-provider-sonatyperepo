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
	"os"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSecurityRealms = "sonatyperepo_security_realms"
	resourceNameSecurityRealms = "sonatyperepo_security_realms.test"
)

func TestAccSecurityRealmsResource(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing - single realm
				{
					Config: getSecurityRealmsResourceConfig(randomString, []string{"NexusAuthenticatingRealm"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify single realm configuration
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.#", "1"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.0", "NexusAuthenticatingRealm"),
						resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, "id"),
					),
				},
				// Update and Read testing - multiple realms
				{
					Config: getSecurityRealmsResourceConfig(randomString, []string{"DockerToken", "NexusAuthenticatingRealm"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify multiple realm configuration
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.#", "2"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.0", "DockerToken"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.1", "NexusAuthenticatingRealm"),
						resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, "id"),
					),
				},
				// Update and Read testing - reorder realms
				{
					Config: getSecurityRealmsResourceConfig(randomString, []string{"NexusAuthenticatingRealm","DockerToken"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify reordered realm configuration
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.#", "2"),						
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.0", "NexusAuthenticatingRealm"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.1", "DockerToken"),
						resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, "id"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func TestAccSecurityRealmsResource_MinimalConfig(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing - minimal configuration
				{
					Config: getSecurityRealmsMinimalResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify minimal realm configuration
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.#", "1"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.0", "NexusAuthenticatingRealm"),
						resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, "id"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func TestAccSecurityRealmsResource_CommonRealms(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
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
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.#", "5"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.0", "NexusAuthenticatingRealm"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.1", "LdapRealm"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.2", "DockerToken"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.3", "NpmToken"),
						resource.TestCheckResourceAttr(resourceNameSecurityRealms, "active.4", "NuGetApiKey"),
						resource.TestCheckResourceAttrSet(resourceNameSecurityRealms, "id"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func getSecurityRealmsResourceConfig(randomString string, activeRealms []string) string {
	realmsConfig := ""
	for _, realm := range activeRealms {
		realmsConfig += fmt.Sprintf(`    "%s",` + "\n", realm)
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