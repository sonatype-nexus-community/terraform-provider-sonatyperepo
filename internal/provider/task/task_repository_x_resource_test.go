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

func TestAccTaskRepositoryDockerGcResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeTaskBlobstoreCompact := "sonatyperepo_task_repository_docker_gc"
	resourceNameTaskBlobstoreCompact := fmt.Sprintf("%s.test_task", resourceTypeTaskBlobstoreCompact)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "sonatyperepo_repository_docker_hosted" "repo" {
  name = "docker-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  docker = {
    force_basic_auth = true
    v1_enabled = true
  }
}

resource "%s" "test_task" {
  name = "test-repository-docker-gc-%s"
  enabled = true
  alert_email = ""
  notification_condition = "FAILURE"
  frequency = {
    schedule = "manual"
  }
  properties = {
    repository_name = sonatyperepo_repository_docker_hosted.repo.name
  }
}
`, randomString, resourceTypeTaskBlobstoreCompact, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameTaskBlobstoreCompact, "id"),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "name", fmt.Sprintf("test-repository-docker-gc-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "alert_email", ""),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "notification_condition", common.NOTIFICATION_CONDITION_FAILURE),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "frequency.schedule", common.FREQUENCY_SCHEDULE_MANUAL),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "properties.deploy_offset", fmt.Sprintf("%d", common.TASK_REPOSITORY_DOCKER_GC_DEFAULT_DEPLOY_OFFSET)),
					resource.TestCheckResourceAttr(resourceNameTaskBlobstoreCompact, "properties.repository_name", fmt.Sprintf("docker-hosted-repo-%s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTaskRepositoryDockerUploadPurgeResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_task_repository_docker_upload_purge"
	resourceName := fmt.Sprintf("%s.test_task", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test_task" {
  name = "test-repository-docker-upload-purge-%s"
  enabled = true
  alert_email = ""
  notification_condition = "FAILURE"
  frequency = {
    schedule = "manual"
  }
  properties = {}
}
`, resourceType, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-repository-docker-upload-purge-%s", randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_email", ""),
					resource.TestCheckResourceAttr(resourceName, "notification_condition", common.NOTIFICATION_CONDITION_FAILURE),
					resource.TestCheckResourceAttr(resourceName, "frequency.schedule", common.FREQUENCY_SCHEDULE_MANUAL),
					resource.TestCheckResourceAttr(resourceName, "properties.age", fmt.Sprintf("%d", common.TASK_REPOSITORY_DOCKER_UPLOAD_PURGE_DEFAULT_AGE)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
