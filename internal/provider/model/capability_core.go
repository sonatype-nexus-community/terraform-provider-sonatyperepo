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
	"strconv"
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Properties for Core: Base URL
// ----------------------------------------
type CapabilityPropertiesCoreBaseUrl struct {
	Url types.String `tfsdk:"url" nxrm:"url"`
}

func (p *CapabilityPropertiesCoreBaseUrl) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: Core: Base URL
// ----------------------------------------
type CapabilityCoreBaseUrlModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesCoreBaseUrl `tfsdk:"properties"`
}

func (m *CapabilityCoreBaseUrlModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesCoreBaseUrl{}
	m.Properties.Url = types.StringValue((*api.Properties)["url"])
}

func (m *CapabilityCoreBaseUrlModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_CORE_BASE_URL.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *CapabilityCoreBaseUrlModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

// Properties for Core: Outreach Management
// ----------------------------------------
type CapabilityPropertiesCoreOutreach struct {
	AlwaysRemote types.Bool   `tfsdk:"always_remote" nxrm:"alwaysRemote"`
	OverrideUrl  types.String `tfsdk:"override_url" nxrm:"overrideUrl"`
}

func (p *CapabilityPropertiesCoreOutreach) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: Core: Outreach Management
// ----------------------------------------
type CapabilityCoreOutreachModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesCoreOutreach `tfsdk:"properties"`
}

func (m *CapabilityCoreOutreachModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesCoreOutreach{}
	m.Properties.OverrideUrl = types.StringValue((*api.Properties)["overrideUrl"])
	alwaysRemote, err := strconv.ParseBool((*api.Properties)["alwaysRemote"])
	if err != nil {
		alwaysRemote = false
	}
	m.Properties.AlwaysRemote = types.BoolValue(alwaysRemote)
}

func (m *CapabilityCoreOutreachModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_OUTREACH.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *CapabilityCoreOutreachModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
