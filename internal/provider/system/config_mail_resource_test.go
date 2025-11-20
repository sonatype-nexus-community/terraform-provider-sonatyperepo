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
	resourceTypeConfigMail = "sonatyperepo_system_config_mail"
	resourceNameConfigMail = "sonatyperepo_system_config_mail.email"
)

func TestAccSystemConfigMailResource(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: getSytemConfigMailResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify
						resource.TestCheckResourceAttr(resourceNameConfigMail, "enabled", "false"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "host", fmt.Sprintf("something.tld.%s", randomString)),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "port", "587"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "username", "someone"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "password", "sensitive-value"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "from_address", "no-where@somewhere.tld"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "start_tls_enabled", "false"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "start_tls_required", "false"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "ssl_on_connect_enabled", "false"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "ssl_server_identity_check_enabled", "false"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "nexus_trust_store_enabled", "false"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "subject_prefix", "TESTING"),
					),
				},
				// ImportState testing
				{
					ResourceName:                         resourceNameConfigMail,
					ImportState:                          true,
					ImportStateId:                        "system-email-config",
					ImportStateVerify:                    true,
					ImportStateVerifyIgnore:              []string{"password", "last_updated"},
					ImportStateVerifyIdentifierAttribute: "host",
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func TestAccSystemConfigMailResourceImport(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create initial resource
				{
					Config: getSytemConfigMailResourceConfig(randomString),
				},
				// Test import with different import IDs
				{
					ResourceName:                         resourceNameConfigMail,
					ImportState:                          true,
					ImportStateId:                        "email-config",
					ImportStateVerify:                    true,
					ImportStateVerifyIgnore:              []string{"password", "last_updated"},
					ImportStateVerifyIdentifierAttribute: "host",
				},
				// Test import with another ID to verify flexibility
				{
					ResourceName:                         resourceNameConfigMail,
					ImportState:                          true,
					ImportStateId:                        "system-mail-settings",
					ImportStateVerify:                    true,
					ImportStateVerifyIgnore:              []string{"password", "last_updated"},
					ImportStateVerifyIdentifierAttribute: "host",
				},
			},
		})
	}
}

func TestAccSystemConfigMailResourceUpdateAfterImport(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
		updatedRandomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create initial resource
				{
					Config: getSytemConfigMailResourceConfig(randomString),
				},
				// Import the resource
				{
					ResourceName:                         resourceNameConfigMail,
					ImportState:                          true,
					ImportStateId:                        "imported-email-config",
					ImportStateVerify:                    true,
					ImportStateVerifyIgnore:              []string{"password", "last_updated"},
					ImportStateVerifyIdentifierAttribute: "host",
				},
				// Update after import
				{
					Config: getSytemConfigMailResourceConfigUpdated(updatedRandomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceNameConfigMail, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "host", fmt.Sprintf("updated.tld.%s", updatedRandomString)),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "port", "465"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "username", "updated-user"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "from_address", "updated@somewhere.tld"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "start_tls_enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "ssl_on_connect_enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameConfigMail, "subject_prefix", "UPDATED"),
					),
				},
			},
		})
	}
}

func getSytemConfigMailResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "email" {
	enabled = false
	host = "something.tld.%s"
	port = 587
	username = "someone"
	password = "sensitive-value"
	from_address = "no-where@somewhere.tld"
	start_tls_enabled = false
	start_tls_required = false
	ssl_on_connect_enabled = false
	ssl_server_identity_check_enabled = false
	nexus_trust_store_enabled = false
	subject_prefix = "TESTING"
}
`, resourceTypeConfigMail, randomString)
}

func getSytemConfigMailResourceConfigUpdated(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "email" {
	enabled = true
	host = "updated.tld.%s"
	port = 465
	username = "updated-user"
	password = "updated-sensitive-value"
	from_address = "updated@somewhere.tld"
	start_tls_enabled = true
	start_tls_required = false
	ssl_on_connect_enabled = true
	ssl_server_identity_check_enabled = false
	nexus_trust_store_enabled = false
	subject_prefix = "UPDATED"
}
`, resourceTypeConfigMail, randomString)
}
