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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	defaultIqServerUrl          = "http://localhost:8070"
	resourceTypeSysIqConnection = "sonatyperepo_system_iq_connection"
	resourceNameSysIqConnection = "sonatyperepo_system_iq_connection.iq"
)

func TestAccSystemIqConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: systemIqConnectionResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", defaultIqServerUrl),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", ""),
					resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
				),
			},
			// Test enabling show_iq_server_link
			{
				Config: systemIqConnectionResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", defaultIqServerUrl),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", ""),
					resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
				),
			},
			// Test adding Properties
			{
				Config: systemIqConnectionWithPropertiesResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", defaultIqServerUrl),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", "key1=value1&key2=value2"),
					resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
				),
			},
			// Test that nexus_trust_store_enabled is applied
			{
				Config: systemIqConnectionWithTrustStoreConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSysIqConnection,
						"nexus_trust_store_enabled", "true"),
				),
			},
			// Verify it can be set back to false
			{
				Config: systemIqConnectionWithTrustStoreConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSysIqConnection,
						"nexus_trust_store_enabled", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceNameSysIqConnection,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "url",
				ImportStateVerifyIgnore: []string{
					"password",
					"last_updated",
				},
				ImportStateId: "system-iq-config",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func systemIqConnectionResourceConfig(showIqLink bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = true
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "%s"
  username = "user"
  password = "token"
  show_iq_server_link = %t
}
`, resourceTypeSysIqConnection, defaultIqServerUrl, showIqLink)
}

func systemIqConnectionWithPropertiesResourceConfig(showIqLink bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = true
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "%s"
  username = "user"
  password = "token"
  show_iq_server_link = %t
  properties = "key1=value1&key2=value2"
}
`, resourceTypeSysIqConnection, defaultIqServerUrl, showIqLink)
}

func systemIqConnectionWithTrustStoreConfig(nexusTrustStoreEnabled bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method     = "USER"
  enabled                   = true
  fail_open_mode_enabled    = false
  nexus_trust_store_enabled = %t
  url 					    = "%s"
  username 					= "user"
  password                  = "token"
  show_iq_server_link       = true
}
`, resourceTypeSysIqConnection, nexusTrustStoreEnabled, defaultIqServerUrl)
}
