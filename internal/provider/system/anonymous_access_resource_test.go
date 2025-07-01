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
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSysAnonymousAccess = "sonatyperepo_system_anonymous_access"
	resourceNameSysAnonymousAccess = "sonatyperepo_system_anonymous_access.aa"
)

func TestAccSystemAnonymousAccessResource(t *testing.T) {
	if os.Getenv("TF_ACC_SINGLE_HIT") == "1" {
		randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: getSystemAnonymousAccessResourceConfig(randomString),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "realm_name", common.DEFAULT_REALM_NAME),
						resource.TestCheckResourceAttr(resourceNameSysAnonymousAccess, "user_id", fmt.Sprintf("anonymous-%s", randomString)),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func getSystemAnonymousAccessResourceConfig(randomString string) string {
	return fmt.Sprintf(utils.ProviderConfig+`
resource "%s" "aa" {
  enabled = true
  realm_name = "%s"
  user_id = "anonymous-%s"
}
`, resourceTypeSysAnonymousAccess, common.DEFAULT_REALM_NAME, randomString)
}
