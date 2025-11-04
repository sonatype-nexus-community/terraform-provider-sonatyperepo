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

package task_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTaskRepairRebuildBrowseNodesResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_task_repair_create_browse_nodes"
	resourceName := fmt.Sprintf(resourceNameF, resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_repository_npm_hosted" "repo" {
   name = "npm-hosted-repo-test-%s"
   online = false
   component = {
     proprietary_components = true
   }
   storage = {
     blob_store_name = "default"
     strict_content_type_validation = true
 	  write_policy = "ALLOW_ONCE"
   }
}
resource "%s" "test_task" {
  name = "test-repair-browse-nodes-%s"
  enabled = true
  alert_email = ""
  notification_condition = "FAILURE"
  frequency = {
    schedule = "manual"
  }
  properties = {
    repository_name = sonatyperepo_repository_npm_hosted.repo.name
  }
}
`, randomString, resourceType, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-repair-browse-nodes-%s", randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_email", ""),
					resource.TestCheckResourceAttr(resourceName, "notification_condition", common.NOTIFICATION_CONDITION_FAILURE),
					resource.TestCheckResourceAttr(resourceName, fieldFrequencySchedule, common.FREQUENCY_SCHEDULE_MANUAL),
					resource.TestCheckResourceAttr(resourceName, "properties.repository_name", fmt.Sprintf("npm-hosted-repo-test-%s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
