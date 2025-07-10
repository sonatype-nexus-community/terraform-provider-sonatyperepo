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

package privilege_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"
	"terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceTypePrivilegeApplication = "sonatyperepo_privilege_application"
	resourceNamePrivilegeApplication = fmt.Sprintf("%s.p", resourceTypePrivilegeApplication)
)

func TestAccPrivilegeApplicationResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-app-%s"
	description = "some description"
	domain = "rubbish"
	actions = [
    	"ALL"
  	]
}`, resourceTypePrivilegeApplication, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNamePrivilegeApplication, "name", fmt.Sprintf("test-priv-app-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNamePrivilegeApplication, "description", "some description"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeApplication, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeApplication, "type", privilege_type.TypeApplication.String()),
					resource.TestCheckResourceAttr(resourceNamePrivilegeApplication, "domain", "rubbish"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeApplication, "actions.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
