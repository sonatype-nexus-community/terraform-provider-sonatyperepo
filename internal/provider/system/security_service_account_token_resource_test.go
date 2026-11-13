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

package system_test

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

const (
	serviceAccountTokenResourceType = "sonatyperepo_security_service_account_token"
	serviceAccountTokenResourceName = serviceAccountTokenResourceType + ".test"

	// A built-in role that exists on every Nexus Repository, so the tests do not depend on
	// (or interfere with) other tests' roles.
	serviceAccountTokenTestRole = "nx-anonymous"
)

// skipIfServiceAccountTokensUnavailable skips tests that call the service account token API:
// it needs NXRM 3.97.0+ and, per Sonatype's documentation, Nexus Repository Pro (Community
// Edition answers with an error). It is not listed in the feature matrix.
func skipIfServiceAccountTokensUnavailable(t *testing.T) {
	t.Helper()
	testutil.SkipOnCommunityEdition(t, "the service account token API requires Nexus Repository Pro")
	skipIfNxrmBefore397(t)
}

func skipIfNxrmBefore397(t *testing.T) {
	t.Helper()
	testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
		Major: 3,
		Minor: 0,
		Patch: 0,
	}, &common.SystemVersion{
		Major: 3,
		Minor: 96,
		Patch: 99,
	})
}

func serviceAccountTokenConfig(name string, extraArguments ...string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  name    = "%s"
  role_id = "%s"
%s}
`, serviceAccountTokenResourceType, name, serviceAccountTokenTestRole, strings.Join(extraArguments, ""))
}

// TestAccSecurityServiceAccountTokenResourceValidation checks the plan-time validation. It
// never reaches the service account token API, so it needs no licence and runs everywhere.
func TestAccSecurityServiceAccountTokenResourceValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      serviceAccountTokenConfig("has a space"),
				ExpectError: regexp.MustCompile(`must\s+contain\s+only\s+letters`),
			},
			{
				Config:      serviceAccountTokenConfig(strings.Repeat("a", 129)),
				ExpectError: regexp.MustCompile(`string\s+length\s+must\s+be\s+between\s+1\s+and\s+128`),
			},
			{
				Config:      serviceAccountTokenConfig("ok", fmt.Sprintf("  description = %q\n", strings.Repeat("d", 257))),
				ExpectError: regexp.MustCompile(`string\s+length\s+must\s+be\s+at\s+most\s+256`),
			},
			{
				Config:      serviceAccountTokenConfig("ok", "  expiration_days = 0\n"),
				ExpectError: regexp.MustCompile(`value\s+must\s+be\s+at\s+least\s+1`),
			},
			{
				Config: utils_test.ProviderConfig + fmt.Sprintf(`
