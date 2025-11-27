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
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositorDockerResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeGroup := "sonatyperepo_repository_docker_group"
	resourceTypeHosted := "sonatyperepo_repository_docker_hosted"
	resourceTypeProxy := "sonatyperepo_repository_docker_proxy"
	resourceGroupName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeGroup)
	resourceHostedName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHosted)
	resourceProxyName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeProxy)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = []
  }
  docker = {
    force_basic_auth = false
    v1_enabled = false
  }
}
`, resourceTypeGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-hosted-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	latest_policy = true
	strict_content_type_validation = true
	write_policy = "ALLOW_ONCE"
  }
  docker = {
    force_basic_auth = true
    v1_enabled = true
  }
}

resource "%s" "repo" {
  name = "docker-proxy-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://registry-1.docker.io"
    content_max_age = 1442
    metadata_max_age = 1400
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
	connection = {
		enable_cookies = true
		retries = 9
		timeout = 999
		use_trust_store = true
		user_agent_suffix = "terraform"
	}
	authentication = {
		username = "user"
		password = "pass"
		preemptive = true
		type = "username"
	}
  }
  docker = {
    force_basic_auth = true
    v1_enabled = true
  }
  docker_proxy = {  }
}

resource "%s" "repo" {
  name = "docker-group-repo-%s"
  online = true
  storage = {
	blob_store_name = "default"
	strict_content_type_validation = true
  }
  group = {
	member_names = ["docker-proxy-repo-%s"]
  }
  docker = {
    force_basic_auth = false
    v1_enabled = false
  }

  depends_on = [
	%s.repo
  ]
}
`, resourceTypeHosted, randomString, resourceTypeProxy, randomString, resourceTypeGroup, randomString, randomString, resourceTypeProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceHostedName, "name", fmt.Sprintf("docker-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceHostedName, "url"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.latest_policy", "true"),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_DOCKER_V1_ENABLED, "true"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceProxyName, "name", fmt.Sprintf("docker-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceProxyName, "url"),
					resource.TestCheckResourceAttr(resourceProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.remote_url", "https://registry-1.docker.io"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourceProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceProxyName, "replication.asset_path_regex"),
					resource.TestCheckResourceAttr(resourceProxyName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceProxyName, RES_ATTR_DOCKER_V1_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceProxyName, "docker_proxy.cache_foreign_layers", "false"),
					resource.TestCheckResourceAttr(resourceProxyName, "docker_proxy.foreign_layer_url_whitelist.#", "0"),
					resource.TestCheckResourceAttr(resourceProxyName, "docker_proxy.index_type", common.DOCKER_PROXY_INDEX_TYPE_REGISTRY),

					// Verify Group
					resource.TestCheckResourceAttr(resourceGroupName, "name", fmt.Sprintf("docker-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceGroupName, "url"),
					resource.TestCheckResourceAttr(resourceGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceGroupName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
					resource.TestCheckResourceAttr(resourceGroupName, RES_ATTR_DOCKER_V1_ENABLED, "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositorDockerPathEnabledResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceTypeHosted := "sonatyperepo_repository_docker_hosted"
	resourceHostedName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHosted)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// This is new functionality in NXRM 3.83.0+
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 0,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 82,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
	resource "%s" "repo" {
	name = "docker-hosted-repo-%s"
	online = true
	storage = {
		blob_store_name = "default"
		strict_content_type_validation = true
		write_policy = "ALLOW_ONCE"
	}
	docker = {
		force_basic_auth = true
		path_enabled = true
		v1_enabled = true
	}
	}`, resourceTypeHosted, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceHostedName, "name", fmt.Sprintf("docker-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceHostedName, "url"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_DOCKER_PATH_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceHostedName, RES_ATTR_DOCKER_V1_ENABLED, "true"),
				),
				// Delete testing automatically occurs in TestCase
			},
		},
	})
}

func TestAccRepositoryDockerHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_docker_hosted"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("docker-hosted-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  component = {
    proprietary_components = true
  }
  storage = {
    blob_store_name = "default"
    latest_policy                  = true
    strict_content_type_validation = true
	write_policy = "ALLOW"
  }
  docker = {
    force_basic_auth = false
    v1_enabled = false
  }
}
`, resourceType, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", repoName),
					resource.TestCheckResourceAttr(resourceName, "online", "true"),
					resource.TestCheckResourceAttr(resourceName, "component.proprietary_components", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "storage.latest_policy", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage.write_policy", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "docker.force_basic_auth", "false"),
					resource.TestCheckResourceAttr(resourceName, "docker.v1_enabled", "false"),
				),
			},
			// Import and verify no changes
			{
				ResourceName: resourceName,
				ImportState:  true,
				// Cannot test for valid import state due to API not returning `latest_policy` when reading
				// Docker Registries
				//
				// ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryDockerProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := "sonatyperepo_repository_docker_proxy"
	resourceName := fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceType)
	repoName := fmt.Sprintf("docker-proxy-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://registry-1.docker.io"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = 1440
  }
  http_client = {
    blocked = false
    auto_block = true
  }
  docker = {
    force_basic_auth = true
    v1_enabled = false
  }
  docker_proxy = {}
}
`, resourceType, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", repoName),
					resource.TestCheckResourceAttr(resourceName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
