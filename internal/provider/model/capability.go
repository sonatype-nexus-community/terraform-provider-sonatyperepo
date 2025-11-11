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
	"github.com/hashicorp/terraform-plugin-framework/types"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// CapabilitiesListModel - used for data source
// ----------------------------------------
type CapabilitiesListModel struct {
	Capabilities []CapabilityModel `tfsdk:"capabilities"`
}

// CapabilitCommonModel - used for data source and resource
// ----------------------------------------
type CapabilitCommonModel struct {
	Id      types.String `tfsdk:"id"`
	Notes   types.String `tfsdk:"notes"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

// CapabilitiesListModel - used for data source
// ----------------------------------------
type CapabilityModel struct {
	CapabilitCommonModel
	Type       types.String      `tfsdk:"type"`
	Properties map[string]string `tfsdk:"properties"`
}

// Base Capability Model - used for create and update
// ----------------------------------------
type BaseCapabilityModel struct {
	CapabilitCommonModel
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (m *BaseCapabilityModel) MapFromApi(api *v3.CapabilityDTO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Notes = types.StringPointerValue(api.Notes)
	m.Enabled = types.BoolPointerValue(api.Enabled)
}

func (m *BaseCapabilityModel) toApiCreateModel() *v3.CapabilityDTO {
	api := v3.NewCapabilityDTOWithDefaults()
	api.Notes = m.Notes.ValueStringPointer()
	api.Enabled = m.Enabled.ValueBoolPointer()
	return api
}

func (m *BaseCapabilityModel) toApiUpdateModel() *v3.CapabilityDTO {
	api := v3.NewCapabilityDTOWithDefaults()
	api.Id = m.Id.ValueStringPointer()
	api.Notes = m.Notes.ValueStringPointer()
	api.Enabled = m.Enabled.ValueBoolPointer()
	return api
}
