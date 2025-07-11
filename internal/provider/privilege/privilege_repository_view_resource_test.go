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

package privilege

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"
	"terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceTypePrivilegeRepoView = "sonatyperepo_privilege_repository_view"
	resourceNamePrivilegeRepoView = fmt.Sprintf("%s.p", resourceTypePrivilegeRepoView)
)

func TestAccPrivilegeRepositoryViewResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils.ProviderConfig+`
resource "%s" "p" {
	name = "test-priv-repo-view-%s"
	description = "a description"
	actions = [
    	"BROWSE",
		"READ"
  	]
	format = "npm"
	repository = "npm-all"
}`, resourceTypePrivilegeRepoView, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "name", fmt.Sprintf("test-priv-repo-view-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "description", "a description"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "read_only", "false"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "type", privilege_type.TypeRepositoryView.String()),
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "format", "npm"),
					resource.TestCheckResourceAttr(resourceNamePrivilegeRepoView, "repository", "npm-all"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
