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

// Conan Hosted
// ----------------------------------------
type RepositorConanHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositorConanHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositorConanHostedModel) ToApiCreateModel() sonatyperepo.ConanHostedRepositoryApiRequest {
	apiModel := sonatyperepo.ConanHostedRepositoryApiRequest{
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

func (m *RepositorConanHostedModel) ToApiUpdateModel() sonatyperepo.ConanHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Conan Proxy
// ----------------------------------------
type RepositoryConanProxyModel struct {
	RepositoryProxyModel
	Conan                      *conanProxyAttributesModel       `tfsdk:"conan"`
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryConanProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryConanProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)
}

func (m *RepositoryConanProxyModel) FromApiModel(api sonatyperepo.ConanProxyApiRepository) {
	m.Name = types.StringValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringValue(api.Url)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = NewRepositoryCleanupModel()
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
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

	// Conan
	if api.ConanProxy != nil {
		m.Conan = &conanProxyAttributesModel{}
		m.Conan.MapFromApi(api.ConanProxy)
	}
}

func (m *RepositoryConanProxyModel) ToApiCreateModel() sonatyperepo.ConanProxyRepositoryApiRequest {
	apiModel := sonatyperepo.ConanProxyRepositoryApiRequest{
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

	// Conan Specific
	if m.Conan != nil {
		apiModel.ConanProxy = &sonatyperepo.ConanProxyAttributes{}
		m.Conan.MapToApi(apiModel.ConanProxy)
	}

	return apiModel
}

func (m *RepositoryConanProxyModel) ToApiUpdateModel() sonatyperepo.ConanProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Conan Group
// ----------------------------------------
type RepositoryConanGroupModel struct {
	RepositoryGroupDeployModel
}

func (m *RepositoryConanGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupDeployRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryConanGroupModel) ToApiCreateModel() sonatyperepo.ConanGroupRepositoryApiRequest {
	apiModel := sonatyperepo.ConanGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)
	return apiModel
}

func (m *RepositoryConanGroupModel) ToApiUpdateModel() sonatyperepo.ConanGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// conanProxyAttributesModel
type conanProxyAttributesModel struct {
	ConanVersion types.String `tfsdk:"conan_version"`
}

func (m *conanProxyAttributesModel) MapFromApi(api *sonatyperepo.ConanProxyAttributes) {
	m.ConanVersion = types.StringPointerValue(api.ConanVersion)
}

func (m *conanProxyAttributesModel) MapToApi(api *sonatyperepo.ConanProxyAttributes) {
	api.ConanVersion = m.ConanVersion.ValueStringPointer()
}
