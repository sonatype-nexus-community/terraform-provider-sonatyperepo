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

	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
)

const (
	securitySamlResourceName = "sonatyperepo_security_saml.test"
	securitySamlResourceType = "sonatyperepo_security_saml"
)

func TestAccSecuritySamlResourceBasic(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSecuritySamlResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securitySamlResourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "first_name_attribute", "firstName"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "last_name_attribute", "lastName"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "email_attribute", "email"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "groups_attribute", "groups"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_response_signature", "true"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_assertion_signature", "false"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "entity_id", fmt.Sprintf("test-entity-%s", randomSuffix)),
					resource.TestCheckResourceAttrSet(securitySamlResourceName, "idp_metadata"),
				),
			},
		},
	})
}

func TestAccSecuritySamlResourceUpdate(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create initial resource
			{
				Config: testAccSecuritySamlResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securitySamlResourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_response_signature", "true"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "entity_id", fmt.Sprintf("test-entity-%s", randomSuffix)),
				),
			},
			// Update resource
			{
				Config: testAccSecuritySamlResourceConfigUpdated(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securitySamlResourceName, "username_attribute", "updatedUsername"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "first_name_attribute", "updatedFirstName"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "last_name_attribute", "updatedLastName"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "email_attribute", "updatedEmail"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "groups_attribute", "updatedGroups"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_response_signature", "false"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_assertion_signature", "true"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "entity_id", fmt.Sprintf("updated-entity-%s", randomSuffix)),
				),
			},
		},
	})
}

func TestAccSecuritySamlResourceMinimal(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecuritySamlResourceConfigMinimal(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securitySamlResourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttrSet(securitySamlResourceName, "idp_metadata"),
				),
			},
		},
	})
}

func TestAccSecuritySamlResourceImport(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create initial resource
			{
				Config: testAccSecuritySamlResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securitySamlResourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "first_name_attribute", "firstName"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "last_name_attribute", "lastName"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "email_attribute", "email"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "groups_attribute", "groups"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_response_signature", "true"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "validate_assertion_signature", "false"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "entity_id", fmt.Sprintf("test-entity-%s", randomSuffix)),
					resource.TestCheckResourceAttrSet(securitySamlResourceName, "idp_metadata"),
				),
			},
			// Import existing resource - will have placeholder values initially
			{
				ResourceName:            securitySamlResourceName,
				ImportState:             true,
				ImportStateId:           "saml-configuration",
				ImportStateVerify:       false, // Can't verify since imported values are placeholders
				ImportStateVerifyIgnore: []string{},
			},
			// Apply configuration after import to set actual values
			{
				Config: testAccSecuritySamlResourceConfigBasic(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(securitySamlResourceName, "username_attribute", "username"),
					resource.TestCheckResourceAttr(securitySamlResourceName, "entity_id", fmt.Sprintf("test-entity-%s", randomSuffix)),
				),
			},
		},
	})
}

func testAccSecuritySamlResourceConfigBasic(randomSuffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_metadata = %s
  username_attribute = "username"
  first_name_attribute = "firstName"
  last_name_attribute = "lastName"
  email_attribute = "email"
  groups_attribute = "groups"
  validate_response_signature = true
  validate_assertion_signature = false
  entity_id = "test-entity-%s"
}
`, securitySamlResourceType, testSamlMetadata(), randomSuffix)
}

func testAccSecuritySamlResourceConfigUpdated(randomSuffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_metadata = %s
  username_attribute = "updatedUsername"
  first_name_attribute = "updatedFirstName"
  last_name_attribute = "updatedLastName"
  email_attribute = "updatedEmail"
  groups_attribute = "updatedGroups"
  validate_response_signature = false
  validate_assertion_signature = true
  entity_id = "updated-entity-%s"
}
`, securitySamlResourceType, testSamlMetadata(), randomSuffix)
}

func testAccSecuritySamlResourceConfigMinimal(randomSuffix string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  idp_metadata = %s
  username_attribute = "username"
}
`, securitySamlResourceType, testSamlMetadata())
}

func testSamlMetadata() string {
	return `<<-EOT
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" ID="test-id" entityID="https://test.example.com/saml">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <KeyDescriptor use="signing">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
        <X509Data>
          <X509Certificate>
MIIEETCCAvmgAwIBAgIUB85DwdgvP7Lp6hIv4R6F84p3H3wwDQYJKoZIhvcNAQEL
BQAwgZcxCzAJBgNVBAYTAkFVMS4wLAYDVQQIDCUjIFNoYXJlIG9ubHkgdGhlIGNl
cnRpZmljYXRlIHdpdGggSWRQMTUwMwYDVQQHDCwjIEtlZXAgdGhlIHByaXZhdGUg
a2V5IHNlY3VyZSBvbiB5b3VyIHNlcnZlcjEhMB8GA1UECgwYSW50ZXJuZXQgV2lk
Z2l0cyBQdHkgTHRkMB4XDTI1MDkxNDIzNDgwN1oXDTI2MDkxNDIzNDgwN1owgZcx
CzAJBgNVBAYTAkFVMS4wLAYDVQQIDCUjIFNoYXJlIG9ubHkgdGhlIGNlcnRpZmlj
YXRlIHdpdGggSWRQMTUwMwYDVQQHDCwjIEtlZXAgdGhlIHByaXZhdGUga2V5IHNl
Y3VyZSBvbiB5b3VyIHNlcnZlcjEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQ
dHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1nIHufmHKbyj
PNDb/KaZ6ppXuMdy87/yXsTikH9dDQe5Ya/reihsfGNY3+uEOhRygKiqasHWkseZ
YLn1rLnhDAwGp/8MTcLYqY3S0F2eF/1q1SjiRJz9++svMLFwoMLHie64zqJJINqA
B0Vp0hu7HVCU8ZSKPBfS6PW5w5doO9huCCyRviabojG9RfqSk5VkQLG3TsjpMwaI
UzeLlMC9GTqmW3moUwkwoMnbKT4YCq/tNHwpDGq0Du1y7B2o+0Z+pihCLN5WYEZT
jiLPvGOpzQcSI3n8Te+qoU7J1qQ6c1BJhs2GUF2cOhLQIvLgHDNB0v94Xfz3ma5Y
ft5REQAHwQIDAQABo1MwUTAdBgNVHQ4EFgQUWwVsy9QgHO/bS81VCkINE9UujPMw
HwYDVR0jBBgwFoAUWwVsy9QgHO/bS81VCkINE9UujPMwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAQEASjabAHp0KqyyyTSVta6x4QAL5Vxg9vcTAklK
iXjvMOzExXYm1zI9eOvIzj1fn0e0tLOA03LGLHkLdQK9A5x/MqXysmK8sXsy5JBv
5WwvFE/7kKefaGsHSOLiXgTpgJtf3X0Z3B4s8y3Y2i6KvPA6uGElNbtd7WvTTFwU
oMhYx6Dkmek4o31fqmvuo3JbOPz2+zJOcD1Y3qL7He7WaSBLA9s3CQwYcTPBjy6i
HYoEGvzPl9flJfVRQ/GdVQxYZTAtmOgy+nSZwkm4oz2Xoz3q8/kHoOIBuPB6iWzf
XBWQDGEjEr6wHCaLl3CLN3yIKxXzJW2wqFvpm0MNcM0E9W9gXQ==
          </X509Certificate>
        </X509Data>
      </KeyInfo>
    </KeyDescriptor>
    <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://test.example.com/saml/logout"/>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://test.example.com/saml/sso"/>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://test.example.com/saml/sso"/>
  </IDPSSODescriptor>
</EntityDescriptor>
EOT`
}