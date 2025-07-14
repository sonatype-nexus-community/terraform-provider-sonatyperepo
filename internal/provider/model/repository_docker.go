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

// dockerHostedStorageModel
// ----------------------------------------
type dockerHostedStorageModel struct {
	repositoryStorageModelHosted
	LatestPolicy types.Bool `tfsdk:"latest_policy"`
}

func (m *dockerHostedStorageModel) MapFromApi(api *sonatyperepo.DockerHostedStorageAttributes) {
	m.BlobStoreName = types.StringValue(api.BlobStoreName)
	m.StrictContentTypeValidation = types.BoolValue(api.StrictContentTypeValidation)
	m.WritePolicy = types.StringValue(api.WritePolicy)
	m.LatestPolicy = types.BoolPointerValue(api.LatestPolicy)
}

func (m *dockerHostedStorageModel) MapToApi(api *sonatyperepo.DockerHostedStorageAttributes) {
	api.BlobStoreName = m.BlobStoreName.ValueString()
	api.StrictContentTypeValidation = m.StrictContentTypeValidation.ValueBool()
	api.WritePolicy = m.WritePolicy.ValueString()
	api.LatestPolicy = m.LatestPolicy.ValueBoolPointer()
}

// dockerAttributesModel
// ----------------------------------------
type dockerAttributesModel struct {
	ForceBasicAuth types.Bool   `tfsdk:"force_basic_auth"`
	HttpPort       types.Int32  `tfsdk:"http_port"`
	HttpsPort      types.Int32  `tfsdk:"https_port"`
	Subdomain      types.String `tfsdk:"subdomain"`
	V1Enabled      types.Bool   `tfsdk:"v1_enabled"`
}

func (m *dockerAttributesModel) MapFromApi(api *sonatyperepo.DockerAttributes) {
	m.ForceBasicAuth = types.BoolValue(api.ForceBasicAuth)
	m.HttpPort = types.Int32PointerValue(api.HttpPort)
	m.HttpsPort = types.Int32PointerValue(api.HttpsPort)
	m.Subdomain = types.StringPointerValue(api.Subdomain)
	m.V1Enabled = types.BoolValue(api.V1Enabled)
}

func (m *dockerAttributesModel) MapToApi(api *sonatyperepo.DockerAttributes) {
	api.ForceBasicAuth = m.ForceBasicAuth.ValueBool()
	api.HttpPort = m.HttpPort.ValueInt32Pointer()
	api.HttpsPort = m.HttpsPort.ValueInt32Pointer()
	api.Subdomain = m.Subdomain.ValueStringPointer()
	api.V1Enabled = m.V1Enabled.ValueBool()
}

// Docker Hosted
// ----------------------------------------
type RepositoryDockerHostedModel struct {
	BasicRepositoryModel
	Storage   dockerHostedStorageModel  `tfsdk:"storage"`
	Component *RepositoryComponentModel `tfsdk:"component"`
	Docker    dockerAttributesModel     `tfsdk:"docker"`
}

func (m *RepositoryDockerHostedModel) FromApiModel(api sonatyperepo.DockerHostedApiRepository) {
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

	// Docker Specific
	m.Docker.MapFromApi(&api.Docker)
}

func (m *RepositoryDockerHostedModel) ToApiCreateModel() sonatyperepo.DockerHostedRepositoryApiRequest {
	apiModel := sonatyperepo.DockerHostedRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.DockerHostedStorageAttributes{},
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

	// Docker Specific
	m.Docker.MapToApi(&apiModel.Docker)

	return apiModel
}

func (m *RepositoryDockerHostedModel) ToApiUpdateModel() sonatyperepo.DockerHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Docker Proxy
// ----------------------------------------
type dockerProxyAttributesModel struct {
	CacheForeignLayers       types.Bool     `tfsdk:"cache_foreign_layers"`
	ForeignLayerUrlWhitelist []types.String `tfsdk:"foreign_layer_url_whitelist"`
	IndexType                types.String   `tfsdk:"index_type"`
	IndexUrl                 types.String   `tfsdk:"index_url"`
}

func (m *dockerProxyAttributesModel) MapFromApi(api *sonatyperepo.DockerProxyAttributes) {
	m.CacheForeignLayers = types.BoolPointerValue(api.CacheForeignLayers)
	m.ForeignLayerUrlWhitelist = make([]types.String, 0)
	for _, l := range api.ForeignLayerUrlWhitelist {
		m.ForeignLayerUrlWhitelist = append(m.ForeignLayerUrlWhitelist, types.StringValue(l))
	}
	m.IndexType = types.StringPointerValue(api.IndexType)
	m.IndexUrl = types.StringPointerValue(api.IndexUrl)
}

func (m *dockerProxyAttributesModel) MapToApi(api *sonatyperepo.DockerProxyAttributes) {
	api.CacheForeignLayers = m.CacheForeignLayers.ValueBoolPointer()
	api.ForeignLayerUrlWhitelist = make([]string, 0)
	for _, l := range m.ForeignLayerUrlWhitelist {
		api.ForeignLayerUrlWhitelist = append(api.ForeignLayerUrlWhitelist, l.ValueString())
	}
	api.IndexType = m.IndexType.ValueStringPointer()
	api.IndexUrl = m.IndexUrl.ValueStringPointer()
}

type RepositoryDockerProxyModel struct {
	RepositoryProxyModel
	Docker      dockerAttributesModel      `tfsdk:"docker"`
	DockerProxy dockerProxyAttributesModel `tfsdk:"docker_proxy"`
}

func (m *RepositoryDockerProxyModel) FromApiModel(api sonatyperepo.DockerProxyApiRepository) {
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

	// Docker Specific
	m.Docker.MapFromApi(&api.Docker)
	m.DockerProxy.MapFromApi(&api.DockerProxy)
}

func (m *RepositoryDockerProxyModel) ToApiCreateModel() sonatyperepo.DockerProxyRepositoryApiRequest {
	apiModel := sonatyperepo.DockerProxyRepositoryApiRequest{
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

	// Docker Specific
	m.Docker.MapToApi(&apiModel.Docker)
	m.DockerProxy.MapToApi(&apiModel.DockerProxy)

	return apiModel
}

func (m *RepositoryDockerProxyModel) ToApiUpdateModel() sonatyperepo.DockerProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// // Docker Group
// // ----------------------------------------
// type RepositoryDockerroupModel struct {
// 	RepositoryGroupDeployModel
// }

// func (m *RepositoryDockerroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupDeployRepository) {
// 	m.Name = types.StringPointerValue(api.Name)
// 	m.Online = types.BoolValue(api.Online)
// 	m.Url = types.StringPointerValue(api.Url)
// 	m.Storage.MapFromApi(&api.Storage)
// 	m.Group.MapFromApi(&api.Group)
// }

// func (m *RepositoryDockerroupModel) ToApiCreateModel() sonatyperepo.NpmGroupRepositoryApiRequest {
// 	apiModel := sonatyperepo.NpmGroupRepositoryApiRequest{
// 		Name:    m.Name.ValueString(),
// 		Online:  m.Online.ValueBool(),
// 		Storage: sonatyperepo.StorageAttributes{},
// 	}
// 	m.Storage.MapToApi(&apiModel.Storage)
// 	m.Group.MapToApi(&apiModel.Group)
// 	return apiModel
// }

// func (m *RepositoryDockerroupModel) ToApiUpdateModel() sonatyperepo.NpmGroupRepositoryApiRequest {
// 	return m.ToApiCreateModel()
// }
