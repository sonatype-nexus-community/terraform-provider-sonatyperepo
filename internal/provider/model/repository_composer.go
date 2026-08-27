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

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Composer Proxy
// ---------------------------------------
type RepositoryComposerProxyModel struct {
	RepositoryProxyModel
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryComposerProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryComposerProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)

	// NXRM 3.94+'s Composer proxy repository GET response does not reliably return a
	// populated `firewall` field, so the inline firewall mode sent on Create/Update can
	// end up unreadable. Mirror what was configured in the plan instead - this is the
	// only source of truth on NXRM 3.94+; on older versions the Capability-based path in
	// repository_common.go runs afterwards and overwrites this with the real Capability
	// data. Mirrors the same fix applied to Raw (#461). This resolves the immediate
	// "Provider produced inconsistent result after apply" error, but a subsequent
	// refresh can still show drift (Terraform wanting to re-add repository_firewall) -
	// it's not yet confirmed whether that's a permanent read-side gap like Raw's, or the
	// write itself not taking effect server-side for Composer. See
	// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/471
	if planModel.FirewallAuditAndQuarantine != nil {
		firewall := *planModel.FirewallAuditAndQuarantine
		firewall.CapabilityId = types.StringNull()
		m.FirewallAuditAndQuarantine = &firewall
	} else {
		m.FirewallAuditAndQuarantine = nil
	}
}

func (m *RepositoryComposerProxyModel) FromApiModel(api sonatyperepo.SimpleApiProxyRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = &RepositoryCleanupModel{}
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
	} else {
		m.Cleanup = nil
	}

	// Storage
	m.Storage.MapFromApi(&api.Storage)

	// Proxy Specific
	m.Proxy.MapFromApi(&api.Proxy)
	m.NegativeCache.MapFromApi(&api.NegativeCache)
	m.HttpClient.MapFromApiHttpClientAttributes(&api.HttpClient)
	m.RoutingRule = types.StringPointerValue(api.RoutingRuleName)
	if api.Replication != nil {
		m.Replication = &RepositoryReplicationModel{}
		m.Replication.MapFromApi(api.Replication)
	} else {
		// Set default values when API doesn't provide replication data
		m.Replication = &RepositoryReplicationModel{
			PreemptivePullEnabled: types.BoolValue(false),
			AssetPathRegex:        types.StringNull(),
		}
	}
}

func (m *RepositoryComposerProxyModel) ToApiCreateModel() sonatyperepo.ComposerProxyRepositoryApiRequest {
	apiModel := sonatyperepo.ComposerProxyRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
	}
	m.Storage.MapToApi(&apiModel.Storage)

	if m.Cleanup != nil {
		mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	}

	// Proxy Specific
	apiModel.Proxy = sonatyperepo.ProxyAttributes{}
	m.Proxy.MapToApi(&apiModel.Proxy)

	apiModel.NegativeCache = sonatyperepo.NegativeCacheAttributes{}
	m.NegativeCache.MapToApi(&apiModel.NegativeCache)

	apiModel.HttpClient = sonatyperepo.HttpClientAttributes{}
	m.HttpClient.MapToApiHttpClientAttributes(&apiModel.HttpClient)

	if m.Replication != nil {
		apiModel.Replication = &sonatyperepo.ReplicationAttributes{}
		m.Replication.MapToApi(apiModel.Replication)
	}

	apiModel.RoutingRule = m.RoutingRule.ValueStringPointer()

	return apiModel
}

func (m *RepositoryComposerProxyModel) ToApiUpdateModel() sonatyperepo.ComposerProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}
