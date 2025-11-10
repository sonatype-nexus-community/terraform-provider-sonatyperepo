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

package common

type CapabilityType string

const (
	CAPABILITY_TYPE_AUDIT                     CapabilityType = "audit"
	CAPABILITY_TYPE_CORE_BASE_URL             CapabilityType = "baseurl"
	CAPABILITY_TYPE_CUSTOM_S3_REGIONS         CapabilityType = "customs3regions"
	CAPABILITY_TYPE_DEFAULT_ROLE              CapabilityType = "defaultrole"
	CAPABILITY_TYPE_FIREWALL_AUDIT_QUARANTINE CapabilityType = "firewall.audit"
	CAPABILITY_TYPE_OUTREACH                  CapabilityType = "OutreachManagementCapability"
	CAPABILITY_TYPE_RUT_AUTH                  CapabilityType = "rutauth"
	CAPABILITY_TYPE_UI_BRANDING               CapabilityType = "rapture.branding"
	CAPABILITY_TYPE_UI_SETTINGS               CapabilityType = "rapture.settings"
	CAPABILITY_TYPE_WEBHOOK_GLOBAL            CapabilityType = "webhook.global"
	CAPABILITY_TYPE_WEBHOOK_REPOSITORY        CapabilityType = "webhook.repository"
)

func (ct CapabilityType) String() string {
	return string(ct)
}

func (ct CapabilityType) StringPointer() *string {
	str := ct.String()
	return &str
}

type WebhookEventType string

func (wet WebhookEventType) String() string {
	return string(wet)
}

const (
	WEBHOOK_EVENT_TYPE_ASSET      WebhookEventType = "asset"
	WEBHOOK_EVENT_TYPE_AUDIT      WebhookEventType = "audit"
	WEBHOOK_EVENT_TYPE_COMPONENT  WebhookEventType = "component"
	WEBHOOK_EVENT_TYPE_REPOSITORY WebhookEventType = "repository"
)

func AllGlobalWebHookEventTypes() []string {
	return []string{
		WEBHOOK_EVENT_TYPE_AUDIT.String(),
		WEBHOOK_EVENT_TYPE_REPOSITORY.String(),
	}
}

func AllRepositoryWebHookEventTypes() []string {
	return []string{
		WEBHOOK_EVENT_TYPE_ASSET.String(),
		WEBHOOK_EVENT_TYPE_COMPONENT.String(),
	}
}
