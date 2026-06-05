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

package huggingface_test

import (
	"fmt"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeHuggingfaceProxy = "sonatyperepo_repository_huggingface_proxy"
)

var (
	resourceHuggingfaceProxyName = fmt.Sprintf(utils_test.RES_NAME_FORMAT, resourceTypeHuggingfaceProxy)
)

// TestAccRepositoryHuggingfaceProxyWithBearerTokenAuth tests Issue #413
// This test verifies that bearer_token authentication works correctly for HuggingFace proxy repositories.
// The Nexus API does not return the bearer_token value after create/update, so the provider must
// copy the value from the plan to the state to avoid "inconsistent values for sensitive attribute" errors.
func TestAccRepositoryHuggingfaceProxyWithBearerTokenAuth(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("huggingface-proxy-bearer-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create repository without authentication
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
    remote_url = "https://huggingface.co"
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
`, resourceTypeHuggingfaceProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceHuggingfaceProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_PROXY_REMOTE_URL, "https://huggingface.co"),
					resource.TestCheckNoResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
				),
			},
			// Step 2: Update to add bearer token authentication - this is where Issue #413 manifests
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
    remote_url = "https://huggingface.co"
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
    authentication = {
      type = "bearerToken"
      bearer_token = "test-bearer-token-value"
    }
  }
}
`, resourceTypeHuggingfaceProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, common.HTTP_AUTH_TYPE_BEARER_TOKEN),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_BEARER_TOKEN, "test-bearer-token-value"),
				),
			},
			// Step 3: Update bearer token value - verify it persists correctly
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
    remote_url = "https://huggingface.co"
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
    authentication = {
      type = "bearerToken"
      bearer_token = "updated-bearer-token-value"
    }
  }
}
`, resourceTypeHuggingfaceProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_BEARER_TOKEN, "updated-bearer-token-value"),
				),
			},
			// Step 4: Remove authentication - verify clean removal
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
    remote_url = "https://huggingface.co"
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
`, resourceTypeHuggingfaceProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckNoResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION),
				),
			},
		},
	})
}

// TestAccRepositoryHuggingfaceProxyWithBearerTokenAuthFromCreate tests creating a repository
// with bearer token authentication from the start (not adding it later).
func TestAccRepositoryHuggingfaceProxyWithBearerTokenAuthFromCreate(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("huggingface-proxy-bearer-create-%s", randomString)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with bearer token authentication from the start
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
    remote_url = "https://huggingface.co"
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
    authentication = {
      type = "bearerToken"
      bearer_token = "initial-bearer-token"
    }
  }
}
`, resourceTypeHuggingfaceProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_ONLINE, "true"),
					resource.TestCheckResourceAttrSet(resourceHuggingfaceProxyName, repotest.RES_ATTR_URL),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_PROXY_REMOTE_URL, "https://huggingface.co"),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_TYPE, common.HTTP_AUTH_TYPE_BEARER_TOKEN),
					resource.TestCheckResourceAttr(resourceHuggingfaceProxyName, repotest.RES_ATTR_HTTP_CLIENT_AUTHENTICATION_BEARER_TOKEN, "initial-bearer-token"),
				),
			},
		},
	})
}
