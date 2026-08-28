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
	"os"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	defaultIqServerUrl          = "http://localhost:8070"
	resourceTypeSysIqConnection = "sonatyperepo_system_iq_connection"
	resourceNameSysIqConnection = "sonatyperepo_system_iq_connection.iq"
)

// settleIqConnectionWrite pauses briefly before returning, for use as the first entry in a
// step's Check chain (which - see testing_new_config.go in terraform-plugin-testing - runs
// immediately after apply and before the framework's own automatic post-apply refresh/empty-plan
// verification). Observed in CI (twice, on different NXRM versions, ~1 in 28 jobs each time):
// that automatic refresh occasionally reads back a stale value for a field this test just wrote
// (e.g. show_iq_server_link), even with no other writer active at the time - consistent with a
// short-lived read-after-write staleness on NXRM's side for this singleton, not a genuine
// provider bug or a race with another test. See
// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285.
func settleIqConnectionWrite(_ *terraform.State) error {
	time.Sleep(3 * time.Second)
	return nil
}

func TestAccSystemIqConnectionResource(t *testing.T) {
	// This test's steps intentionally write conflicting sonatyperepo_system_iq_connection
	// configuration (fake credentials, toggling show_iq_server_link/nexus_trust_store_enabled),
	// and resource.Test's implicit end-of-test destroy calls DisableIq - all directly against
	// the same singleton that internal/provider/provider_test.go's TestMain bootstraps at the
	// start of its own package's test binary. Since `go test ./...` compiles each package to a
	// separate binary and can run them concurrently, hold a cross-process lock for the entire
	// test (steps + destroy + restore) so that bootstrap can't interleave with it and produce
	// spurious refresh drift (observed in CI - see
	// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285).
	if os.Getenv("TF_ACC_IQ_SERVER") == "1" {
		unlock, err := testutil.LockIqConnection(2 * time.Minute)
		if err != nil {
			t.Fatalf("Failed to acquire IQ connection lock: %v", err)
		}
		// Cleanups run last-added-first: register the unlock first (so it runs last, releasing
		// only after the restore below has run), then the restore (so it runs first, while the
		// lock is still held).
		t.Cleanup(unlock)
		t.Cleanup(func() {
			if err := testutil.ConfigureIqConnection(); err != nil {
				t.Logf("Failed to restore Sonatype IQ Server connection after test: %v", err)
			}
		})
	}

	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: systemIqConnectionResourceConfig(false),
			Check: resource.ComposeAggregateTestCheckFunc(
				settleIqConnectionWrite,
				// Verify
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "true"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", defaultIqServerUrl),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", ""),
				resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
			),
		},
		// Test enabling show_iq_server_link
		{
			Config: systemIqConnectionResourceConfig(true),
			Check: resource.ComposeAggregateTestCheckFunc(
				settleIqConnectionWrite,
				// Verify
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "true"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", defaultIqServerUrl),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "true"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", ""),
				resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
			),
		},
		// Test adding Properties
		{
			Config: systemIqConnectionWithPropertiesResourceConfig(true),
			Check: resource.ComposeAggregateTestCheckFunc(
				settleIqConnectionWrite,
				// Verify
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "authentication_method", "USER"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "enabled", "true"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "fail_open_mode_enabled", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "nexus_trust_store_enabled", "false"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "url", defaultIqServerUrl),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "username", "user"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "password", "token"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "show_iq_server_link", "true"),
				resource.TestCheckResourceAttr(resourceNameSysIqConnection, "properties", "key1=value1&key2=value2"),
				resource.TestCheckResourceAttrSet(resourceNameSysIqConnection, "last_updated"),
			),
		},
		// Test that nexus_trust_store_enabled is applied
		{
			Config: systemIqConnectionWithTrustStoreConfig(true),
			Check: resource.ComposeAggregateTestCheckFunc(
				settleIqConnectionWrite,
				resource.TestCheckResourceAttr(resourceNameSysIqConnection,
					"nexus_trust_store_enabled", "true"),
			),
		},
		// Verify it can be set back to false
		{
			Config: systemIqConnectionWithTrustStoreConfig(false),
			Check: resource.ComposeAggregateTestCheckFunc(
				settleIqConnectionWrite,
				resource.TestCheckResourceAttr(resourceNameSysIqConnection,
					"nexus_trust_store_enabled", "false"),
			),
		},
		// ImportState testing
		{
			ResourceName:                         resourceNameSysIqConnection,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "url",
			ImportStateVerifyIgnore: []string{
				"password",
				"last_updated",
			},
			ImportStateId: "system-iq-config",
		},
		// Delete testing automatically occurs in TestCase
	}

	// Against a real, licensed Sonatype IQ Server (TF_ACC_IQ_SERVER=1), confirm the connection
	// NXRM ends up with actually works - the steps above only ever exercise fake credentials
	// against an unreachable URL, which NXRM accepts without validating at write time.
	if os.Getenv("TF_ACC_IQ_SERVER") == "1" {
		steps = append(steps, resource.TestStep{
			Config: systemIqConnectionRealResourceConfig(),
			Check: func(_ *terraform.State) error {
				success, reason, err := testutil.VerifyIqConnection()
				if err != nil {
					return fmt.Errorf("error verifying Sonatype IQ Server connection: %w", err)
				}
				if !success {
					return fmt.Errorf("Sonatype IQ Server connection was not successful: %s", reason)
				}
				return nil
			},
		})
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// systemIqConnectionRealResourceConfig points sonatyperepo_system_iq_connection at the real IQ
// Server under test (testutil.IqServerUrl/Username/Password), unlike the fake-credential
// configs below which only ever exercise attribute CRUD against NXRM's unvalidated write path.
func systemIqConnectionRealResourceConfig() string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = true
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "%s"
  username = "%s"
  password = "%s"
  show_iq_server_link = true
}
`, resourceTypeSysIqConnection, testutil.IqServerUrl(), testutil.IqServerUsername(), testutil.IqServerPassword())
}

func systemIqConnectionResourceConfig(showIqLink bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = true
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "%s"
  username = "user"
  password = "token"
  show_iq_server_link = %t
}
`, resourceTypeSysIqConnection, defaultIqServerUrl, showIqLink)
}

func systemIqConnectionWithPropertiesResourceConfig(showIqLink bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method = "USER"
  enabled = true
  fail_open_mode_enabled = false
  nexus_trust_store_enabled = false
  url = "%s"
  username = "user"
  password = "token"
  show_iq_server_link = %t
  properties = "key1=value1&key2=value2"
}
`, resourceTypeSysIqConnection, defaultIqServerUrl, showIqLink)
}

func systemIqConnectionWithTrustStoreConfig(nexusTrustStoreEnabled bool) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "iq" {
  authentication_method     = "USER"
  enabled                   = true
  fail_open_mode_enabled    = false
  nexus_trust_store_enabled = %t
  url 					    = "%s"
  username 					= "user"
  password                  = "token"
  show_iq_server_link       = true
}
`, resourceTypeSysIqConnection, nexusTrustStoreEnabled, defaultIqServerUrl)
}
