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

package oci_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeOciGroup  = "sonatyperepo_repository_oci_group"
	resourceTypeOciHosted = "sonatyperepo_repository_oci_hosted"
	resourceTypeOciProxy  = "sonatyperepo_repository_oci_proxy"
)

var (
	resourceOciGroupName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeOciGroup)
	resourceOciHostedName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeOciHosted)
	resourceOciProxyName  = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeOciProxy)
)

// OCI repositories require NXRM 3.94.0+
func skipIfNxrmTooOldForOci(t *testing.T) {
	testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
		Major: 3,
		Minor: 0,
		Patch: 0,
	}, &common.SystemVersion{
		Major: 3,
		Minor: 93,
		Patch: 99,
	})
}

func TestAccRepositoryOciResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "oci-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://registry-1.docker.io"
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
  oci = {
    force_basic_auth = false
    v1_enabled = false
  }
  oci_proxy = {
    index_type = "REGISTRY"
    index_url = "https://index.docker.io"
  }
}

resource "%s" "repo" {
  name = "oci-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["oci-proxy-repo-%s"]
  }
  oci = {
    force_basic_auth = true
    v1_enabled = false
  }

  depends_on = [
    %s.repo
  ]
}
`, resourceTypeOciProxy, randomString, resourceTypeOciGroup, randomString, randomString, resourceTypeOciProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Proxy
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_NAME, fmt.Sprintf("oci-proxy-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceOciProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "proxy.remote_url", "https://registry-1.docker.io"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "proxy.content_max_age", "1441"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "proxy.metadata_max_age", "1440"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "negative_cache.enabled", "true"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "negative_cache.time_to_live", "1440"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "http_client.blocked", "false"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "http_client.auto_block", "true"),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_OCI_FORCE_BASIC_AUTH, "false"),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_OCI_V1_ENABLED, "false"),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_OCI_PROXY_INDEX_TYPE, "REGISTRY"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "oci_proxy.index_url", "https://index.docker.io"),
					// cosign block omitted from config entirely -- verifies the Computed default
					// (NXRM itself defaults to this shape; see format/oci.go's ociCosignSchemaAttributes)
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_OCI_COSIGN_ENFORCEMENT, "NONE"),
					resource.TestCheckNoResourceAttr(resourceOciProxyName, "cosign.identity_regex"),
					resource.TestCheckNoResourceAttr(resourceOciProxyName, "cosign.issuer_regex"),
					resource.TestCheckNoResourceAttr(resourceOciProxyName, "routing_rule"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "replication.preemptive_pull_enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceOciProxyName, "replication.asset_path_regex"),

					// Verify Group
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_NAME, fmt.Sprintf("oci-group-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceOciGroupName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceOciGroupName, "group.member_names.#", "1"),
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_OCI_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_OCI_COSIGN_ENFORCEMENT, "NONE"),
				),
			},
			// Update - explicitly set Cosign (still NONE -- KEYLESS is rejected server-side as of
			// NXRM 3.95.0-07: "Cosign keyless enforcement is not yet available; choose NONE until a
			// real verifier is shipped" -- but exercises identity_regex/issuer_regex round-tripping)
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "oci-proxy-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  proxy = {
    remote_url = "https://registry-1.docker.io"
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
  }
  oci = {
    force_basic_auth = false
    v1_enabled = false
  }
  oci_proxy = {
    index_type = "REGISTRY"
  }
  cosign = {
    enforcement    = "NONE"
    identity_regex = "^https://github.com/example/.*$"
    issuer_regex   = "^https://token.actions.githubusercontent.com$"
  }
}

resource "%s" "repo" {
  name = "oci-group-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
  }
  group = {
    member_names = ["oci-proxy-repo-%s"]
  }
  oci = {
    force_basic_auth = true
    v1_enabled = false
  }

  depends_on = [
    %s.repo
  ]
}
`, resourceTypeOciProxy, randomString, resourceTypeOciGroup, randomString, randomString, resourceTypeOciProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_OCI_COSIGN_ENFORCEMENT, "NONE"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "cosign.identity_regex", "^https://github.com/example/.*$"),
					resource.TestCheckResourceAttr(resourceOciProxyName, "cosign.issuer_regex", "^https://token.actions.githubusercontent.com$"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryOciGroupImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("oci-group-import-%s", randomString)
	memberName := fmt.Sprintf("oci-proxy-member-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
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
  oci = {
    force_basic_auth = false
    v1_enabled = false
  }
  oci_proxy = {
    index_type = "REGISTRY"
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
  oci = {
    force_basic_auth = true
    v1_enabled = false
  }
  depends_on = [%s.member]
}
`, resourceTypeOciProxy, memberName, resourceTypeOciGroup, repoName, memberName, resourceTypeOciProxy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceOciGroupName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceOciGroupName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryOciProxyImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("oci-proxy-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
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
  oci = {
    force_basic_auth = false
    v1_enabled = false
  }
  oci_proxy = {
    index_type = "REGISTRY"
  }
}
`, resourceTypeOciProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceOciProxyName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}

func TestAccRepositoryOciHostedResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "repo" {
  name = "oci-hosted-repo-%s"
  online = true
  storage = {
    blob_store_name = "default"
    strict_content_type_validation = true
    write_policy = "ALLOW"
  }
  component = {
    proprietary_components = false
  }
  oci = {
    force_basic_auth = true
    v1_enabled = false
  }
  cosign = {
    // NOTE: KEYLESS enforcement is documented by NXRM's OpenAPI schema but rejected server-side
    // as of NXRM 3.95.0-07 ("Cosign keyless enforcement is not yet available; choose NONE until
    // a real verifier is shipped") -- NONE is exercised here to still cover the cosign wiring.
    enforcement    = "NONE"
    identity_regex = "^https://github.com/example/.*$"
  }
}
`, resourceTypeOciHosted, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_NAME, fmt.Sprintf("oci-hosted-repo-%s", randomString)),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceOciHostedName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_STORAGE_BLOB_STORE_NAME, common.DEFAULT_BLOB_STORE_NAME),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_STORAGE_STRICT_CONTENT_TYPE_VALIDATION, "true"),
					resource.TestCheckResourceAttr(resourceOciHostedName, "storage.write_policy", "ALLOW"),
					resource.TestCheckResourceAttr(resourceOciHostedName, "component.proprietary_components", "false"),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_OCI_FORCE_BASIC_AUTH, "true"),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_OCI_V1_ENABLED, "false"),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_OCI_COSIGN_ENFORCEMENT, "NONE"),
					resource.TestCheckResourceAttr(resourceOciHostedName, "cosign.identity_regex", "^https://github.com/example/.*$"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRepositoryOciHostedImport(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("oci-hosted-import-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
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
  oci = {
    force_basic_auth = true
    v1_enabled = false
  }
}
`, resourceTypeOciHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Import and verify no changes
			{
				ResourceName:                         resourceOciHostedName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        repoName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
		},
	})
}
