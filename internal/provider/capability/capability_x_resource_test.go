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

package capability_test

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceNameF = "%s.cap"
	notesFString  = "example-notes-%s"
	propertiesUrl = "properties.url"
	urlFString    = "https://%s.tld"
)

func TestAccCapabilityAuditResource(t *testing.T) {
	resourceName := fmt.Sprintf("%s.this", resourceAudit)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Import & Read Testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
			data "sonatyperepo_capabilities" "capabilities" {}

			import {
			  for_each = [for c in data.sonatyperepo_capabilities.capabilities.capabilities : c.id if c.type == "audit"]

			  id = [for c in data.sonatyperepo_capabilities.capabilities.capabilities : c.id if c.type == "audit"][0]
			  to = sonatyperepo_capability_audit.this
			}

			resource "%s" "this" {
			  enabled = true
			  
			}
			`, resourceAudit),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", ""),
				),
			},
			// Update Testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
			resource "%s" "this" {
			  enabled = true
			  notes = "Managed by Terraform"
			}
			`, resourceAudit),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by Terraform"),
				),
			},
		},
	})
}

func TestAccCapabilityCoreBaseUrlResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceBaseUrl)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  notes = "example-notes-%s"
  enabled = true
  properties = {
    url = "https://%s.tld"
  }
}
`, resourceBaseUrl, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, propertiesUrl, fmt.Sprintf(urlFString, randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityCoreStorageSettingsResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceStorageSettings)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  notes = "example-notes-%s"
  enabled = true
  properties = {
    last_downloaded_interval = 24
  }
}
`, resourceStorageSettings, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.last_downloaded_interval", "24"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityCustomS3RegionsResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceCustomS3Regions)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  notes = "example-notes-%s"
  enabled = true
  properties = {
    regions = [
		"somewhere-1-%s",
		"somewhere-2-%s"
	]
  }
}
`, resourceCustomS3Regions, randomString, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "properties.regions.*", fmt.Sprintf("somewhere-1-%s", randomString)),
					resource.TestCheckTypeSetElemAttr(resourceName, "properties.regions.*", fmt.Sprintf("somewhere-2-%s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityDefaultRoleResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceDefaultRole)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  notes = "example-notes-%s"
  enabled = false
  properties = {
    role = "nx-anonymous"
  }
}
`, resourceDefaultRole, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "properties.role", "nx-anonymous"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityFirewallAuditQuarantineResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceFirewallAuditQuarantine)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  notes = "example-notes-%s"
  enabled = true
  properties = {
    repository = "maven-central"
    quarantine = true
  }
}
`, resourceFirewallAuditQuarantine, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.repository", "maven-central"),
					resource.TestCheckResourceAttr(resourceName, "properties.quarantine", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityHealthcheckResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceHealthcheck)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  notes = "example-notes-%s"
  enabled = true
  properties = {
    configured_for_all_proxies = false
    use_nexus_truststore       = true
  }
}
`, resourceHealthcheck, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.configured_for_all_proxies", "false"),
					resource.TestCheckResourceAttr(resourceName, "properties.use_nexus_truststore", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilitySecurityRutAuthResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceSecurityRutAuth)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  enabled = true
  notes   = "example-notes-%s"
  properties = {
    http_header    = "TESTING 1 2 3 %s"
  }
}
`, resourceSecurityRutAuth, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.http_header", fmt.Sprintf("TESTING 1 2 3 %s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityUiBrandingResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceUiBranding)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  enabled = true
  notes   = "example-notes-%s"
  properties = {
    header_enabled = true
    header_html    = "TESTING 1 2 3 %s"
  }
}
`, resourceUiBranding, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.header_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "properties.header_html", fmt.Sprintf("TESTING 1 2 3 %s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityWebhookGlobalResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceWebhookGlobal)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  enabled = true
  notes   = "example-notes-%s"
  properties = {
    names = [
      "repository"
    ]
    url    = "https://%s.tld"
    secret = "super-secret-key-%s"
  }
}
`, resourceWebhookGlobal, randomString, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "properties.names.*", "repository"),
					resource.TestCheckResourceAttr(resourceName, propertiesUrl, fmt.Sprintf(urlFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "properties.secret", fmt.Sprintf("super-secret-key-%s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCapabilityWebhookRepositoryResource(t *testing.T) {

	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := fmt.Sprintf(resourceNameF, resourceWebhookRepository)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			// Not supported prior to NXRM 3.84.0
			testutil.SkipIfNxrmVersionInRange(t, &common.SystemVersion{
				Major: 3,
				Minor: 0,
				Patch: 0,
			}, &common.SystemVersion{
				Major: 3,
				Minor: 83,
				Patch: 99,
			})
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "cap" {
  enabled = true
  notes   = "example-notes-%s"
  properties = {
    names = [
      "asset"
    ]
    url    = "https://%s.tld"
    secret = "super-secret-key-%s"
	repository = "maven-central"
  }
}
`, resourceWebhookRepository, randomString, randomString, randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notes", fmt.Sprintf(notesFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckTypeSetElemAttr(resourceName, "properties.names.*", "asset"),
					resource.TestCheckResourceAttr(resourceName, propertiesUrl, fmt.Sprintf(urlFString, randomString)),
					resource.TestCheckResourceAttr(resourceName, "properties.secret", fmt.Sprintf("super-secret-key-%s", randomString)),
					resource.TestCheckResourceAttr(resourceName, "properties.repository", "maven-central"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
