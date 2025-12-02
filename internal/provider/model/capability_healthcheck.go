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

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Properties for Healthcheck
// ----------------------------------------
type CapabilityPropertiesHealthcheck struct {
	ConfiguredForAllProxies types.Bool `tfsdk:"configured_for_all_proxies" nxrm:"configuredForAll"`
	UseNexusTruststore      types.Bool `tfsdk:"use_nexus_truststore" nxrm:"useTrustStore"`
}

func (p *CapabilityPropertiesHealthcheck) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: Healthcheck
// ----------------------------------------
type CapabilityHealthcheckModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesHealthcheck `tfsdk:"properties" nxrm:"properties"`
}

func (m *CapabilityHealthcheckModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Notes = types.StringPointerValue(api.Notes)
	m.Enabled = types.BoolPointerValue(api.Enabled)
	m.Properties = &CapabilityPropertiesHealthcheck{}
	m.Properties.ConfiguredForAllProxies = types.BoolValue(ParseBool(
		(*api.Properties)["configuredForAll"],
		common.CAPABILITY_HEALTHCHECK_DEFAULT_CONFIGURED_FOR_ALL,
	))
	m.Properties.UseNexusTruststore = types.BoolValue(ParseBool(
		(*api.Properties)["useTrustStore"],
		common.CAPABILITY_HEALTHCHECK_DEFAULT_USE_NEXUS_TRUSTSTORE,
	))
}

func (m *CapabilityHealthcheckModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_HEALTHCHECK.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *CapabilityHealthcheckModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
