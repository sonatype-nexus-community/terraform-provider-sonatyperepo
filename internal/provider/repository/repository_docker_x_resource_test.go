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

const (
	resourceTypeDockerGroup  = "sonatyperepo_repository_docker_group"
	resourceTypeDockerHosted = "sonatyperepo_repository_docker_hosted"
	resourceTypeDockerProxy  = "sonatyperepo_repository_docker_proxy"
)

var (
	resourceDockerGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeDockerGroup)
	resourceDockerHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeDockerHosted)
	resourceDockerProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeDockerProxy)
)

func TestAccRepositorDockerResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
`, resourceTypeDockerGroup, randomString),
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
`, resourceTypeDockerHosted, randomString, resourceTypeDockerProxy, randomString, resourceTypeDockerGroup, randomString, randomString, resourceTypeDockerProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceDockerHostedName, "name", fmt.Sprintf("docker-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceDockerHostedName, "url"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.latest_policy", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceDockerHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_DOCKER_V1_ENABLED, "true"),

					// Verify Proxy
					resource.TestCheckResourceAttr(resourceDockerProxyName, "name", fmt.Sprintf("docker-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceDockerProxyName, "url"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "proxy.remote_url", "https://registry-1.docker.io"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "proxy.content_max_age", "1442"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "proxy.metadata_max_age", "1400"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "http_client.connection.user_agent_suffix", "terraform"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceDockerGroupName, "name", fmt.Sprintf("docker-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceDockerGroupName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceDockerGroupName, "url"),
					resource.TestCheckResourceAttr(resourceDockerGroupName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceDockerGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceDockerGroupName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "false"),
					resource.TestCheckResourceAttr(resourceDockerGroupName, RES_ATTR_DOCKER_V1_ENABLED, "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositorDockerPathEnabledResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
	}`, resourceTypeDockerHosted, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Hosted
					resource.TestCheckResourceAttr(resourceDockerHostedName, "name", fmt.Sprintf("docker-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "online", "true"),
					resource.TestCheckResourceAttrSet(resourceDockerHostedName, "url"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.write_policy", common.WRITE_POLICY_ALLOW_ONCE),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "component.proprietary_components", "false"),
					resource.TestCheckNoResourceAttr(resourceDockerHostedName, "cleanup"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_DOCKER_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_DOCKER_PATH_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, RES_ATTR_DOCKER_V1_ENABLED, "true"),
				),
				// Delete testing automatically occurs in TestCase
			},
		},
	})
}

func TestAccRepositoryDockerHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
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
`, resourceTypeDockerHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceDockerHostedName, "name", repoName),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "online", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "component.proprietary_components", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.blob_store_name", "default"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.latest_policy", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.strict_content_type_validation", "true"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "storage.write_policy", "ALLOW"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "docker.force_basic_auth", "false"),
					resource.TestCheckResourceAttr(resourceDockerHostedName, "docker.v1_enabled", "false"),
				),
			},
			// Import and verify no changes
			{
				ResourceName: resourceDockerHostedName,
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
`, resourceTypeDockerProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceDockerProxyName, "name", repoName),
					resource.TestCheckResourceAttr(resourceDockerProxyName, "online", "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceDockerProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
func TestAccRepositoryDockerProxyInvalidRemoteUrl(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid remote URL (missing protocol)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "invalid-url-without-protocol"
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
    force_basic_auth = false
    enable_v1 = false
  }
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be a valid URL|must be a valid HTTP URL"),
			},
		},
	})
}

func TestAccRepositoryDockerHostedInvalidBlobStore(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid blob store name (non-existent)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "non-existent-blob-store"
    strict_content_type_validation = true
  }
  docker = {
    v1_enabled = true
  }
}
`, resourceTypeDockerHosted, randomString),
				ExpectError: regexp.MustCompile("Blob store.*not found|Blob store.*does not exist"),
			},
		},
	})
}

func TestAccRepositoryDockerHostedMissingStorage(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing storage block (required field)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-hosted-repo-%s"
  online = true
  # Missing storage block
}
`, resourceTypeDockerHosted, randomString),
				ExpectError: regexp.MustCompile("Attribute storage is required"),
			},
		},
	})
}

func TestAccRepositoryDockerProxyInvalidTimeoutTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too large, max is 3600)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-timeout-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      timeout = 3601
    }
  }
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 3600"),
			},
		},
	})
}

func TestAccRepositoryDockerProxyInvalidTimeoutTooSmall(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid timeout (too small, min is 1)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-timeout-small-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      timeout = 0
    }
  }
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 1"),
			},
		},
	})
}

func TestAccRepositoryDockerProxyInvalidRetriesTooLarge(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (too large, max is 10)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-retries-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      retries = 11
    }
  }
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be less than or equal to 10"),
			},
		},
	})
}

func TestAccRepositoryDockerProxyInvalidRetriesNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid retries (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-retries-neg-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
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
    connection = {
      retries = -1
    }
  }
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be between|must be greater than or equal to 0"),
			},
		},
	})
}

func TestAccRepositoryDockerProxyInvalidMaxAgeNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid content_max_age (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-maxage-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
    content_max_age = -1
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
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}

func TestAccRepositoryDockerProxyInvalidTimeToLiveNegative(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid time_to_live (negative)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "docker-proxy-repo-ttl-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://repo.example.com"
    content_max_age = 1440
    metadata_max_age = 1440
  }
  negative_cache = {
    enabled = true
    time_to_live = -1
  }
  http_client = {
    blocked = false
    auto_block = true
  }
}
`, resourceTypeDockerProxy, randomString),
				ExpectError: regexp.MustCompile("must be greater than or equal to|cannot be negative"),
			},
		},
	})
}
