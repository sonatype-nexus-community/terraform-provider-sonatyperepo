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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/testutil"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

const (
	securityOAuth2ResourceName = "sonatyperepo_security_oauth2.test"
	securityOAuth2ResourceType = "sonatyperepo_security_oauth2"
)

// skipIfOAuth2Unsupported skips the test if NXRM does not yet support the OAuth2/OIDC
// configuration API. Unlike the issue this resource implements (which cites 3.93.0), the
// underlying nexus-repo-api-client-go bindings only exist in the V395 client generation -
// nothing through v3.94.0 of the V382 generation has them - so 3.94.0 is the real floor this
// provider can support, not 3.93.0.
func skipIfOAuth2Unsupported(t *testing.T) {
	t.Helper()
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

func TestAccSecurityOAuth2ResourceBasic(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfOAuth2Unsupported(t)
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSecurityOAuth2ResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "idp_jws_algorithm", "RS256"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "username_claim", "preferred_username"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "first_name_claim", "given_name"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "last_name_claim", "family_name"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "email_claim", "email"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "groups_claim", "groups"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "use_trust_store", "false"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "client_id", fmt.Sprintf("test-client-%s", randomSuffix)),
					resource.TestCheckResourceAttrSet(securityOAuth2ResourceName, "idp_authorization_url"),
					resource.TestCheckResourceAttrSet(securityOAuth2ResourceName, "idp_token_url"),
					resource.TestCheckResourceAttrSet(securityOAuth2ResourceName, "idp_logout_url"),
					resource.TestCheckResourceAttrSet(securityOAuth2ResourceName, "idp_jwks_url"),
				),
			},
		},
	})
}

func TestAccSecurityOAuth2ResourceUpdate(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfOAuth2Unsupported(t)
		},
		Steps: []resource.TestStep{
			// Create initial resource
			{
				Config: testAccSecurityOAuth2ResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "username_claim", "preferred_username"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "client_id", fmt.Sprintf("test-client-%s", randomSuffix)),
				),
			},
			// Update resource
			{
				Config: testAccSecurityOAuth2ResourceConfigUpdated(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "idp_jws_algorithm", "RS384"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "username_claim", "sub"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "first_name_claim", "updatedFirstName"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "last_name_claim", "updatedLastName"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "email_claim", "updatedEmail"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "groups_claim", "updatedGroups"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "use_trust_store", "true"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "client_id", fmt.Sprintf("updated-client-%s", randomSuffix)),
				),
			},
		},
	})
}

func TestAccSecurityOAuth2ResourceMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfOAuth2Unsupported(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityOAuth2ResourceConfigMinimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "idp_jws_algorithm", "RS256"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "username_claim", "sub"),
				),
			},
		},
	})
}

func TestAccSecurityOAuth2ResourceCustomParams(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfOAuth2Unsupported(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityOAuth2ResourceConfigCustomParams(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "authorization_custom_params.prompt", "consent"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "token_request_custom_params.audience", "nexus-repository"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "exact_match_claims.aud", fmt.Sprintf("test-client-%s", randomSuffix)),
				),
			},
		},
	})
}

func TestAccSecurityOAuth2ResourceImport(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			skipIfOAuth2Unsupported(t)
		},
		Steps: []resource.TestStep{
			// Create initial resource
			{
				Config: testAccSecurityOAuth2ResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "username_claim", "preferred_username"),
					resource.TestCheckResourceAttr(securityOAuth2ResourceName, "client_id", fmt.Sprintf("test-client-%s", randomSuffix)),
				),
			},
			// Import existing resource - client_secret is write-only and never returned by the
			// API, so it can't be verified post-import.
			{
				ResourceName:            securityOAuth2ResourceName,
				ImportState:             true,
				ImportStateId:           "oauth2-configuration",
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
		},
	})
}

func testAccSecurityOAuth2ResourceConfigBasic(randomSuffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_jws_algorithm      = "RS256"
  username_claim         = "preferred_username"
  first_name_claim       = "given_name"
  last_name_claim        = "family_name"
  email_claim            = "email"
  groups_claim           = "groups"
  use_trust_store        = false
  client_id              = "test-client-%s"
  client_secret          = "test-secret-%s"
  idp_authorization_url  = "https://idp.example.com/oauth2/authorize"
  idp_token_url          = "https://idp.example.com/oauth2/token"
  idp_logout_url         = "https://idp.example.com/oauth2/logout"
  idp_jwks_url           = "https://idp.example.com/oauth2/jwks"
}
`, securityOAuth2ResourceType, randomSuffix, randomSuffix)
}

func testAccSecurityOAuth2ResourceConfigUpdated(randomSuffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_jws_algorithm      = "RS384"
  username_claim         = "sub"
  first_name_claim       = "updatedFirstName"
  last_name_claim        = "updatedLastName"
  email_claim            = "updatedEmail"
  groups_claim           = "updatedGroups"
  use_trust_store        = true
  client_id              = "updated-client-%s"
  client_secret          = "updated-secret-%s"
  idp_authorization_url  = "https://idp.example.com/oauth2/authorize"
  idp_token_url          = "https://idp.example.com/oauth2/token"
  idp_logout_url         = "https://idp.example.com/oauth2/logout"
  idp_jwks_url           = "https://idp.example.com/oauth2/jwks"
}
`, securityOAuth2ResourceType, randomSuffix, randomSuffix)
}

func testAccSecurityOAuth2ResourceConfigMinimal() string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_jws_algorithm = "RS256"
  username_claim    = "sub"
}
`, securityOAuth2ResourceType)
}

func testAccSecurityOAuth2ResourceConfigCustomParams(randomSuffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_jws_algorithm = "RS256"
  username_claim    = "preferred_username"
  client_id         = "test-client-%s"

  authorization_custom_params = {
    prompt = "consent"
  }
  token_request_custom_params = {
    audience = "nexus-repository"
  }
  exact_match_claims = {
    aud = "test-client-%s"
  }
}
`, securityOAuth2ResourceType, randomSuffix, randomSuffix)
}
