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

// Properties for: RUT Auth
// ----------------------------------------
type CapabilityPropertiesRutAuth struct {
	HttpHeader types.String `tfsdk:"http_header" nxrm:"http_header"`
}

func (p *CapabilityPropertiesRutAuth) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: RUT Auth
// ----------------------------------------
type CapabilityRutAuthModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesRutAuth `tfsdk:"properties"`
}

func (m *CapabilityRutAuthModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesRutAuth{}
	m.Properties.HttpHeader = types.StringValue((*api.Properties)["http_header"])
}

func (m *CapabilityRutAuthModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_RUT_AUTH.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *CapabilityRutAuthModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
