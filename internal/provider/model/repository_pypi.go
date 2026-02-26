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
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Hosted PyPi
// --------------------------------------------
type RepositoryPyPiHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositoryPyPiHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositoryPyPiHostedModel) ToApiCreateModel() sonatyperepo.PypiHostedRepositoryApiRequest {
	apiModel := sonatyperepo.PypiHostedRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{},
		Component: &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: common.NewFalse(),
		},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	m.Component.MapToApi(apiModel.Component)

	return apiModel
}

func (m *RepositoryPyPiHostedModel) ToApiUpdateModel() sonatyperepo.PypiHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Proxy PyPI
// --------------------------------------------
type RepositoryPyPiProxyModel struct {
	RepositoryProxyModel
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineWithPccsModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryPyPiProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryPyPiProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)
}

func (m *RepositoryPyPiProxyModel) FromApiModel(api sonatyperepo.PyPiProxyApiRepository) {
	m.Name = types.StringValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringValue(api.Url)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = NewRepositoryCleanupModel()
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
			PreemptivePullEnabled: types.BoolValue(common.DEFAULT_PROXY_PREEMPTIVE_PULL),
			AssetPathRegex:        types.StringNull(),
		}
	}

	// Firewall Audit and Quarantine (with PCCS)
	if m.FirewallAuditAndQuarantine != nil && api.Pypi.RemoveQuarantined {
		m.FirewallAuditAndQuarantine.PccsEnabled = types.BoolValue(api.Pypi.RemoveQuarantined)
	}
}

func (m *RepositoryPyPiProxyModel) ToApiCreateModel() sonatyperepo.PypiProxyRepositoryApiRequest {
	apiModel := sonatyperepo.PypiProxyRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Pypi: &sonatyperepo.PyPiProxyAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	if m.Cleanup != nil {
		mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	}

	// Proxy
	apiModel.HttpClient = sonatyperepo.HttpClientAttributes{}
	m.HttpClient.MapToApiHttpClientAttributes(&apiModel.HttpClient)
	m.NegativeCache.MapToApi(&apiModel.NegativeCache)
	m.Proxy.MapToApi(&apiModel.Proxy)
	if m.Replication != nil {
		apiModel.Replication = &sonatyperepo.ReplicationAttributes{}
		m.Replication.MapToApi(apiModel.Replication)
	}

	apiModel.RoutingRule = m.RoutingRule.ValueStringPointer()

	return apiModel
}

func (m *RepositoryPyPiProxyModel) ToApiUpdateModel() sonatyperepo.PypiProxyRepositoryApiRequest {
	model := m.ToApiCreateModel()
	if m.FirewallAuditAndQuarantine != nil && !m.FirewallAuditAndQuarantine.PccsEnabled.IsNull() && m.FirewallAuditAndQuarantine.PccsEnabled.ValueBool() {
		if model.Pypi == nil {
			model.Pypi = &sonatyperepo.PyPiProxyAttributes{}
		}
		model.Pypi.RemoveQuarantined = m.FirewallAuditAndQuarantine.PccsEnabled.ValueBool()
	}
	return model
}

// Group PyPI
// --------------------------------------------
type RepositoryPyPiGroupModel struct {
	RepositoryGroupDeployModel
}

func (m *RepositoryPyPiGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupDeployRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Storage
	m.Storage.MapFromApi(&api.Storage)

	// Group Attributes
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryPyPiGroupModel) ToApiCreateModel() sonatyperepo.PypiGroupRepositoryApiRequest {
	apiModel := sonatyperepo.PypiGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)

	// Group
	m.Group.MapToApi(&apiModel.Group)

	return apiModel
}

func (m *RepositoryPyPiGroupModel) ToApiUpdateModel() sonatyperepo.PypiGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}
