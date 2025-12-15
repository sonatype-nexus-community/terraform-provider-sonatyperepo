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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	dataSourceTypeSecurityUserToken = "sonatyperepo_security_user_tokens"
	dataSourceNameSecurityUserToken = "data.sonatyperepo_security_user_tokens.test"
)

func TestAccSecurityUserTokenDataSource(t *testing.T) {
	// This test reads the current security user token configuration
	// Note: Resource creation/modification is not tested as it would lock out the test user
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test 1: Read current data source state
			{
				Config: getSecurityUserTokenDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify data source can read the current state
					resource.TestCheckResourceAttrSet(dataSourceNameSecurityUserToken, "enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameSecurityUserToken, "expiration_days"),
					resource.TestCheckResourceAttrSet(dataSourceNameSecurityUserToken, "expiration_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameSecurityUserToken, "protect_content"),
				),
			},
		},
	})
}

func getSecurityUserTokenDataSourceConfig() string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
data "%s" "test" {
}
`, dataSourceTypeSecurityUserToken)
}
