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
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRepositoryOciHostedResource_Issue499 reproduces
// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/499
//
// Explicitly configuring an empty `cleanup` object (no cleanup policies attached) planned a
// non-null `cleanup` object, but the provider mapped the API's empty policy list back to a null
// `cleanup` after apply - failing with "Provider produced inconsistent result after apply".
func TestAccRepositoryOciHostedResource_Issue499(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("oci-hosted-issue499-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
    force_basic_auth = false
    v1_enabled = false
  }
  cleanup = {}
}
`, resourceTypeOciHosted, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciHostedName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttrSet(resourceOciHostedName, "cleanup.%"),
				),
			},
		},
	})
}

// TestAccRepositoryOciProxyResource_Issue499 is the proxy equivalent of
// TestAccRepositoryOciHostedResource_Issue499 - and uses an empty (non-null) `policy_names` set,
// as produced by the reporter's `for` expression over an empty list.
func TestAccRepositoryOciProxyResource_Issue499(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("oci-proxy-issue499-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
  cleanup = {
    policy_names = []
  }
}
`, resourceTypeOciProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceOciProxyName, "cleanup.policy_names.#", "0"),
				),
			},
		},
	})
}
