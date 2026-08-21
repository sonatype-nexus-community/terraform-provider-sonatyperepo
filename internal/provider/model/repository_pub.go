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

// Pub Proxy
// ----------------------------------------
type RepositoryPubProxyModel struct {
	RepositoryProxyModel
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryPubProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryPubProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)
}

func (m *RepositoryPubProxyModel) FromApiModel(api sonatyperepo.SimpleApiProxyRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

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
}

func (m *RepositoryPubProxyModel) ToApiCreateModel() sonatyperepo.PubProxyRepositoryApiRequest {
	apiModel := sonatyperepo.PubProxyRepositoryApiRequest{
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

func (m *RepositoryPubProxyModel) ToApiUpdateModel() sonatyperepo.PubProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Pub Group
// ----------------------------------------
type RepositoryPubGroupModel struct {
	RepositoryGroupModel
}

func (m *RepositoryPubGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryPubGroupModel) ToApiCreateModel() sonatyperepo.PubGroupRepositoryApiRequest {
	apiModel := sonatyperepo.PubGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)
	return apiModel
}

func (m *RepositoryPubGroupModel) ToApiUpdateModel() sonatyperepo.PubGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Pub Hosted
// ----------------------------------------
type RepositoryPubHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositoryPubHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = &RepositoryCleanupModel{}
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
	} else {
		m.Cleanup = nil
	}

	// Component
	if api.Component != nil {
		m.Component = &RepositoryComponentModel{}
		m.Component.MapFromApi(api.Component)
	}
}

func (m *RepositoryPubHostedModel) ToApiCreateModel() sonatyperepo.PubHostedRepositoryApiRequest {
	apiModel := sonatyperepo.PubHostedRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)

	if m.Cleanup != nil {
		apiModel.Cleanup = &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		}
		mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	}

	if m.Component != nil {
		apiModel.Component = &sonatyperepo.ComponentAttributes{}
		m.Component.MapToApi(apiModel.Component)
	}

	return apiModel
}

func (m *RepositoryPubHostedModel) ToApiUpdateModel() sonatyperepo.PubHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}
