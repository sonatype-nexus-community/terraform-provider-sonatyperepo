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

package system

import (
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameLdap1 = "sonatyperepo_system_config_ldap_connection.ldap1"
	resourceNameLdap2 = "sonatyperepo_system_config_ldap_connection.ldap2"
)

func TestAccSystemConfigLdapResource(t *testing.T) {
	// randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: utils_test.ProviderConfig + `
resource "sonatyperepo_system_config_ldap_connection" "ldap1" {
  name = "Test LDAP Connection"
  protocol = "LDAP"
  hostname = "ldap.somewhere.tld"
  port = 389
  auth_scheme = "NONE"
  search_base = "something e"
  user_object_class = "inetOrgPerson"
  user_id_attribute = "uid"
  user_real_name_attribute = "name"
  user_email_name_attribute = "mail"
  connection_timeout = 999
}				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameLdap1, "name", "Test LDAP Connection"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "protocol", common.PROTOCOL_LDAP),
					resource.TestCheckResourceAttr(resourceNameLdap1, "hostname", "ldap.somewhere.tld"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "port", "389"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "auth_scheme", common.AUTH_SCHEME_NONE),
					resource.TestCheckResourceAttr(resourceNameLdap1, "search_base", "something e"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "user_object_class", "inetOrgPerson"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "user_id_attribute", "uid"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "user_real_name_attribute", "name"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "user_email_name_attribute", "mail"),
					resource.TestCheckResourceAttr(resourceNameLdap1, "connection_timeout", "999"),
				),
			},
			{
				Config: utils_test.ProviderConfig + `
resource "sonatyperepo_system_config_ldap_connection" "ldap2" {
  name = "Test LDAP Connection"
  protocol = "LDAPS"
  hostname = "ldap.somewhere.tld"
  port     = 636
  auth_scheme = "NONE"
  connection_retry_delay = 60
  connection_timeout     = 10
  nexus_trust_store_enabled = true
  map_ldap_groups_to_roles = true
  search_base   = "a-base"
  group_subtree = false
  group_type    = "DYNAMIC"
  user_base_dn                 = "ou=people"
  user_email_name_attribute = "mail"
  user_id_attribute            = "uid"
  user_ldap_filter             = ""
  user_member_of_attribute     = "memberOf"
  user_object_class            = "inetOrgPerson"
  user_password_attribute      = ""
  user_real_name_attribute     = "cn"
  user_subtree                 = false
}				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameLdap2, "name", "Test LDAP Connection"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "protocol", common.PROTOCOL_LDAPS),
					resource.TestCheckResourceAttr(resourceNameLdap2, "hostname", "ldap.somewhere.tld"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "port", "636"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "auth_scheme", common.AUTH_SCHEME_NONE),
					resource.TestCheckResourceAttr(resourceNameLdap2, "connection_retry_delay", "60"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "connection_timeout", "10"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "nexus_trust_store_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "map_ldap_groups_to_roles", "true"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "search_base", "a-base"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "group_subtree", "false"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "group_type", common.LDAP_GROUP_MAPPING_DYNAMIC),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_base_dn", "ou=people"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_email_name_attribute", "mail"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_id_attribute", "uid"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_ldap_filter", ""),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_member_of_attribute", "memberOf"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_object_class", "inetOrgPerson"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_password_attribute", ""),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_real_name_attribute", "cn"),
					resource.TestCheckResourceAttr(resourceNameLdap2, "user_subtree", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}
