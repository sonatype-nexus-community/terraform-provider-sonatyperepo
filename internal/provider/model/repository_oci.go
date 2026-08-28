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

// ociAttributesModel
// ----------------------------------------
type ociAttributesModel struct {
	ForceBasicAuth types.Bool   `tfsdk:"force_basic_auth"`
	HttpPort       types.Int32  `tfsdk:"http_port"`
	HttpsPort      types.Int32  `tfsdk:"https_port"`
	PathEnabled    types.Bool   `tfsdk:"path_enabled"`
	Subdomain      types.String `tfsdk:"subdomain"`
	V1Enabled      types.Bool   `tfsdk:"v1_enabled"`
}

func (m *ociAttributesModel) MapFromApi(api *common.OciAttributes) {
	m.ForceBasicAuth = types.BoolValue(api.ForceBasicAuth)
	m.HttpPort = types.Int32PointerValue(api.HttpPort)
	m.HttpsPort = types.Int32PointerValue(api.HttpsPort)
	m.PathEnabled = types.BoolPointerValue(api.PathEnabled)
	m.Subdomain = types.StringPointerValue(api.Subdomain)
	m.V1Enabled = types.BoolValue(api.V1Enabled)
}

func (m *ociAttributesModel) MapToApi(api *common.OciAttributes) {
	api.ForceBasicAuth = m.ForceBasicAuth.ValueBool()
	api.HttpPort = m.HttpPort.ValueInt32Pointer()
	api.HttpsPort = m.HttpsPort.ValueInt32Pointer()
	api.PathEnabled = m.PathEnabled.ValueBoolPointer()
	api.Subdomain = m.Subdomain.ValueStringPointer()
	api.V1Enabled = m.V1Enabled.ValueBool()
}

// ociProxyAttributesModel
// ----------------------------------------
type ociProxyAttributesModel struct {
	CacheForeignLayers       types.Bool     `tfsdk:"cache_foreign_layers"`
	ForeignLayerUrlWhitelist []types.String `tfsdk:"foreign_layer_url_whitelist"`
	IndexType                types.String   `tfsdk:"index_type"`
	IndexUrl                 types.String   `tfsdk:"index_url"`
}

func (m *ociProxyAttributesModel) MapFromApi(api *common.OciProxyAttributes) {
	m.CacheForeignLayers = types.BoolPointerValue(api.CacheForeignLayers)
	m.ForeignLayerUrlWhitelist = make([]types.String, 0)
	for _, l := range api.ForeignLayerUrlWhitelist {
		m.ForeignLayerUrlWhitelist = append(m.ForeignLayerUrlWhitelist, types.StringValue(l))
	}
	m.IndexType = types.StringPointerValue(api.IndexType)
	m.IndexUrl = types.StringPointerValue(api.IndexUrl)
}

func (m *ociProxyAttributesModel) MapToApi(api *common.OciProxyAttributes) {
	api.CacheForeignLayers = m.CacheForeignLayers.ValueBoolPointer()
	api.ForeignLayerUrlWhitelist = make([]string, 0)
	for _, l := range m.ForeignLayerUrlWhitelist {
		api.ForeignLayerUrlWhitelist = append(api.ForeignLayerUrlWhitelist, l.ValueString())
	}
	api.IndexType = m.IndexType.ValueStringPointer()
	api.IndexUrl = m.IndexUrl.ValueStringPointer()
}

// ociCosignConfigurationModel
// ----------------------------------------
type ociCosignConfigurationModel struct {
	Enforcement   types.String `tfsdk:"enforcement"`
	IdentityRegex types.String `tfsdk:"identity_regex"`
	IssuerRegex   types.String `tfsdk:"issuer_regex"`
}

// MapFromApi intentionally leaves IdentityRegex/IssuerRegex untouched when the API returns them
// null instead of forcing them null, mirroring the RepositoryHttpClientAuthenticationModel
// idiom for write-only/not-always-echoed fields. Confirmed against a live NXRM 3.95.0-07 server:
// with `enforcement` NONE, the API silently drops these two fields from its response even when
// the request set them, so overwriting unconditionally would erase a value that survived from a
// prior Read into the current state -- across Read (refresh), m is the existing state, so simply
// not assigning here is enough for the last-known value to stick; Create/Update instead fill any
// gap from the plan via MapMissingApiFieldsFromPlan (see UpdateStateFromPlanForNonApiFields).
func (m *ociCosignConfigurationModel) MapFromApi(api *common.OciCosignConfiguration) {
	m.Enforcement = types.StringValue(api.Enforcement)
	if api.IdentityRegex != nil {
		m.IdentityRegex = types.StringValue(*api.IdentityRegex)
	}
	if api.IssuerRegex != nil {
		m.IssuerRegex = types.StringValue(*api.IssuerRegex)
	}
}