resource "%s" "test" {
  name = "ok"
}
`, serviceAccountTokenResourceType),
				ExpectError: regexp.MustCompile(`(?s)role_id.*required`),
			},
		},
	})
}

// TestAccSecurityServiceAccountTokenResourceUnsupportedVersion checks the clear error on
// servers that predate service account tokens. It only runs against NXRM before 3.97.0, and
// needs no licence because the provider rejects the request before calling the server.
func TestAccSecurityServiceAccountTokenResourceUnsupportedVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 97,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 99,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			{
				Config:      serviceAccountTokenConfig("tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)),
				ExpectError: regexp.MustCompile(`(?s)3\.97\.0\s+or\s+later`),
			},
		},
	})
}

func TestAccSecurityServiceAccountTokenResource(t *testing.T) {
	name := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	replaced := func() resource.ConfigPlanChecks {
		return resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(serviceAccountTokenResourceName, plancheck.ResourceActionReplace)},
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfServiceAccountTokensUnavailable(t)
		},
		Steps: []resource.TestStep{
			// Create with required arguments only: never expires
			{
				Config: serviceAccountTokenConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(serviceAccountTokenResourceName, "id"),
					resource.TestCheckResourceAttr(serviceAccountTokenResourceName, "name", name),
					resource.TestCheckResourceAttr(serviceAccountTokenResourceName, "role_id", serviceAccountTokenTestRole),
					resource.TestCheckResourceAttrWith(serviceAccountTokenResourceName, "token", func(value string) error {
						if !strings.HasPrefix(value, "sat.") {
							return fmt.Errorf("expected a token starting with %q", "sat.")
						}
						return nil
					}),
					resource.TestCheckNoResourceAttr(serviceAccountTokenResourceName, "description"),
					resource.TestCheckNoResourceAttr(serviceAccountTokenResourceName, "expiration_days"),
					resource.TestCheckNoResourceAttr(serviceAccountTokenResourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(serviceAccountTokenResourceName, "created_at"),
					resource.TestCheckResourceAttrSet(serviceAccountTokenResourceName, "created_by"),
				),
			},
			// Same configuration again: nothing to do (and the token must survive refresh)
			{
				Config:   serviceAccountTokenConfig(name),
				PlanOnly: true,
			},
			// Any change replaces the token
			{
				Config:           serviceAccountTokenConfig(name, "  description = \"managed by terraform\"\n"),
				ConfigPlanChecks: replaced(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(serviceAccountTokenResourceName, "description", "managed by terraform"),
					resource.TestCheckResourceAttrSet(serviceAccountTokenResourceName, "token"),
				),
			},
			// Expiry
			{
				Config:           serviceAccountTokenConfig(name, "  description = \"managed by terraform\"\n", "  expiration_days = 30\n"),
				ConfigPlanChecks: replaced(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(serviceAccountTokenResourceName, "expiration_days", "30"),
					resource.TestCheckResourceAttrSet(serviceAccountTokenResourceName, "expires_at"),
				),
			},
			// Importing is not supported: the token value cannot be read back
			{
				ResourceName:  serviceAccountTokenResourceName,
				ImportState:   true,
				ImportStateId: "anything",
				ExpectError:   regexp.MustCompile(`Import\s+is\s+not\s+supported`),
			},
		},
	})
}

func TestAccSecurityServiceAccountTokenResourceDuplicateName(t *testing.T) {
	name := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfServiceAccountTokensUnavailable(t)
		},
		Steps: []resource.TestStep{
			{
				Config: serviceAccountTokenConfig(name) + fmt.Sprintf(`
resource "%s" "duplicate" {
  name       = %q
  role_id    = %q
  depends_on = [%s.test]
}
`, serviceAccountTokenResourceType, name, serviceAccountTokenTestRole, serviceAccountTokenResourceType),
				ExpectError: regexp.MustCompile(`(?s)already\s+in\s+use`),
			},
		},
	})
}

func TestAccSecurityServiceAccountTokenResourceUnknownRole(t *testing.T) {
	name := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfServiceAccountTokensUnavailable(t)
		},
		Steps: []resource.TestStep{
			{
				Config: utils_test.ProviderConfig + fmt.Sprintf(`
resource "%s" "test" {
  name    = %q
  role_id = "no-such-role-%s"
}
`, serviceAccountTokenResourceType, name, name),
				ExpectError: regexp.MustCompile(`(?s)Unknown\s+role`),
			},
		},
	})
}

// TestAccSecurityServiceAccountTokenResourceDeletedOutsideTerraform checks that a token
// revoked outside Terraform is noticed on refresh (it would be recreated), and that destroying
// a token that is already gone is not an error.
func TestAccSecurityServiceAccountTokenResourceDeletedOutsideTerraform(t *testing.T) {
	name := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	// The id of the token created by step 1, so that step 2 can revoke it.
	var tokenId string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfServiceAccountTokensUnavailable(t)
		},
		Steps: []resource.TestStep{
			{
				Config: serviceAccountTokenConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(serviceAccountTokenResourceName, "id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[serviceAccountTokenResourceName]
						if !ok {
							return fmt.Errorf("resource %s not found in state", serviceAccountTokenResourceName)
						}
						tokenId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					if err := revokeServiceAccountToken(tokenId); err != nil {
						t.Fatalf("Unable to revoke the token outside Terraform: %v", err)
					}
				},
				Config:             serviceAccountTokenConfig(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// revokeServiceAccountToken deletes a token directly through the NXRM API.
func revokeServiceAccountToken(id string) error {
	url := fmt.Sprintf("%s/service/rest/v1/security/service-account-tokens/%s",
		strings.TrimRight(os.Getenv("NXRM_SERVER_URL"), "/"), id)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(os.Getenv("NXRM_SERVER_USERNAME"), os.Getenv("NXRM_SERVER_PASSWORD"))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected 204, got %d", response.StatusCode)
	}
	return nil
}
