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

// NPM Hosted
type RepositoryNpmHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositoryNpmHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositoryNpmHostedModel) ToApiCreateModel() sonatyperepo.NpmHostedRepositoryApiRequest {
	apiModel := sonatyperepo.NpmHostedRepositoryApiRequest{
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

func (m *RepositoryNpmHostedModel) ToApiUpdateModel() sonatyperepo.NpmHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// NPM Proxy Specific Attributes
// ----------------------------------------
type npmSpecificProxyModel struct {
	RemoveQuarrantined types.Bool `tfsdk:"remove_quarrantined"`
}

func (m *npmSpecificProxyModel) MapFromApi(api *sonatyperepo.NpmAttributes) {
	m.RemoveQuarrantined = types.BoolValue(api.RemoveQuarantined)
}

func (m *npmSpecificProxyModel) MapToApi(api *sonatyperepo.NpmAttributes) {
	api.RemoveQuarantined = m.RemoveQuarrantined.ValueBool()
}

type RepositoryNpmProxyModel struct {
	RepositoryProxyModel
	Npm *npmSpecificProxyModel `tfsdk:"npm"`
}

func (m *RepositoryNpmProxyModel) FromApiModel(api sonatyperepo.NpmProxyApiRepository) {
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

	// NPM Specific
	m.Npm.MapFromApi(api.Npm)
}

func (m *RepositoryNpmProxyModel) ToApiCreateModel() sonatyperepo.NpmProxyRepositoryApiRequest {
	apiModel := sonatyperepo.NpmProxyRepositoryApiRequest{
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

	// NPM Specific
	if m.Npm != nil {
		apiModel.Npm = &sonatyperepo.NpmAttributes{}
		m.Npm.MapToApi(apiModel.Npm)
	}

	return apiModel
}

func (m *RepositoryNpmProxyModel) ToApiUpdateModel() sonatyperepo.NpmProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}