func (m *ociCosignConfigurationModel) MapToApi(api *common.OciCosignConfiguration) {
	api.Enforcement = m.Enforcement.ValueString()
	api.IdentityRegex = m.IdentityRegex.ValueStringPointer()
	api.IssuerRegex = m.IssuerRegex.ValueStringPointer()
}

// MapMissingApiFieldsFromPlan mirrors identity_regex/issuer_regex from the plan into state
// whenever the API came back with them null. Confirmed against a live NXRM 3.95.0-07 server:
// when `enforcement` is NONE, the API silently drops identity_regex/issuer_regex from its
// response even if the request set them (matching the field docs: "ignored when enforcement is
// NONE") -- a permanent server behavior, not a client gap, so this mirrors the plan-mirroring
// workaround already used elsewhere in this codebase (see RepositoryHttpClientAuthenticationModel)
// rather than trusting the API's own (silently nulling) response.
func (m *ociCosignConfigurationModel) MapMissingApiFieldsFromPlan(planModel *ociCosignConfigurationModel) {
	if planModel == nil {
		return
	}
	if m.IdentityRegex.IsNull() {
		m.IdentityRegex = planModel.IdentityRegex
	}
	if m.IssuerRegex.IsNull() {
		m.IssuerRegex = planModel.IssuerRegex
	}
}

// newDefaultOciCosignConfigurationModel mirrors the default NXRM itself applies when a
// create/update request omits `cosign` entirely (confirmed against a live NXRM 3.95.0-07
// server, which always echoes this shape back on GET) -- used as the FromApiModel fallback so
// state always matches the schema's own Computed default (see format/oci.go).
func newDefaultOciCosignConfigurationModel() *ociCosignConfigurationModel {
	return &ociCosignConfigurationModel{
		Enforcement:   types.StringValue(common.OCI_COSIGN_ENFORCEMENT_NONE),
		IdentityRegex: types.StringNull(),
		IssuerRegex:   types.StringNull(),
	}
}

// mapOciCosignFromApi reuses existing (the model's current Cosign, which on a Read/refresh is
// the prior state -- possibly still carrying an identity_regex/issuer_regex value the API just
// silently omitted) rather than always allocating a fresh struct, so MapFromApi's deliberate
// no-op-when-nil behavior actually has a prior value to leave alone. Falls back to the same
// NXRM-observed default as an explicit Cosign-less API response when api is nil.
func mapOciCosignFromApi(existing *ociCosignConfigurationModel, api *common.OciCosignConfiguration) *ociCosignConfigurationModel {
	if api == nil {
		return newDefaultOciCosignConfigurationModel()
	}
	if existing == nil {
		existing = &ociCosignConfigurationModel{
			IdentityRegex: types.StringNull(),
			IssuerRegex:   types.StringNull(),
		}
	}
	existing.MapFromApi(api)
	return existing
}

// OCI Hosted
// ----------------------------------------
type RepositoryOciHostedModel struct {
	BasicRepositoryModel
	Storage   dockerHostedStorageModel     `tfsdk:"storage"`
	Component *RepositoryComponentModel    `tfsdk:"component"`
	Oci       ociAttributesModel           `tfsdk:"oci"`
	Cosign    *ociCosignConfigurationModel `tfsdk:"cosign"`
}

func (m *RepositoryOciHostedModel) MapMissingApiFieldsFromPlan(planModel RepositoryOciHostedModel) {
	if m.Cosign != nil {
		m.Cosign.MapMissingApiFieldsFromPlan(planModel.Cosign)
	}
}

func (m *RepositoryOciHostedModel) FromApiModel(api common.OciHostedApiRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = NewRepositoryCleanupModel()
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
	} else {
		m.Cleanup = nil
	}

	m.Storage.MapFromApi(&api.Storage)

	if api.Component != nil {
		m.Component = &RepositoryComponentModel{}
		m.Component.MapFromApi(api.Component)
	} else {
		m.Component = &RepositoryComponentModel{
			ProprietaryComponents: types.BoolValue(false),
		}
	}

	// OCI Specific
	m.Oci.MapFromApi(&api.Oci)

	m.Cosign = mapOciCosignFromApi(m.Cosign, api.Cosign)
}

