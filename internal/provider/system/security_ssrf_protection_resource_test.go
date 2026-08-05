/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSecuritySsrfProtection = "sonatyperepo_security_ssrf_protection"
	resourceNameSecuritySsrfProtection = "sonatyperepo_security_ssrf_protection.ssrf"
)

// skipIfSsrfProtectionUnsupported skips the test if NXRM does not yet support SSRF Protection (added in 3.92.0).
func skipIfSsrfProtectionUnsupported(t *testing.T) {
	t.Helper()
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

func TestAccSecuritySsrfProtectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.92.0
			skipIfSsrfProtectionUnsupported(t)
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getSecuritySsrfProtectionResourceConfig(true, []string{"internal.example.com"}, []string{"10.0.0.5"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceNameSecuritySsrfProtection, "allowed_domains.*", "internal.example.com"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_ips.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceNameSecuritySsrfProtection, "allowed_ips.*", "10.0.0.5"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceNameSecuritySsrfProtection,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				// Ignore last_updated since it will be different after import
				ImportStateVerifyIgnore: []string{"last_updated"},
				ImportStateVerify:       true,
				ImportStateId:           "ssrf_protection", // Can be any string for this singleton resource
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSecuritySsrfProtectionResourceImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.92.0
			skipIfSsrfProtectionUnsupported(t)
		},
		Steps: []resource.TestStep{
			// First, create a resource
			{
				Config: getSecuritySsrfProtectionResourceConfig(true, []string{"internal.example.com"}, []string{"10.0.0.5"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "enabled", "true"),
				),
			},
			// Test import with different import IDs (all should work for singleton resource)
			{
				ResourceName:                         resourceNameSecuritySsrfProtection,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "ssrf_protection",
			},
			{
				ResourceName:                         resourceNameSecuritySsrfProtection,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "import",
			},
			{
				ResourceName:                         resourceNameSecuritySsrfProtection,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "any-string-works",
			},
		},
	})
}

func TestAccSecuritySsrfProtectionResourceUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.92.0
			skipIfSsrfProtectionUnsupported(t)
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getSecuritySsrfProtectionResourceConfig(true, []string{"internal.example.com"}, []string{"10.0.0.5"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_ips.#", "1"),
				),
			},
			// Update domains/IPs and Read testing
			{
				Config: getSecuritySsrfProtectionResourceConfig(true, []string{"internal.example.com", "other.example.com"}, []string{"10.0.0.5", "192.168.1.1"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_domains.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceNameSecuritySsrfProtection, "allowed_domains.*", "other.example.com"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_ips.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceNameSecuritySsrfProtection, "allowed_ips.*", "192.168.1.1"),
				),
			},
			// Test import after update
			{
				ResourceName:                         resourceNameSecuritySsrfProtection,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "post-update-import",
			},
			// Test disabling SSRF Protection
			{
				Config: getSecuritySsrfProtectionResourceConfig(false, []string{}, []string{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_domains.#", "0"),
					resource.TestCheckResourceAttr(resourceNameSecuritySsrfProtection, "allowed_ips.#", "0"),
				),
			},
			// Test import when disabled
			{
				ResourceName:                         resourceNameSecuritySsrfProtection,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "enabled",
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated"},
				ImportStateId:                        "disabled-state-import",
			},
		},
	})
}

func getSecuritySsrfProtectionResourceConfig(enabled bool, allowedDomains []string, allowedIPs []string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "ssrf" {
	enabled = %t
	allowed_domains = %s
	allowed_ips = %s
}
`, resourceTypeSecuritySsrfProtection, enabled, toHclStringList(allowedDomains), toHclStringList(allowedIPs))
}

func toHclStringList(values []string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = fmt.Sprintf("%q", v)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}
