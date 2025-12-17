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

import "fmt"

const (
	RES_ATTR_AUTHENTICATION  string = "authentication"
	RES_ATTR_ENABLED         string = "enabled"
	RES_ATTR_HTTP_PROXY      string = "http_proxy"
	RES_ATTR_HTTPS_PROXY     string = "https_proxy"
	RES_ATTR_HOST            string = "host"
	RES_ATTR_NTLM_DOMAIN     string = "ntlm_domain"
	RES_ATTR_NTLM_HOST       string = "ntlm_host"
	RES_ATTR_NON_PROXY_HOSTS string = "non_proxy_hosts"
	RES_ATTR_PASSWORD        string = "password"
	RES_ATTR_PORT            string = "port"
	RES_ATTR_RETRIES         string = "retries"
	RES_ATTR_TIMEOUT         string = "timeout"
	RES_ATTR_USER_AGENT      string = "user_agent"
	RES_ATTR_USERNAME        string = "username"
	RES_NAME_FMT             string = "%s.test"
	RES_TYPE_CONFIG_HTTP     string = "sonatyperepo_system_config_http"
)

var (
	RES_NAME_CONFIG_HTTP string = fmt.Sprintf(RES_NAME_FMT, RES_TYPE_CONFIG_HTTP)
)
