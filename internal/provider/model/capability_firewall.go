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

// Properties for Firewall Audit & Quarantine
// ----------------------------------------
type CapabilityPropertiesFirewallAuditQuarantine struct {
	Quarantine types.Bool   `tfsdk:"quarantine" nxrm:"quarantine"`
	Repository types.String `tfsdk:"repository" nxrm:"repository"`
}

func (p *CapabilityPropertiesFirewallAuditQuarantine) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: Firewall Audit & Quarantine
// ----------------------------------------
type CapabilityFirewallAuditQuarantineModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesFirewallAuditQuarantine `tfsdk:"properties" nxrm:"properties"`
}

func (m *CapabilityFirewallAuditQuarantineModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesFirewallAuditQuarantine{}
	m.Properties.Quarantine = types.BoolValue(ParseBool(
		(*api.Properties)["quarantine"],
		common.CAPABILITY_FIREWALL_AUDIT_QUARANTINE_DEFAULT_QUARANTINE,
	))
	m.Properties.Repository = types.StringValue((*api.Properties)["repository"])
}

func (m *CapabilityFirewallAuditQuarantineModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_FIREWALL_AUDIT_QUARANTINE.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *CapabilityFirewallAuditQuarantineModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
