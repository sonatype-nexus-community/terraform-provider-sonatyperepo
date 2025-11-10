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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataSourceTypeSecurityUserToken = "sonatyperepo_security_user_token"
	dataSourceNameSecurityUserToken = "data.sonatyperepo_security_user_token.test"
)

func TestAccSecurityUserTokenDataSource(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create a resource first
				{
					Config: getSecurityUserTokenDataSourceConfigWithResource(true, 45, true, false),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify the resource
						resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_days", "45"),
						resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "expiration_enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSecurityUserToken, "protect_content", "false"),
						// Verify the data source
						resource.TestCheckResourceAttr(dataSourceNameSecurityUserToken, "enabled", "true"),
						resource.TestCheckResourceAttr(dataSourceNameSecurityUserToken, "expiration_days", "45"),
						resource.TestCheckResourceAttr(dataSourceNameSecurityUserToken, "expiration_enabled", "true"),
						resource.TestCheckResourceAttr(dataSourceNameSecurityUserToken, "protect_content", "false"),
					),
				},
			},
		})
	}
}

func TestAccSecurityUserTokenDataSourceStandalone(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read the current configuration using data source only
				{
					Config: getSecurityUserTokenDataSourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Just verify that the data source can read the current state
						resource.TestCheckResourceAttrSet(dataSourceNameSecurityUserToken, "enabled"),
					),
				},
			},
		})
	}
}

func getSecurityUserTokenDataSourceConfigWithResource(enabled bool, expirationDays int, expirationEnabled bool, protectContent bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
	enabled = %t
	expiration_days = %d
	expiration_enabled = %t
	protect_content = %t
}

data "%s" "test" {
	depends_on = [%s.test]
}
`, resourceTypeSecurityUserToken, enabled, expirationDays, expirationEnabled, protectContent,
		dataSourceTypeSecurityUserToken, resourceTypeSecurityUserToken)
}

func getSecurityUserTokenDataSourceConfig() string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
data "%s" "test" {
}
`, dataSourceTypeSecurityUserToken)
}
