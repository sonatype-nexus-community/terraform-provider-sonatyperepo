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
	"terraform-provider-sonatyperepo/internal/provider/common"
	utils_test "terraform-provider-sonatyperepo/internal/provider/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemConfigHttpResource(t *testing.T) {
	randomString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: utils_test.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: sytemConfigHttpResourceHttpsOnlySimple(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_ENABLED, "false"),
					resource.TestCheckNoResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_AUTHENTICATION),
					resource.TestCheckNoResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_HOST),
					resource.TestCheckNoResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_PORT),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_HOST, fmt.Sprintf("my.proxy-%s.tld", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_PORT, "8080"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_USERNAME, fmt.Sprintf("proxy-user-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_PASSWORD, fmt.Sprintf("proxy-password-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_NTLM_DOMAIN, ""),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_NTLM_HOST, ""),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_NON_PROXY_HOSTS+".#", "0"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_RETRIES, fmt.Sprintf("%d", common.HTTP_SETTINGS_DEFAULT_RETRIES)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_TIMEOUT, fmt.Sprintf("%d", common.HTTP_SETTINGS_DEFAULT_TIMEOUT)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_USER_AGENT, ""),
				),
			},
			// Update - Full Config
			{
				Config: sytemConfigHttpResourceFull(randomString),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_HOST, fmt.Sprintf("my.http-proxy-%s.tld", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_PORT, "8080"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_USERNAME, fmt.Sprintf("http-user-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_PASSWORD, fmt.Sprintf("http-password-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_NTLM_DOMAIN, fmt.Sprintf("http.test-%s.tld", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTP_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_NTLM_HOST, fmt.Sprintf("http-host-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_HOST, fmt.Sprintf("my.https-proxy-%s.tld", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_PORT, "8080"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_ENABLED, "true"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_USERNAME, fmt.Sprintf("https-user-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_PASSWORD, fmt.Sprintf("https-password-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_NTLM_DOMAIN, fmt.Sprintf("https.test-%s.tld", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_HTTPS_PROXY+"."+RES_ATTR_AUTHENTICATION+"."+RES_ATTR_NTLM_HOST, fmt.Sprintf("https-host-%s", randomString)),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_NON_PROXY_HOSTS+".#", "2"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_RETRIES, "3"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_TIMEOUT, "30"),
					resource.TestCheckResourceAttr(RES_NAME_CONFIG_HTTP, RES_ATTR_USER_AGENT, fmt.Sprintf("test-%s", randomString)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func sytemConfigHttpResourceHttpsOnlySimple(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
	http_proxy = {
		enabled = false
	}
	https_proxy = {
		enabled = true
		host = "my.proxy-%s.tld"
		port = 8080
		authentication = {
			enabled = true
			username = "proxy-user-%s"
			password = "proxy-password-%s"
		}
	}
}
`, RES_TYPE_CONFIG_HTTP, randomString, randomString, randomString)
}

func sytemConfigHttpResourceFull(randomString string) string {
	return fmt.Sprintf(utils_test.ProviderConfig+`
resource "%s" "test" {
	http_proxy = {
		enabled = true
		host = "my.http-proxy-%s.tld"
		port = 8080
		authentication = {
			enabled = true
			username = "http-user-%s"
			password = "http-password-%s"
			ntlm_host = "http-host-%s"
			ntlm_domain = "http.test-%s.tld"
		}
	}
	https_proxy = {
		enabled = true
		host = "my.https-proxy-%s.tld"
		port = 8080
		authentication = {
			enabled = true
			username = "https-user-%s"
			password = "https-password-%s"
			ntlm_host = "https-host-%s"
			ntlm_domain = "https.test-%s.tld"
		}
	}
	non_proxy_hosts = ["127.0.0.1", "test-%s.localhost"]
	retries = 3
	timeout = 30
	user_agent = "test-%s"
}
`, RES_TYPE_CONFIG_HTTP, randomString, randomString, randomString, randomString, randomString, randomString, randomString, randomString, randomString, randomString, randomString, randomString)
}
