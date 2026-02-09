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

// Terraform Specific Model
type repositoryTerraformSpecificModel struct {
	RequireAuthentication types.Bool `tfsdk:"require_authentication"`
}

func (m *repositoryTerraformSpecificModel) FromApiModel(api *sonatyperepo.TerraformAttributes) {
	if api != nil {
		m.RequireAuthentication = types.BoolPointerValue(api.RequireAuthentication)
	}
}

func (m *repositoryTerraformSpecificModel) MapToApi(api *sonatyperepo.TerraformAttributes) {
	api.RequireAuthentication = m.RequireAuthentication.ValueBoolPointer()
}

// Terraform Proxy
// ----------------------------------------
type RepositoryTerraformProxyModel struct {
	RepositoryProxyModel
	Terraform repositoryTerraformSpecificModel `tfsdk:"terraform"`
}

func (m *RepositoryTerraformProxyModel) FromApiModel(api sonatyperepo.TerraformProxyApiRepository) {
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

	// Terraform Specific
	if api.Terraform != nil {
		m.Terraform.FromApiModel(api.Terraform)
	}
}

func (m *RepositoryTerraformProxyModel) ToApiCreateModel() sonatyperepo.TerraformProxyRepositoryApiRequest {
	apiModel := sonatyperepo.TerraformProxyRepositoryApiRequest{
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

	// Terraform Specific
	apiModel.Terraform = &sonatyperepo.TerraformAttributes{}
	m.Terraform.MapToApi(apiModel.Terraform)

	return apiModel
}

func (m *RepositoryTerraformProxyModel) ToApiUpdateModel() sonatyperepo.TerraformProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}
