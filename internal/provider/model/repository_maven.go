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

type repositoryMavenSpecificModel struct {
	VersionPolicy      types.String `tfsdk:"version_policy"`
	LayoutPolicy       types.String `tfsdk:"layout_policy"`
	ContentDisposition types.String `tfsdk:"content_disposition"`
}

// Hosted Maven
// --------------------------------------------
type RepositoryMavenHostedModel struct {
	RepositoryHostedModel
	Maven repositoryMavenSpecificModel `tfsdk:"maven"`
}

func (m *RepositoryMavenHostedModel) FromApiModel(api sonatyperepo.MavenHostedApiRepository) {
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

	// Component
	if api.Component != nil {
		m.Component = &RepositoryComponentModel{}
		m.Component.MapFromApi(api.Component)
	}

	// Maven Specific
	m.Maven = repositoryMavenSpecificModel{}
	m.Maven.MapFromApi(&api.Maven)
}

func (m *RepositoryMavenHostedModel) ToApiCreateModel() sonatyperepo.MavenHostedRepositoryApiRequest {
	apiModel := sonatyperepo.MavenHostedRepositoryApiRequest{
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

	// Maven Specific
	m.Maven.MapToApi(&apiModel.Maven)

	return apiModel
}

func (m *RepositoryMavenHostedModel) ToApiUpdateModel() sonatyperepo.MavenHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Proxy Maven
// --------------------------------------------
type RepositoryMavenProxyModel struct {
	RepositoryProxyModel
	Maven                      repositoryMavenSpecificModel     `tfsdk:"maven"`
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryMavenProxyModel) FromApiModel(api sonatyperepo.MavenProxyApiRepository) {
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
		m.Replication = &RepositoryReplicationModel{}
		m.Replication.MapFromApi(api.Replication)
	} else {
		// Set default values when API doesn't provide replication data
		m.Replication = &RepositoryReplicationModel{
			PreemptivePullEnabled: types.BoolValue(common.DEFAULT_PROXY_PREEMPTIVE_PULL),
			AssetPathRegex:        types.StringNull(),
		}
	}

	// Maven Specific
	m.Maven.MapFromApi(&api.Maven)
}

func (m *RepositoryMavenProxyModel) ToApiCreateModel() sonatyperepo.MavenProxyRepositoryApiRequest {
	apiModel := sonatyperepo.MavenProxyRepositoryApiRequest{
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

	// Proxy
	apiModel.HttpClient = sonatyperepo.HttpClientAttributesWithPreemptiveAuth{}
	m.HttpClient.MapToApiHttpClientAttributesWithPreemptiveAuth(&apiModel.HttpClient)
	m.NegativeCache.MapToApi(&apiModel.NegativeCache)
	m.Proxy.MapToApi(&apiModel.Proxy)
	if m.Replication != nil {
		apiModel.Replication = &sonatyperepo.ReplicationAttributes{}
		m.Replication.MapToApi(apiModel.Replication)
	}

	apiModel.RoutingRule = m.RoutingRule.ValueStringPointer()

	// Maven
	m.Maven.MapToApi(&apiModel.Maven)

	return apiModel
}

func (m *RepositoryMavenProxyModel) ToApiUpdateModel() sonatyperepo.MavenProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

func (m *repositoryMavenSpecificModel) MapFromApi(api *sonatyperepo.MavenAttributes) {
	m.ContentDisposition = types.StringPointerValue(api.ContentDisposition)
	m.LayoutPolicy = types.StringPointerValue(api.LayoutPolicy)
	m.VersionPolicy = types.StringPointerValue(api.VersionPolicy)
}

func (m *repositoryMavenSpecificModel) MapToApi(api *sonatyperepo.MavenAttributes) {
	api.ContentDisposition = m.ContentDisposition.ValueStringPointer()
	api.LayoutPolicy = m.LayoutPolicy.ValueStringPointer()
	api.VersionPolicy = m.VersionPolicy.ValueStringPointer()
}

// Group Maven
// --------------------------------------------
type RepositoryMavenGroupModel struct {
	RepositoryGroupModel
}

func (m *RepositoryMavenGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Storage
	m.Storage.MapFromApi(&api.Storage)

	// Group Attributes
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryMavenGroupModel) ToApiCreateModel() sonatyperepo.MavenGroupRepositoryApiRequest {
	apiModel := sonatyperepo.MavenGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)

	// Group
	m.Group.MapToApi(&apiModel.Group)

	// Maven - Injected to keep NXRM 3.88 happy (they are in API, but not used)
	// NXRM 3.89.0 removed these from the API...
	// apiModel.Maven = *sonatyperepo.NewMavenAttributesWithDefaults()
	// apiModel.Maven.ContentDisposition = common.StringPointer(common.MAVEN_CONTENT_DISPOSITION_INLINE)
	// apiModel.Maven.LayoutPolicy = common.StringPointer(common.MAVEN_LAYOUT_PERMISSIVE)
	// apiModel.Maven.VersionPolicy = common.StringPointer(common.MAVEN_VERSION_POLICY_MIXED)

	return apiModel
}

func (m *RepositoryMavenGroupModel) ToApiUpdateModel() sonatyperepo.MavenGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}
