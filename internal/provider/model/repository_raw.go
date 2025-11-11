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

// Raw Hosted
// ----------------------------------------
type RepositoryRawHostedModel struct {
	RepositoryHostedModel
	Raw RawRepositoryAttributesModel `tfsdk:"raw"`
}

func (m *RepositoryRawHostedModel) FromApiModel(api sonatyperepo.RawHostedApiRepository) {
	m.Name = types.StringValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = &RepositoryCleanupModel{}
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
	}

	// Component
	if api.Component != nil {
		m.Component = &RepositoryComponentModel{}
		m.Component.MapFromApi(api.Component)
	}

	// Raw Specific
	m.Raw.MapFromApi(&api.Raw)
}

func (m *RepositoryRawHostedModel) ToApiCreateModel() sonatyperepo.RawHostedRepositoryApiRequest {
	apiModel := sonatyperepo.RawHostedRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{},
		Component: &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: common.NewFalse(),
		},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Raw: sonatyperepo.NewRawAttributes(),
	}
	m.Storage.MapToApi(&apiModel.Storage)
	mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	m.Component.MapToApi(apiModel.Component)

	// Raw
	m.Raw.MapToApi(apiModel.Raw)

	return apiModel
}

func (m *RepositoryRawHostedModel) ToApiUpdateModel() sonatyperepo.RawHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Raw Proxy
// ----------------------------------------
type RepositoryRawProxyModel struct {
	RepositoryProxyModel
	Raw RawRepositoryAttributesModel `tfsdk:"raw"`
}

func (m *RepositoryRawProxyModel) FromApiModel(api sonatyperepo.RawProxyApiRepository) {
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
			PreemptivePullEnabled: types.BoolValue(common.DEFAULT_PROXY_PREEMPTIVE_PULL),
			AssetPathRegex:        types.StringNull(),
		}
	}

	// Raw Specific
	m.Raw.MapFromApi(&api.Raw)
}

func (m *RepositoryRawProxyModel) ToApiCreateModel() sonatyperepo.RawProxyRepositoryApiRequest {
	apiModel := sonatyperepo.RawProxyRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Raw: sonatyperepo.NewRawAttributes(),
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

	// Raw Specific
	m.Raw.MapToApi(apiModel.Raw)

	return apiModel
}

func (m *RepositoryRawProxyModel) ToApiUpdateModel() sonatyperepo.RawProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Raw Group
// ----------------------------------------
type RepositoryRawGroupModel struct {
	RepositoryGroupModel
	Raw RawRepositoryAttributesModel `tfsdk:"raw"`
}

func (m *RepositoryRawGroupModel) FromApiModel(api sonatyperepo.RawGroupApiRepository) {
	m.Name = types.StringValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)
	m.Group.MapFromApi(&api.Group)
	m.Raw.MapFromApi(&api.Raw)
}

func (m *RepositoryRawGroupModel) ToApiCreateModel() sonatyperepo.RawGroupRepositoryApiRequest {
	apiModel := sonatyperepo.RawGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Raw:     sonatyperepo.NewRawAttributes(),
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)
	m.Raw.MapToApi(apiModel.Raw)
	return apiModel
}

func (m *RepositoryRawGroupModel) ToApiUpdateModel() sonatyperepo.RawGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// RawRepositoryAttributesModel
// ----------------------------------------
type RawRepositoryAttributesModel struct {
	ContentDisposition types.String `tfsdk:"content_disposition"`
}

func (m *RawRepositoryAttributesModel) MapFromApi(api *sonatyperepo.RawAttributes) {
	m.ContentDisposition = types.StringPointerValue(api.ContentDisposition)
}

func (m *RawRepositoryAttributesModel) MapToApi(api *sonatyperepo.RawAttributes) {
	api.ContentDisposition = m.ContentDisposition.ValueStringPointer()
}
