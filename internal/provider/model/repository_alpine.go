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

// Alpine Hosted
// ----------------------------------------
type RepositoryAlpineHostedModel struct {
	RepositoryHostedModel
	Alpine *alpineSigningModel `tfsdk:"alpine"`
}

func (m *RepositoryAlpineHostedModel) MapMissingApiFieldsFromPlan(planModel RepositoryAlpineHostedModel) {
	// Alpine signing fields are not returned by the GET API; preserve from plan
	m.Alpine = planModel.Alpine
}

func (m *RepositoryAlpineHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositoryAlpineHostedModel) ToApiCreateModel() sonatyperepo.AlpineHostedRepositoryApiRequest {
	apiModel := sonatyperepo.AlpineHostedRepositoryApiRequest{
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

	// Alpine
	if m.Alpine != nil {
		apiModel.AlpineSigning = &sonatyperepo.AlpineSigningRepositoriesAttributes{}
		m.Alpine.MapToApi(apiModel.AlpineSigning)
	}

	return apiModel
}

func (m *RepositoryAlpineHostedModel) ToApiUpdateModel() sonatyperepo.AlpineHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Alpine Proxy
// ----------------------------------------
type RepositoryAlpineProxyModel struct {
	RepositoryProxyModel
	Alpine                     *alpineSigningModel              `tfsdk:"alpine"`
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryAlpineProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryAlpineProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)
	// Alpine signing fields are not returned by the GET API; preserve from plan
	m.Alpine = planModel.Alpine
}

func (m *RepositoryAlpineProxyModel) FromApiModel(api sonatyperepo.SimpleApiProxyRepository) {
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

func (m *RepositoryAlpineProxyModel) ToApiCreateModel() sonatyperepo.AlpineProxyRepositoryApiRequest {
	apiModel := sonatyperepo.AlpineProxyRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Proxy:         sonatyperepo.ProxyAttributes{},
		NegativeCache: sonatyperepo.NegativeCacheAttributes{},
		HttpClient:    sonatyperepo.HttpClientAttributes{},
		AlpineSigning: sonatyperepo.AlpineSigningRepositoriesAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)

	if m.Cleanup != nil {
		mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	}

	// Proxy Specific
	m.Proxy.MapToApi(&apiModel.Proxy)
	m.NegativeCache.MapToApi(&apiModel.NegativeCache)
	m.HttpClient.MapToApiHttpClientAttributes(&apiModel.HttpClient)

	if m.Replication != nil {
		apiModel.Replication = &sonatyperepo.ReplicationAttributes{}
		m.Replication.MapToApi(apiModel.Replication)
	}

	apiModel.RoutingRule = m.RoutingRule.ValueStringPointer()

	// Alpine Specific
	if m.Alpine != nil {
		m.Alpine.MapToApi(&apiModel.AlpineSigning)
	}

	return apiModel
}

func (m *RepositoryAlpineProxyModel) ToApiUpdateModel() sonatyperepo.AlpineProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Alpine Group
// ----------------------------------------
type RepositoryAlpineGroupModel struct {
	RepositoryGroupModel
	Alpine *alpineSigningModel `tfsdk:"alpine"`
}

func (m *RepositoryAlpineGroupModel) MapMissingApiFieldsFromPlan(planModel RepositoryAlpineGroupModel) {
	// Alpine signing fields are not returned by the GET API; preserve from plan
	m.Alpine = planModel.Alpine
}

func (m *RepositoryAlpineGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryAlpineGroupModel) ToApiCreateModel() sonatyperepo.AlpineGroupRepositoryApiRequest {
	apiModel := sonatyperepo.AlpineGroupRepositoryApiRequest{
		Name:          m.Name.ValueString(),
		Online:        m.Online.ValueBool(),
		Storage:       sonatyperepo.StorageAttributes{},
		AlpineSigning: sonatyperepo.AlpineSigningRepositoriesAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)

	// Alpine Specific
	if m.Alpine != nil {
		m.Alpine.MapToApi(&apiModel.AlpineSigning)
	}

	return apiModel
}

func (m *RepositoryAlpineGroupModel) ToApiUpdateModel() sonatyperepo.AlpineGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// alpineSigningModel
// ----------------------------------------
type alpineSigningModel struct {
	KeyPair    types.String `tfsdk:"key_pair"`
	Passphrase types.String `tfsdk:"passphrase"`
}

func (m *alpineSigningModel) MapFromApi(api *sonatyperepo.AlpineSigningRepositoriesAttributes) {
	m.KeyPair = types.StringPointerValue(api.Keypair)
	// m.Passphrase = types.StringPointerValue(api.Passphrase)
}

func (m *alpineSigningModel) MapToApi(api *sonatyperepo.AlpineSigningRepositoriesAttributes) {
	api.Keypair = m.KeyPair.ValueStringPointer()
	api.Passphrase = m.Passphrase.ValueStringPointer()
}
