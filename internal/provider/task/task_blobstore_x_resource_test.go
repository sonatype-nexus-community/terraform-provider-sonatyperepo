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

func TestAccTaskBlobstoreCompactResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeTaskBlobstoreCompact := "sonatyperepo_task_blobstore_compact"
	resourceNameTaskBlobstoreCompact := fmt.Sprintf("%s.test_task", resourceTypeTaskBlobstoreCompact)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test_task" {
  name = "test-blobstore-compact-%s"
  enabled = true
  alert_email = ""
  notification_condition = "FAILURE"
  frequency = {
    schedule = "manual"
  }
  properties = {
    blob_store_name = "default"
    blobs_older_than = 8
  }
}
`, resourceTypeTaskBlobstoreCompact, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameTaskBlobstoreCompact, "id"),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "name", fmt.Sprintf("test-blobstore-compact-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "alert_email", ""),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "notification_condition", common.NOTIFICATION_CONDITION_FAILURE),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "frequency.schedule", common.FREQUENCY_SCHEDULE_MANUAL),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "properties.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "properties.blobs_older_than", "8"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
