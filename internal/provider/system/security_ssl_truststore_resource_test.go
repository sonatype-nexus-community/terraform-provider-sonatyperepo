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
	"regexp"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resourceTypeSecuritySslTruststore = "sonatyperepo_security_ssl_truststore"
	resourceNameSecuritySslTruststore = "sonatyperepo_security_ssl_truststore.test"

	// Self-signed test certificate (generated for testing purposes only)
	testPemCertificate = `-----BEGIN CERTIFICATE-----
MIIOODCCDSCgAwIBAgIRAKzmW6rKKuHYCZjKWbMmnVMwDQYJKoZIhvcNAQELBQAw
OzELMAkGA1UEBhMCVVMxHjAcBgNVBAoTFUdvb2dsZSBUcnVzdCBTZXJ2aWNlczEM
MAoGA1UEAxMDV1IyMB4XDTI2MDMxNjA4MzYzMloXDTI2MDYwODA4MzYzMVowFzEV
MBMGA1UEAwwMKi5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
Qpq8TmMXJmbelhFp9HK1tG9c4PZnNN4KBJUwieoxjKHurNnQFzPNjSlV74r+X5TP
FdLBTd1mRb27Y75v0roqJaOCDCQwggwgMA4GA1UdDwEB/wQEAwIHgDATBgNVHSUE
DDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQlAU4wPdtfQXNa
XSFo3aYZer1X0zAfBgNVHSMEGDAWgBTeGx7teRXUPjckwyG77DQ5bUKyMDBYBggr
BgEFBQcBAQRMMEowIQYIKwYBBQUHMAGGFWh0dHA6Ly9vLnBraS5nb29nL3dyMjAl
BggrBgEFBQcwAoYZaHR0cDovL2kucGtpLmdvb2cvd3IyLmNydDCCCfgGA1UdEQSC
Ce8wggnrggwqLmdvb2dsZS5jb22CFiouYXBwZW5naW5lLmdvb2dsZS5jb22CCSou
YmRuLmRldoIVKi5vcmlnaW4tdGVzdC5iZG4uZGV2ghIqLmNsb3VkLmdvb2dsZS5j
b22CGCouY3Jvd2Rzb3VyY2UuZ29vZ2xlLmNvbYIYKi5kYXRhY29tcHV0ZS5nb29n
bGUuY29tggsqLmdvb2dsZS5jYYILKi5nb29nbGUuY2yCDiouZ29vZ2xlLmNvLmlu
gg4qLmdvb2dsZS5jby5qcIIOKi5nb29nbGUuY28udWuCDyouZ29vZ2xlLmNvbS5h
coIPKi5nb29nbGUuY29tLmF1gg8qLmdvb2dsZS5jb20uYnKCDyouZ29vZ2xlLmNv
bS5jb4IPKi5nb29nbGUuY29tLm14gg8qLmdvb2dsZS5jb20udHKCDyouZ29vZ2xl
LmNvbS52boILKi5nb29nbGUuZGWCCyouZ29vZ2xlLmVzggsqLmdvb2dsZS5mcoIL
Ki5nb29nbGUuaHWCCyouZ29vZ2xlLml0ggsqLmdvb2dsZS5ubIILKi5nb29nbGUu
cGyCCyouZ29vZ2xlLnB0gg8qLmdvb2dsZWFwaXMuY26CDCouZ3N0YXRpYy5jboIQ
Ki5nc3RhdGljLWNuLmNvbYIPZ29vZ2xlY25hcHBzLmNughEqLmdvb2dsZWNuYXBw
cy5jboIRZ29vZ2xlYXBwcy1jbi5jb22CEyouZ29vZ2xlYXBwcy1jbi5jb22CDGdr
ZWNuYXBwcy5jboIOKi5na2VjbmFwcHMuY26CEmdvb2dsZWRvd25sb2Fkcy5jboIU
Ki5nb29nbGVkb3dubG9hZHMuY26CEHJlY2FwdGNoYS5uZXQuY26CEioucmVjYXB0
Y2hhLm5ldC5jboIQcmVjYXB0Y2hhLWNuLm5ldIISKi5yZWNhcHRjaGEtY24ubmV0
ggt3aWRldmluZS5jboINKi53aWRldmluZS5jboIRYW1wcHJvamVjdC5vcmcuY26C
EyouYW1wcHJvamVjdC5vcmcuY26CEWFtcHByb2plY3QubmV0LmNughMqLmFtcHBy
b2plY3QubmV0LmNughdnb29nbGUtYW5hbHl0aWNzLWNuLmNvbYIZKi5nb29nbGUt
YW5hbHl0aWNzLWNuLmNvbYIXZ29vZ2xlYWRzZXJ2aWNlcy1jbi5jb22CGSouZ29v
Z2xlYWRzZXJ2aWNlcy1jbi5jb22CEWdvb2dsZXZhZHMtY24uY29tghMqLmdvb2ds
ZXZhZHMtY24uY29tghFnb29nbGVhcGlzLWNuLmNvbYITKi5nb29nbGVhcGlzLWNu
LmNvbYIVZ29vZ2xlb3B0aW1pemUtY24uY29tghcqLmdvb2dsZW9wdGltaXplLWNu
LmNvbYISZG91YmxlY2xpY2stY24ubmV0ghQqLmRvdWJsZWNsaWNrLWNuLm5ldIIY
Ki5mbHMuZG91YmxlY2xpY2stY24ubmV0ghYqLmcuZG91YmxlY2xpY2stY24ubmV0
gg5kb3VibGVjbGljay5jboIQKi5kb3VibGVjbGljay5jboIUKi5mbHMuZG91Ymxl
Y2xpY2suY26CEiouZy5kb3VibGVjbGljay5jboIRZGFydHNlYXJjaC1jbi5uZXSC
EyouZGFydHNlYXJjaC1jbi5uZXSCHWdvb2dsZXRyYXZlbGFkc2VydmljZXMtY24u
Y29tgh8qLmdvb2dsZXRyYXZlbGFkc2VydmljZXMtY24uY29tghhnb29nbGV0YWdz
ZXJ2aWNlcy1jbi5jb22CGiouZ29vZ2xldGFnc2VydmljZXMtY24uY29tghdnb29n
bGV0YWdtYW5hZ2VyLWNuLmNvbYIZKi5nb29nbGV0YWdtYW5hZ2VyLWNuLmNvbYIY
Z29vZ2xlc3luZGljYXRpb24tY24uY29tghoqLmdvb2dsZXN5bmRpY2F0aW9uLWNu
LmNvbYIkKi5zYWZlZnJhbWUuZ29vZ2xlc3luZGljYXRpb24tY24uY29tghZhcHAt
bWVhc3VyZW1lbnQtY24uY29tghgqLmFwcC1tZWFzdXJlbWVudC1jbi5jb22CC2d2
dDEtY24uY29tgg0qLmd2dDEtY24uY29tggtndnQyLWNuLmNvbYINKi5ndnQyLWNu
LmNvbYILMm1kbi1jbi5uZXSCDSouMm1kbi1jbi5uZXSCFGdvb2dsZWZsaWdodHMt
Y24ubmV0ghYqLmdvb2dsZWZsaWdodHMtY24ubmV0ggxhZG1vYi1jbi5jb22CDiou
YWRtb2ItY24uY29tghkqLmdlbWluaS5jbG91ZC5nb29nbGUuY29tghRnb29nbGVz
YW5kYm94LWNuLmNvbYIWKi5nb29nbGVzYW5kYm94LWNuLmNvbYIeKi5zYWZlbnVw
Lmdvb2dsZXNhbmRib3gtY24uY29tgg0qLmdzdGF0aWMuY29tghQqLm1ldHJpYy5n
c3RhdGljLmNvbYIKKi5ndnQxLmNvbYIRKi5nY3BjZG4uZ3Z0MS5jb22CCiouZ3Z0
Mi5jb22CDiouZ2NwLmd2dDIuY29tghAqLnVybC5nb29nbGUuY29tghYqLnlvdXR1
YmUtbm9jb29raWUuY29tggsqLnl0aW1nLmNvbYIKYWkuYW5kcm9pZIILYW5kcm9p
ZC5jb22CDSouYW5kcm9pZC5jb22CEyouZmxhc2guYW5kcm9pZC5jb22CBGcuY26C
BiouZy5jboIEZy5jb4IGKi5nLmNvggZnb28uZ2yCCnd3dy5nb28uZ2yCFGdvb2ds
ZS1hbmFseXRpY3MuY29tghYqLmdvb2dsZS1hbmFseXRpY3MuY29tggpnb29nbGUu
Y29tghJnb29nbGVjb21tZXJjZS5jb22CFCouZ29vZ2xlY29tbWVyY2UuY29tgghn
Z3BodC5jboIKKi5nZ3BodC5jboIKdXJjaGluLmNvbYIMKi51cmNoaW4uY29tggh5
b3V0dS5iZYILeW91dHViZS5jb22CDSoueW91dHViZS5jb22CEW11c2ljLnlvdXR1
YmUuY29tghMqLm11c2ljLnlvdXR1YmUuY29tghR5b3V0dWJlZWR1Y2F0aW9uLmNv
bYIWKi55b3V0dWJlZWR1Y2F0aW9uLmNvbYIPeW91dHViZWtpZHMuY29tghEqLnlv
dXR1YmVraWRzLmNvbYIFeXQuYmWCByoueXQuYmWCGmFuZHJvaWQuY2xpZW50cy5n
b29nbGUuY29tghMqLmFuZHJvaWQuZ29vZ2xlLmNughIqLmNocm9tZS5nb29nbGUu
Y26CFiouZGV2ZWxvcGVycy5nb29nbGUuY26CFSouYWlzdHVkaW8uZ29vZ2xlLmNv
bTATBgNVHSAEDDAKMAgGBmeBDAECATA2BgNVHR8ELzAtMCugKaAnhiVodHRwOi8v
Yy5wa2kuZ29vZy93cjIvb0JGWVlhaHpnVkkuY3JsMIIBBgYKKwYBBAHWeQIEAgSB
9wSB9ADyAHcAlpdkv1VYl633Q4doNwhCd+nwOtX2pPM2bkakPw/KqcYAAAGc9gEi
sQAABAMASDBGAiEAjm7Vok6iTPht4FG6mczcEOwTotmezdd4eQ+XIKTwuFsCIQCE
5zkcieek8bQkp0JiS7cdHLbnAk6Bg0QsXvNjl04mwgB3AEmcm2neHXzs/DbezYdk
prhbrwqHgBnRVVL76esp3fjDAAABnPYBIpIAAAQDAEgwRgIhAJ713zJl9sPcewpZ
bqxp/R0jGjcYJTCmX37IdH0hXg/dAiEAuE/1ccPJZAqzNWbekGDKRo9fkwE+69PO
5E23Z6hpbSYwDQYJKoZIhvcNAQELBQADggEBAE3jhj4JQ8OilUEL08a1T94fYJDO
aUVvivnmNY7xM4yADKUKLCf8x5ZvOKFJ93ufHPRbfe7jzWodI9YBsg9gQN91f9tb
8jx8v23hPQ4+3Zfd5kiIckNBsaSXjhrQvfDqnDTdKnwfr4N9DWwHs6snsvdPh+rl
1AFSD9a6MX72NApdK4zpU5Ean48RcP5R3lHo7TY19zhR+T/oDdFYOPB5A/AzUmOt
aRZEJ/VvrdB+Agaifrj3DttjSkKT3J4xv+S12jYkztCxEdgYYWDQf3P+aBXVGbu8
S+YoF/3aSIwoFmhfnsk4UMFPxBRQQsrtu1ETLm21TZcZxHbN/76eoXmpD7w=
-----END CERTIFICATE-----`
)

func TestAccSecuritySslTruststoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Use PEM with additional trailing \n - triggering validation error
			{
				Config:      getSecuritySslTruststoreResourceConfig(testPemCertificate + "\n"),
				ExpectError: regexp.MustCompile("Attribute pem PEM data not in expected format - only a single newline"),
			},
			// Create and Read testing
			{
				Config: getSecuritySslTruststoreResourceConfig(testPemCertificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "id", "DE:18:EF:31:96:48:61:26:47:8F:5F:EE:5F:C5:7E:A7:0F:31:50:C1"),
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "fingerprint", "DE:18:EF:31:96:48:61:26:47:8F:5F:EE:5F:C5:7E:A7:0F:31:50:C1"),
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "serial_number", "229823302783097913606945207554866060627"),
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "subject_common_name", "*.google.com"),
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "issued_on", "1773650192000"),
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "expires_on", "1780907791000"),
					resource.TestCheckResourceAttr(resourceNameSecuritySslTruststore, "pem", testPemCertificate+"\n"),
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
