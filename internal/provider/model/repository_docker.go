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
// type npmSpecificProxyModel struct {
// 	RemoveQuarrantined types.Bool `tfsdk:"remove_quarrantined"`
// }

// func (m *npmSpecificProxyModel) MapFromApi(api *sonatyperepo.NpmAttributes) {
// 	m.RemoveQuarrantined = types.BoolValue(api.RemoveQuarantined)
// }

// func (m *npmSpecificProxyModel) MapToApi(api *sonatyperepo.NpmAttributes) {
// 	api.RemoveQuarantined = m.RemoveQuarrantined.ValueBool()
// }

// type RepositoryDockerProxyModel struct {
// 	RepositoryProxyModel
// 	// Npm *npmSpecificProxyModel `tfsdk:"npm"`
// }

// func (m *RepositoryDockerProxyModel) FromApiModel(api sonatyperepo.DockerHostedApiRepository) {
// 	m.Name = types.StringPointerValue(api.Name)
// 	m.Online = types.BoolValue(api.Online)
// 	m.Url = types.StringPointerValue(api.Url)

// 	// Cleanup
// 	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
// 		m.Cleanup = NewRepositoryCleanupModel()
// 		mapCleanupFromApi(api.Cleanup, m.Cleanup)
// 	}

// 	// Storage
// 	m.Storage.MapFromApi(&api.Storage)

// 	// Docker Specific
// 	m.Docker.MapFromApi(api.Docker)
// 	// m.Npm.MapFromApi(api.Npm)
// }

// func (m *RepositoryDockerProxyModel) ToApiCreateModel() sonatyperepo.NpmProxyRepositoryApiRequest {
// 	apiModel := sonatyperepo.NpmProxyRepositoryApiRequest{
// 		Name:    m.Name.ValueString(),
// 		Online:  m.Online.ValueBool(),
// 		Storage: sonatyperepo.StorageAttributes{},
// 		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
// 			PolicyNames: make([]string, 0),
// 		},
// 	}
// 	m.Storage.MapToApi(&apiModel.Storage)

// 	if m.Cleanup != nil {
// 		mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
// 	}

// 	// Proxy Specific
// 	apiModel.Proxy = sonatyperepo.ProxyAttributes{}
// 	m.Proxy.MapToApi(&apiModel.Proxy)

// 	apiModel.NegativeCache = sonatyperepo.NegativeCacheAttributes{}
// 	m.NegativeCache.MapToApi(&apiModel.NegativeCache)

// 	apiModel.HttpClient = sonatyperepo.HttpClientAttributes{}
// 	m.HttpClient.MapToApiHttpClientAttributes(&apiModel.HttpClient)

// 	if m.Replication != nil {
// 		apiModel.Replication = &sonatyperepo.ReplicationAttributes{}
// 		m.Replication.MapToApi(apiModel.Replication)
// 	}

// 	apiModel.RoutingRule = m.RoutingRule.ValueStringPointer()

// 	// NPM Specific
// 	if m.Npm != nil {
// 		apiModel.Npm = &sonatyperepo.NpmAttributes{}
// 		m.Npm.MapToApi(apiModel.Npm)
// 	}

// 	return apiModel
// }

// func (m *RepositoryDockerProxyModel) ToApiUpdateModel() sonatyperepo.NpmProxyRepositoryApiRequest {
// 	return m.ToApiCreateModel()
// }

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