func (m *RepositoryOciHostedModel) ToApiCreateModel() common.OciHostedRepositoryApiRequest {
	apiModel := common.OciHostedRepositoryApiRequest{
		Name:   m.Name.ValueString(),
		Online: m.Online.ValueBool(),
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

	// OCI Specific
	m.Oci.MapToApi(&apiModel.Oci)

	if m.Cosign != nil {
		apiModel.Cosign = &common.OciCosignConfiguration{}
		m.Cosign.MapToApi(apiModel.Cosign)
	}

	return apiModel
}

func (m *RepositoryOciHostedModel) ToApiUpdateModel() common.OciHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// OCI Proxy
// ----------------------------------------
type RepositoryOciProxyModel struct {
	RepositoryProxyModel
	Oci                        ociAttributesModel               `tfsdk:"oci"`
	OciProxy                   ociProxyAttributesModel          `tfsdk:"oci_proxy"`
	Cosign                     *ociCosignConfigurationModel     `tfsdk:"cosign"`
	FirewallAuditAndQuarantine *FirewallAuditAndQuarantineModel `tfsdk:"repository_firewall"`
}

func (m *RepositoryOciProxyModel) MapMissingApiFieldsFromPlan(planModel RepositoryOciProxyModel) {
	m.HttpClient.MapMissingApiFieldsFromPlan(planModel.HttpClient)
	if m.Cosign != nil {
		m.Cosign.MapMissingApiFieldsFromPlan(planModel.Cosign)
	}
}

func (m *RepositoryOciProxyModel) FromApiModel(api common.OciProxyApiRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = NewRepositoryCleanupModel()
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
			PreemptivePullEnabled: types.BoolValue(common.DEFAULT_PROXY_PREEMPTIVE_PULL),
			AssetPathRegex:        types.StringNull(),
		}
	}

	// OCI Specific
	m.Oci.MapFromApi(&api.Oci)
	if api.OciProxy != nil {
		m.OciProxy.MapFromApi(api.OciProxy)
	}

	m.Cosign = mapOciCosignFromApi(m.Cosign, api.Cosign)
}

func (m *RepositoryOciProxyModel) ToApiCreateModel() common.OciProxyRepositoryApiRequest {
	apiModel := common.OciProxyRepositoryApiRequest{
		Name:   m.Name.ValueString(),
		Online: m.Online.ValueBool(),
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

	// OCI Specific
	m.Oci.MapToApi(&apiModel.Oci)
	apiModel.OciProxy = &common.OciProxyAttributes{}
	m.OciProxy.MapToApi(apiModel.OciProxy)

	if m.Cosign != nil {
		apiModel.Cosign = &common.OciCosignConfiguration{}
		m.Cosign.MapToApi(apiModel.Cosign)
	}

	return apiModel
}

func (m *RepositoryOciProxyModel) ToApiUpdateModel() common.OciProxyRepositoryApiRequest {
	return m.ToApiCreateModel()
}

// OCI Group
// ----------------------------------------
type RepositoryOciGroupModel struct {
	RepositoryGroupDeployModel
	Oci    ociAttributesModel           `tfsdk:"oci"`
	Cosign *ociCosignConfigurationModel `tfsdk:"cosign"`
}

func (m *RepositoryOciGroupModel) MapMissingApiFieldsFromPlan(planModel RepositoryOciGroupModel) {
	if m.Cosign != nil {
		m.Cosign.MapMissingApiFieldsFromPlan(planModel.Cosign)
	}
}

func (m *RepositoryOciGroupModel) FromApiModel(api common.OciGroupApiRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Storage
	m.Storage.MapFromApi(&api.Storage)

	// Group Attributes
	m.Group.MapFromApi(&api.Group)

	// OCI Specific
	m.Oci.MapFromApi(&api.Oci)

	m.Cosign = mapOciCosignFromApi(m.Cosign, api.Cosign)
}

func (m *RepositoryOciGroupModel) ToApiCreateModel() common.OciGroupRepositoryApiRequest {
	apiModel := common.OciGroupRepositoryApiRequest{
		Name:   m.Name.ValueString(),
		Online: m.Online.ValueBool(),
	}
	m.Storage.MapToApi(&apiModel.Storage)
	m.Group.MapToApi(&apiModel.Group)

	// OCI Specific
	m.Oci.MapToApi(&apiModel.Oci)

	if m.Cosign != nil {
		apiModel.Cosign = &common.OciCosignConfiguration{}
		m.Cosign.MapToApi(apiModel.Cosign)
	}

	return apiModel
}

func (m *RepositoryOciGroupModel) ToApiUpdateModel() common.OciGroupRepositoryApiRequest {
	return m.ToApiCreateModel()
}
