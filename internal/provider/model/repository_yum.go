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

// YUM Hosted
// ----------------------------------------
type RepositoryYumHostedModel struct {
	RepositoryHostedModel
	Yum yumHostedAttributesModel `tfsdk:"yum"`
}

func (m *RepositoryYumHostedModel) FromApiModel(api sonatyperepo.YumHostedApiRepository) {
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

	// YUM Specific
	m.Yum.MapFromApi(&api.Yum)
}

func (m *RepositoryYumHostedModel) ToApiCreateModel() sonatyperepo.YumHostedRepositoryApiRequest {
	apiModel := sonatyperepo.YumHostedRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{},
		Component: &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: common.NewFalse(),
		},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Yum: *sonatyperepo.NewYumAttributesWithDefaults(),
	}
	m.Storage.MapToApi(&apiModel.Storage)
	mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	m.Component.MapToApi(apiModel.Component)

	// YUM
	m.Yum.MapToApi(&apiModel.Yum)

	return apiModel
}

func (m *RepositoryYumHostedModel) ToApiUpdateModel() sonatyperepo.YumHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// YUM Proxy
// ----------------------------------------
type RepositoryYumProxyModel struct {
	RepositoryProxyModel
	Yum                        *yumSigningModel                 `tfsdk:"yum"`
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryYumProxyModel) FromApiModel(api sonatyperepo.SimpleApiProxyRepository) {
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

	// YUM Specific
	// NOT RETURNED BY GET API - yum field is optional and not populated during import

	// Firewall Audit and Quarantine
	// This will be populated separately by the resource helper during Read operations
	if m.FirewallAuditAndQuarantine == nil {
		m.FirewallAuditAndQuarantine = NewFirewallAuditAndQuarantineModelWithDefaults()
	}
}

func (m *RepositoryYumProxyModel) ToApiCreateModel() sonatyperepo.YumProxyRepositoryApiRequest {
	apiModel := sonatyperepo.YumProxyRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
		Proxy:         sonatyperepo.ProxyAttributes{},
		NegativeCache: sonatyperepo.NegativeCacheAttributes{},
		HttpClient:    sonatyperepo.HttpClientAttributes{},
		YumSigning:    &sonatyperepo.YumSigningRepositoriesAttributes{},
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

	// YUM Specific
	if m.Yum != nil {
		m.Yum.MapToApi(apiModel.YumSigning)
	}

	return apiModel
}

func (m *RepositoryYumProxyModel) ToApiUpdateModel() sonatyperepo.YumProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// YUM Group
// ----------------------------------------
type RepositoryYumGroupModel struct {
	RepositoryGroupModel
	Yum *yumSigningModel `tfsdk:"yum"`
}

func (m *RepositoryYumGroupModel) FromApiModel(api sonatyperepo.SimpleApiGroupRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Storage
	m.Storage.MapFromApi(&api.Storage)

	// Group Attributes
	m.Group.MapFromApi(&api.Group)

	// Yum
	// NOT RETURNED BY GET API
}

func (m *RepositoryYumGroupModel) ToApiCreateModel() sonatyperepo.YumGroupRepositoryApiRequest {
	apiModel := sonatyperepo.YumGroupRepositoryApiRequest{
		Name:       m.Name.ValueString(),
		Online:     m.Online.ValueBool(),
		Storage:    sonatyperepo.StorageAttributes{},
		YumSigning: sonatyperepo.NewYumSigningRepositoriesAttributesWithDefaults(),
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)

	// YUM
	if m.Yum != nil {
		m.Yum.MapToApi(apiModel.YumSigning)
	}

	return apiModel
}

func (m *RepositoryYumGroupModel) ToApiUpdateModel() sonatyperepo.YumGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// yumHostedAttributesModel
// ----------------------------------------
type yumHostedAttributesModel struct {
	DeployPolicy  types.String `tfsdk:"deploy_policy"`
	RepoDataDepth types.Int32  `tfsdk:"repo_data_depth"`
}

func (m *yumHostedAttributesModel) MapFromApi(api *sonatyperepo.YumAttributes) {
	m.DeployPolicy = types.StringPointerValue(api.DeployPolicy)
	m.RepoDataDepth = types.Int32Value(api.RepodataDepth)
}

func (m *yumHostedAttributesModel) MapToApi(api *sonatyperepo.YumAttributes) {
	api.DeployPolicy = m.DeployPolicy.ValueStringPointer()
	api.RepodataDepth = m.RepoDataDepth.ValueInt32()
}

// yumSigningAttributesModel
// ----------------------------------------
type yumSigningModel struct {
	KeyPair    types.String `tfsdk:"key_pair"`
	Passphrase types.String `tfsdk:"passphrase"`
}

func (m *yumSigningModel) MapFromApi(api *sonatyperepo.YumSigningRepositoriesAttributes) {
	m.KeyPair = types.StringPointerValue(api.Keypair)
	// m.Passphrase = types.StringPointerValue(api.Passphrase)
}

func (m *yumSigningModel) MapToApi(api *sonatyperepo.YumSigningRepositoriesAttributes) {
	api.Keypair = m.KeyPair.ValueStringPointer()
	api.Passphrase = m.Passphrase.ValueStringPointer()
}
