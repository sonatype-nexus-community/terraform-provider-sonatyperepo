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

// Cargo Hosted
// ----------------------------------------
type RepositorCargoHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositorCargoHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositorCargoHostedModel) ToApiCreateModel() sonatyperepo.CargoHostedRepositoryApiRequest {
	apiModel := sonatyperepo.CargoHostedRepositoryApiRequest{
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

func (m *RepositorCargoHostedModel) ToApiUpdateModel() sonatyperepo.CargoHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Cargo Proxy
// ----------------------------------------
type RepositoryCargoProxyModel struct {
	RepositoryProxyModel
	Cargo                      *cargoAttributesModel            `tfsdk:"cargo"`
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryCargoProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryCargoProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)
}

func (m *RepositoryCargoProxyModel) FromApiModel(api sonatyperepo.CargoProxyApiRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = &RepositoryCleanupModel{}
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
	} else {
		m.Cleanup = nil
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

	// Cargo Specific
	if api.Cargo != nil {
		m.Cargo = &cargoAttributesModel{}
		m.Cargo.MapFromApi(api.Cargo)
	}
}

func (m *RepositoryCargoProxyModel) ToApiCreateModel() sonatyperepo.CargoProxyRepositoryApiRequest {
	apiModel := sonatyperepo.CargoProxyRepositoryApiRequest{
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
	if m.Cargo != nil {
		apiModel.Cargo = &sonatyperepo.CargoAttributes{}
		m.Cargo.MapToApi(apiModel.Cargo)
	}

	return apiModel
}

func (m *RepositoryCargoProxyModel) ToApiUpdateModel() sonatyperepo.CargoProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// Cargo Group
// ----------------------------------------
type RepositoryCargoGroupModel struct {
	RepositoryGroupModel
	Cargo *cargoAttributesModel `tfsdk:"cargo"`
}

func (m *RepositoryCargoGroupModel) FromApiModel(api sonatyperepo.CargoGroupApiRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)
	m.Group.MapFromApi(&api.Group)
}

func (m *RepositoryCargoGroupModel) ToApiCreateModel() sonatyperepo.CargoGroupRepositoryApiRequest {
	apiModel := sonatyperepo.CargoGroupRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)
	if m.Cargo != nil {
		apiModel.Cargo = sonatyperepo.NewCargoAttributes()
		m.Cargo.MapToApi(apiModel.Cargo)
	}
	return apiModel
}

func (m *RepositoryCargoGroupModel) ToApiUpdateModel() sonatyperepo.CargoGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// CargoProxyAttributes
type cargoAttributesModel struct {
	RequiresAuthentication types.Bool `tfsdk:"require_authentication"`
}

func (m *cargoAttributesModel) MapFromApi(api *sonatyperepo.CargoAttributes) {
	m.RequiresAuthentication = types.BoolPointerValue(api.RequireAuthentication)
}

func (m *cargoAttributesModel) MapToApi(api *sonatyperepo.CargoAttributes) {
	api.RequireAuthentication = m.RequiresAuthentication.ValueBoolPointer()
}
