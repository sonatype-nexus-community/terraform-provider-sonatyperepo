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
	resourceTypeSysIqConnection = "sonatyperepo_system_iq_connection"
	resourceNameSysIqConnection = "sonatyperepo_system_iq_connection.iq"
)

func TestAccSystemIqConnectionResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: systemIqConnectionResourceConfig(randomString, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", fmt.Sprintf("https://%s.somewhere.tld", randomString)),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", ""),
					resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
				),
			},
			// Test enabling show_iq_server_link
			{
				Config: systemIqConnectionResourceConfig(randomString, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", fmt.Sprintf("https://%s.somewhere.tld", randomString)),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", ""),
					resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
				),
			},
			// Test adding Properties
			{
				Config: systemIqConnectionWithPropertiesResourceConfig(randomString, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", fmt.Sprintf("https://%s.somewhere.tld", randomString)),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "true"),
					resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", "key1=value1&key2=value2"),
					resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func systemIqConnectionResourceConfig(randomString string, showIqLink bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = false
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "https://%s.somewhere.tld"
  username = "user"
  password = "token"
  show_iq_server_link = %t
}
`, resourceTypeSysIqConnection, randomString, showIqLink)
}

func systemIqConnectionWithPropertiesResourceConfig(randomString string, showIqLink bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = false
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "https://%s.somewhere.tld"
  username = "user"
  password = "token"
  show_iq_server_link = %t
  properties = "key1=value1&key2=value2"
}
`, resourceTypeSysIqConnection, randomString, showIqLink)
}
