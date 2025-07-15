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

// Nuget Hosted
// ----------------------------------------
type RepositoryNugetHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositoryNugetHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositoryNugetHostedModel) ToApiCreateModel() sonatyperepo.NugetHostedRepositoryApiRequest {
	apiModel := sonatyperepo.NugetHostedRepositoryApiRequest{
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

func (m *RepositoryNugetHostedModel) ToApiUpdateModel() sonatyperepo.NugetHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Nuget Proxy
// ----------------------------------------
type RepositoryNugetProxyModel struct {
	RepositoryProxyModel
	NugetProxy *NugetProxyModel `tfsdk:"nuget_proxy"`
}

func (m *RepositoryNugetProxyModel) FromApiModel(api sonatyperepo.NugetProxyApiRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

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
		m.Replication.MapFromApi(api.Replication)
	}

	// Nuget Specific
	m.NugetProxy.MapFromApi(api.NugetProxy)
}

func (m *RepositoryNugetProxyModel) ToApiCreateModel() sonatyperepo.NugetProxyRepositoryApiRequest {
	apiModel := sonatyperepo.NugetProxyRepositoryApiRequest{
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

	// Nuget Specific
	if m.NugetProxy != nil {
		apiModel.NugetProxy = sonatyperepo.NugetAttributes{}
		m.NugetProxy.MapToApi(&apiModel.NugetProxy)
	}

	return apiModel
}

func (m *RepositoryNugetProxyModel) ToApiUpdateModel() sonatyperepo.NugetProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Nuget Group
// ----------------------------------------
type RepositoryNugetGroupModel struct {
	RepositoryGroupModel
}

func (m *RepositoryNugetGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryNugetGroupModel) ToApiCreateModel() sonatyperepo.NugetGroupRepositoryApiRequest {
	apiModel := sonatyperepo.NugetGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)
	return apiModel
}

func (m *RepositoryNugetGroupModel) ToApiUpdateModel() sonatyperepo.NugetGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// NugetProxyModel
// ----------------------------------------
type NugetProxyModel struct {
	NugetVersion         types.String `tfsdk:"nuget_version"`
	QueryCacheItemMaxAge types.Int32  `tfsdk:"query_cache_item_max_age"`
}

func (m *NugetProxyModel) MapFromApi(api sonatyperepo.NugetAttributes) {
	m.NugetVersion = types.StringPointerValue(api.NugetVersion)
	m.QueryCacheItemMaxAge = types.Int32PointerValue(api.QueryCacheItemMaxAge)
}

func (m *NugetProxyModel) MapToApi(api *sonatyperepo.NugetAttributes) {
	api.NugetVersion = m.NugetVersion.ValueStringPointer()
	api.QueryCacheItemMaxAge = m.QueryCacheItemMaxAge.ValueInt32Pointer()
}
