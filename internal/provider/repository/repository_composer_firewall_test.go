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
	"strings"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	repotest "terraform-provider-sonatyperepo/internal/provider/repository/repotest"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRepositoryComposerProxyFirewallAtCreate reproduces
// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/471.
//
// It enables `repository_firewall` on a Composer proxy when the repository is CREATED - exactly
// the configuration in the issue - and then changes it in an update. Against NXRM 3.94+ with a
// connected Sonatype IQ Server it currently FAILS: the provider computes the inline firewall
// mode on the shared ComposerRepositoryFormat, whose firewall getters all return false (the
// real ones are on ComposerRepositoryFormatProxy), so it always sends mode DISABLED. NXRM stores
// DISABLED and reports it back, which fails the apply with "Provider produced inconsistent
// result after apply" and would otherwise show as drift on the next plan.
//
// NXRM itself is not at fault: a Composer proxy created directly through the REST API with
// firewall.mode = QUARANTINE stores and returns it.
//
// This is deliberately separate from TestAccRepositoryGenericProxyFirewallToggle, which never
// enables the firewall at create (its first step has no firewall) and excludes Composer. Once
// Composer is fixed it can be added to that table.
func TestAccRepositoryComposerProxyFirewallAtCreate(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceType := fmt.Sprintf(resourceTypeProxyFString, strings.ToLower(common.REPO_FORMAT_COMPOSER))
	resourceName := fmt.Sprintf(repoNameFString, resourceType)
	repoName := strings.ToLower(fmt.Sprintf(proxyNameFString, common.REPO_FORMAT_COMPOSER, randomString))

	config := func(firewallBlock string) string {
		return repositoryProxyResourceConfigWithFirewall(resourceType, repoName, TEST_DATA_COMPOSER_PROXY_REMOTE_URL, "", firewallBlock)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testutil.SkipIfNoIqServer(t)
			// repository_firewall requires NXRM 3.94+ (inline firewall)
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 93,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// 1. Created with the firewall already enabled (the configuration from #471)
			{
				Config: config(repositoryFirewallBlockHcl(true, true, false)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_NAME, repoName),
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE, "true"),
					resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_CAPABILITY_ID),
				),
			},
			// 2. Changed in an update: audit only
			{
				Config: config(repositoryFirewallBlockHcl(true, false, false)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_ENABLED, "true"),
					resource.TestCheckResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_QUARANTINE, "false"),
					resource.TestCheckNoResourceAttr(resourceName, repotest.RES_ATTR_REPOSITORY_FIREWALL_CAPABILITY_ID),
				),
			},
			// 3. Same configuration again: must be a clean plan (no drift on refresh)
			{
				Config:   config(repositoryFirewallBlockHcl(true, false, false)),
				PlanOnly: true,
			},
		},
	})
}
