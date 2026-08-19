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

package pub_test

import (
	"fmt"
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypePubGroup  = "sonatyperepo_repository_pub_group"
	resourceTypePubHosted = "sonatyperepo_repository_pub_hosted"
	resourceTypePubProxy  = "sonatyperepo_repository_pub_proxy"
)

var (
	resourcePubGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypePubGroup)
	resourcePubHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypePubHosted)
	resourcePubProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypePubProxy)
)

// Requires NXRM 3.92.0+
func skipIfNxrmTooOldForPub(t *testing.T) {
	testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
		Major: 3,
		Minor: 0,
		Patch: 0,
	}, &common.SystemVersion{
		Major: 3,
		Minor: 91,
		Patch: 99,
	})
}

func TestAccRepositoryPubResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForPub(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pub-group-repo-%s"
  online = true
  storage = {
	  blob_store_name = "default"
	  strict_content_type_validation = true
  }
  group = {
	  member_names = []
  }
}
`, resourceTypePubGroup, randomString),
				ExpectError: regexp.MustCompile("Attribute group.member_names list must contain at least 1 elements"),
			},
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pub-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://pub.dev/"
    content_max_age = 1441
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
}

resource "%s" "repo" {
  name = "pub-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
	  member_names = ["pub-proxy-repo-%s"]
  }

  depends_on = [
	  %s.repo
  ]
}
`, resourceTypePubProxy, randomString, resourceTypePubGroup, randomString, randomString, resourceTypePubProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Proxy
					resource.TestCheckResourceAttr(resourcePubProxyName, repotest.RES_ATTR_NAME, fmt.Sprintf("pub-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePubProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourcePubProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourcePubProxyName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePubProxyName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "proxy.remote_url", "https://pub.dev/"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.connection.enable_circular_redirects", "false"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.connection.enable_cookies", "true"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.connection.use_trust_store", "true"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.connection.retries", "9"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.connection.timeout", "999"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.connection.user_agent_suffix", "terraform"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.authentication.username", "user"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.authentication.password", "pass"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.authentication.preemptive", "true"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "http_client.authentication.type", "username"),
					resource.TestCheckNoResourceAttr(resourcePubProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourcePubProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourcePubProxyName, "replication.asset_path_regex"),

					// Verify Group
					resource.TestCheckResourceAttr(resourcePubGroupName, repotest.RES_ATTR_NAME, fmt.Sprintf("pub-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePubGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourcePubGroupName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourcePubGroupName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePubGroupName, "group.member_names.#", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryPubGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("pub-group-import-%s", randomString)
	memberName := fmt.Sprintf("pub-proxy-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForPub(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "member" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://pub.dev/"
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
}

resource "%s" "repo" {
  name = "%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["%s"]
  }
  depends_on = [%s.member]
}
`, resourceTypePubProxy, memberName, resourceTypePubGroup, repoName, memberName, resourceTypePubProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePubGroupName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourcePubGroupName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourcePubGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryPubProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("pub-proxy-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForPub(t) },
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
    remote_url = "https://pub.dev/"
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
}
`, resourceTypePubProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePubProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourcePubProxyName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourcePubProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryPubHostedResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForPub(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "pub-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
    write_policy = "ALLOW"
  }
  component = {
    proprietary_components = false
  }
}
`, resourceTypePubHosted, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePubHostedName, repotest.RES_ATTR_NAME, fmt.Sprintf("pub-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourcePubHostedName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourcePubHostedName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourcePubHostedName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourcePubHostedName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourcePubHostedName, "storage.write_policy", "ALLOW"),
					resource.TestCheckResourceAttr(resourcePubHostedName, "component.proprietary_components", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryPubHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("pub-hosted-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForPub(t) },
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
    write_policy = "ALLOW"
  }
}
`, resourceTypePubHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePubHostedName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourcePubHostedName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourcePubHostedName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
