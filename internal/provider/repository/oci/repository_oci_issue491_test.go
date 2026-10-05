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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRepositoryOciProxyResource_Issue491 reproduces
// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/491
//
// An OCI proxy repository is created via Terraform with http_client.authentication left
// unconfigured (null) and `lifecycle { ignore_changes = [http_client.authentication] }` set.
// Upstream authentication is then configured manually (out of band, simulating the Nexus UI)
// directly against the repository via the raw NXRM REST API - so the provider's state ends up
// with `type`/`username` (from the next Read()) but no `password` (NXRM never returns it, and
// there is no plan value to restore it from, since Terraform never managed it in the first
// place). Applying an unrelated Terraform change (a proxy attribute) previously sent that
// half-populated authentication object straight through, which NXRM's own validation rejects
// with `password must not be null` - surfaced by the provider as a misleading
// "Repository did not exist to update" error.
//
// Fixed in internal/provider/model/repository_proxy.go: MapToApiHttpClientAttributes/
// MapToApiHttpClientAttributesWithPreemptiveAuth now omit `authentication` from the outbound
// request entirely whenever it's missing the secret NXRM requires for its type (see
// hasRequiredSecret), and MapMissingApiFieldsFromPlan mirrors the plan's authentication verbatim
// into state in that case, since NXRM's PUT is a full replace (no partial-update endpoint
// exists) and genuinely clears the field when it's omitted - so state must reflect what was
// actually sent (nothing changed) to stay consistent with the `ignore_changes`-frozen plan.
// This fix lives in code shared by every proxy repository format with
// http_client.authentication, not just OCI - see CHANGELOG.md.
func TestAccRepositoryOciProxyResource_Issue491(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	repoName := fmt.Sprintf("oci-proxy-issue491-%s", randomString)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { skipIfNxrmTooOldForOci(t) },
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create OCI proxy repository with no http_client.authentication configured
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

  lifecycle {
    ignore_changes = [http_client.authentication]
  }
}
`, resourceTypeOciProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_ONLINE, "true"),
				),
			},
			// Step 2: Manually configure upstream authentication out of band (simulating the
			// Nexus UI), then apply an unrelated change (proxy.content_max_age) via Terraform.
			{
				PreConfig: func() {
					manuallyConfigureOciProxyAuthentication(t, repoName)
				},
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

  lifecycle {
    ignore_changes = [http_client.authentication]
  }
}
`, resourceTypeOciProxy, repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceOciProxyName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceOciProxyName, "proxy.content_max_age", "1441"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// manuallyConfigureOciProxyAuthentication sets http_client.authentication on the named OCI proxy
// repository directly via the raw NXRM REST API, bypassing both the provider and the generated Go
// client entirely (OCI has no equivalent methods in the v3/"v382" client generation this provider
// otherwise uses) - simulating a user manually configuring upstream authentication credentials in
// the Nexus UI that Terraform does not manage.
func manuallyConfigureOciProxyAuthentication(t *testing.T, repositoryName string) {
	t.Helper()

	baseUrl := strings.TrimRight(os.Getenv("NXRM_SERVER_URL"), "/")
	username := os.Getenv("NXRM_SERVER_USERNAME")
	password := os.Getenv("NXRM_SERVER_PASSWORD")
	repoUrl := fmt.Sprintf("%s/service/rest/v1/repositories/oci/proxy/%s", baseUrl, repositoryName)

	current := getRawRepository(t, repoUrl, username, password)

	httpClient, ok := current["httpClient"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected httpClient object in repository %q, got: %#v", repositoryName, current["httpClient"])
	}
	httpClient["authentication"] = map[string]interface{}{
		"type":     "username",
		"username": "manually-configured-user",
		"password": "manually-configured-password",
	}

	putRawRepository(t, repoUrl, username, password, current)
}

func getRawRepository(t *testing.T, url, username, password string) map[string]interface{} {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to build GET request for %q: %v", url, err)
	}
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to GET %q: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body from GET %q: %v", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status GET %q: %d - %s", url, resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to unmarshal response body from GET %q: %v", url, err)
	}
	return result
}

func putRawRepository(t *testing.T, url, username, password string, body map[string]interface{}) {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body for PUT %q: %v", url, err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to build PUT request for %q: %v", url, err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to PUT %q: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body from PUT %q: %v", url, err)
	}
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status PUT %q: %d - %s", url, resp.StatusCode, string(respBody))
	}
}
