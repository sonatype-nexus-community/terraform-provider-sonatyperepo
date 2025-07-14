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

// APT Hosted
// ----------------------------------------
// type RepositoryAptHostedModel struct {
// 	RepositoryHostedModel
// }

// func (m *RepositoryAptHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
// 	m.mapSimpleApiHostedRepository(api)
// }

// func (m *RepositoryAptHostedModel) ToApiCreateModel() sonatyperepo.NpmHostedRepositoryApiRequest {
// 	apiModel := sonatyperepo.NpmHostedRepositoryApiRequest{
// 		Name:    m.Name.ValueString(),
// 		Online:  m.Online.ValueBool(),
// 		Storage: sonatyperepo.HostedStorageAttributes{},
// 		Component: &sonatyperepo.ComponentAttributes{
// 			ProprietaryComponents: common.NewFalse(),
// 		},
// 		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
// 			PolicyNames: make([]string, 0),
// 		},
// 	}
// 	m.Storage.MapToApi(&apiModel.Storage)
// 	mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
// 	m.Component.MapToApi(apiModel.Component)
// 	return apiModel
// }

// func (m *RepositoryAptHostedModel) ToApiUpdateModel() sonatyperepo.NpmHostedRepositoryApiRequest {
// 	return m.ToApiCreateModel()
// }

// APT Proxy
// ----------------------------------------
type aptSpecificProxyModel struct {
	Distribution types.String `tfsdk:"distribution"`
	Flat         types.Bool   `tfsdk:"flat"`
}

func (m *aptSpecificProxyModel) MapFromApi(api *sonatyperepo.AptProxyRepositoriesAttributes) {
	m.Distribution = types.StringPointerValue(api.Distribution)
	m.Flat = types.BoolPointerValue(api.Flat)
}

func (m *aptSpecificProxyModel) MapToApi(api *sonatyperepo.AptProxyRepositoriesAttributes) {
	api.Distribution = m.Distribution.ValueStringPointer()
	api.Flat = m.Flat.ValueBoolPointer()
}

type RepositoryAptProxyModel struct {
	RepositoryProxyModel
	Apt aptSpecificProxyModel `tfsdk:"apt"`
}

func (m *RepositoryAptProxyModel) FromApiModel(api sonatyperepo.AptProxyApiRepository) {
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

	// APT Specific
	m.Apt.MapFromApi(&api.Apt)
}

func (m *RepositoryAptProxyModel) ToApiCreateModel() sonatyperepo.AptProxyRepositoryApiRequest {
	apiModel := sonatyperepo.AptProxyRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Proxy:         sonatyperepo.ProxyAttributes{},
		NegativeCache: sonatyperepo.NegativeCacheAttributes{},
		HttpClient:    sonatyperepo.HttpClientAttributes{},
		Apt:           sonatyperepo.AptProxyRepositoriesAttributes{},
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

	// APT Specific
	m.Apt.MapToApi(&apiModel.Apt)

	return apiModel
}

func (m *RepositoryAptProxyModel) ToApiUpdateModel() sonatyperepo.AptProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}
