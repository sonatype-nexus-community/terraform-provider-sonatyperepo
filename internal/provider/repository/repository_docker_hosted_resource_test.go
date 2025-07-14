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

package repository_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameDockerHosted = "sonatyperepo_repository_docker_hosted.repo"
	resourceTypeDockerHosted = "sonatyperepo_repository_docker_hosted"
)

func TestAccRepositoryDockerHostedResourceNoReplication(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getRepositoryDockerHostedResourceConfig(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "name", fmt.Sprintf("docker-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "online", "false"),
					resource.TestCheckResourceAttrSet(resourceNameDockerHosted, "url"),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "component.proprietary_components", "true"),
					resource.TestCheckNoResourceAttr(resourceNameDockerHosted, "cleanup"),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "docker.force_basic_auth", "true"),
					resource.TestCheckResourceAttr(resourceNameDockerHosted, "docker.v1_enabled", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func getRepositoryDockerHostedResourceConfig(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-hosted-repo-%s"
  online = false
  component = {
    proprietary_components = true
  }
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
	  write_policy = "ALLOW_ONCE"
  }
  docker = {
    force_basic_auth = true
    v1_enabled = false
  }
}
`, resourceTypeDockerHosted, randomString)
}
