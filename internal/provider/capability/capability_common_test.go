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

const (
	resourceAudit                   string = "sonatyperepo_capability_audit"
	resourceBaseUrl                 string = "sonatyperepo_capability_base_url"
	resourceCustomS3Regions         string = "sonatyperepo_capability_custom_s3_regions"
	resourceDefaultRole             string = "sonatyperepo_capability_default_role"
	resourceFirewallAuditQuarantine string = "sonatyperepo_capability_firewall_audit_and_quarantine"
	resourceHealthcheck             string = "sonatyperepo_capability_healthcheck"
	resourceOutreachManagement      string = "sonatyperepo_capability_outreach_management"
	resourceSecurityRutAuth         string = "sonatyperepo_capability_rut_auth"
	resourceStorageSettings         string = "sonatyperepo_capability_storage_settings"
	resourceUiBranding              string = "sonatyperepo_capability_ui_branding"
	resourceUiSettings              string = "sonatyperepo_capability_ui_settings"
	resourceWebhookGlobal           string = "sonatyperepo_capability_webhook_global"
	resourceWebhookRepository       string = "sonatyperepo_capability_webhook_repository"

	resourceAttrId      string = "id"
	resourceAttrEnabled string = "enabled"
	resourceAttrNotes   string = "notes"

	resourceAttrPropertiesFormat            string = "properties.%s"
	resourceAttrAlwaysRemote                string = "always_remote"
	resourceAttrConfiguredForAllProxies     string = "configured_for_all_proxies"
	resourceAttrDebugAllowed                string = "debug_allowed"
	resourceAttrHeaderEnabled               string = "header_enabled"
	resourceAttrHeaderHtml                  string = "header_html"
	resourceAttrHttpHeader                  string = "http_header"
	resourceAttrLongRequestTimeout          string = "long_request_timeout"
	resourceAttrNames                       string = "names.*"
	resourceAttrOverrideUrl                 string = "override_url"
	resourceAttrRegionsStar                 string = "regions.*"
	resourceAttrRepository                  string = "repository"
	resourceAttrRequestTimeout              string = "request_timeout"
	resourceAttrRole                        string = "role"
	resourceAttrQuarantine                  string = "quarantine"
	resourceAttrSecret                      string = "secret"
	resourceAttrSessionTimeout              string = "session_timeout"
	resourceAttrStatusIntervalAnonymous     string = "status_interval_anonymous"
	resourceAttrStatusIntervalAuthenticated string = "status_interval_authenticated"
	resourceAttrTitle                       string = "title"
	resourceAttrUseNexusTruststore          string = "use_nexus_truststore"
)
