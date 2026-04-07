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
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSecuritySslTruststore = "sonatyperepo_security_ssl_truststore"
	resourceNameSecuritySslTruststore = "sonatyperepo_security_ssl_truststore.test"

	// Self-signed test certificate (generated for testing purposes only)
	testPemCertificate = `-----BEGIN CERTIFICATE-----
MIICpDCCAYwCCQDU+pYq5raHMjANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAls
b2NhbGhvc3QwHhcNMjUwMTAxMDAwMDAwWhcNMjYwMTAxMDAwMDAwWjAUMRIwEAYD
VQQDDAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7
7RtHm2MTDMbMjZbCQT1LVKX4sZ9MgnDHhA9gXqH6JnGpVHH5FWoI7DCSCCf7wKP
qjsApqGDJqig0A54B0CZOR3YKXrFVcIkCzOFPkVaCfOyAKG9fC8eBuKz6SOEZFNQ
OQHm01lGCUmFQFKR/oqkBBUb1YqR0bfFH91FBIIqoMzFOIEIkFE1zFPMBiYGLCMS
TDTrNHOLJfBLhkfihIQRNI8gVVNRdFo6MMDlshOoWGU4FYFRRNTtMCMGjp/H3dO2
9VTHcRBJoFm4HEJOxmXGAJJMaxMlYEFBK7G1LRbOwS9M7NeI8WFLaWPXs3FMENR7
6FNjG2D/1YDnlEvxLiUhAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAKnEU97vtAqt
Dkq6A1ZFnOchRG7mJNR/KVfEJkFqQIMJbBBjGXfNPDUJMXppkyRw2DLPB7bYkfkH
0aDRcGCqMTQzMKz1RBYlX9VmqfDCVIEDKICr7WsGBgaYFocfaGKLUisFOOWQi3Ye
aVVjOCUGMA4U9RQvL9pHtHbqfiU9sdfIYBSJPKJGzJEviDOJFz9IGaOFMIg4dPMm
nwtB8AsBGRgNiMTeRaqhMik1KFJmqJIEiQ5ywI1S2CAiVFfzmCqBnMJNVAJFM8ZF
S5nm1s0YFJ7xSPFL7Wwn0p5VYdr3MxLhPMXsOJkFVPFe9K9yZr7RWwSezjyBNJmf
TBXWT8BKALw=
-----END CERTIFICATE-----`
)

func TestAccSecuritySslTruststoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: getSecuritySslTruststoreResourceConfig(testPemCertificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "id"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "fingerprint"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "subject_common_name"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "issued_on"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "expires_on"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "pem"),
					resource.TestCheckResourceAttrSet(resourceNameSecuritySslTruststore, "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:            resourceNameSecuritySslTruststore,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func getSecuritySslTruststoreResourceConfig(pem string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
  pem = <<-EOT
%s
EOT
}
`, resourceTypeSecuritySslTruststore, pem)
}
