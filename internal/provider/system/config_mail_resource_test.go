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

package system_test

import (
	"fmt"
	"os"
	"terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceType = "sonatyperepo_system_config_mail"
	resourceName = "sonatyperepo_system_config_mail.email"
)

func TestAccSystemConfigMailResource(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: getSsytemConfigMailResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify
						resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "host", fmt.Sprintf("something.tld.%s", randomString)),
						resource.TestCheckResourceAttr(resourceName, "port", "587"),
						resource.TestCheckResourceAttr(resourceName, "username", "someone"),
						resource.TestCheckResourceAttr(resourceName, "password", "sensitive-value"),
						resource.TestCheckResourceAttr(resourceName, "from_address", "no-where@somewhere.tld"),
						resource.TestCheckResourceAttr(resourceName, "start_tls_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "start_tls_required", "false"),
						resource.TestCheckResourceAttr(resourceName, "ssl_on_connect_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "ssl_server_identity_check_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "nexus_trust_store_enabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "subject_prefix", "TESTING"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func getSsytemConfigMailResourceConfig(randomString string) string {
	return fmt.Sprintf(utils.ProviderConfig+`
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
`, resourceType, randomString)
}
