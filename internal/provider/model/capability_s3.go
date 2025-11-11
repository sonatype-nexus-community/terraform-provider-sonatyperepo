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
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Properties for S3: Custom S3 Regions
// ----------------------------------------
type CapabilityPropertiesCustomS3Regions struct {
	Regions types.Set `tfsdk:"regions" nxrm:"regions"`
}

func (p *CapabilityPropertiesCustomS3Regions) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: S3: Custom S3 Regions
// ----------------------------------------
type CapabilityCustomS3RegionsModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesCustomS3Regions `tfsdk:"properties"`
}

func (m *CapabilityCustomS3RegionsModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesCustomS3Regions{}

	// Regions
	regionsStr := (*api.Properties)["regions"]
	regionsParts := strings.Split(regionsStr, ",")
	var regionsValues []attr.Value
	for _, part := range regionsParts {
		regionsValues = append(regionsValues, types.StringValue(strings.TrimSpace(part)))
	}
	m.Properties.Regions = types.SetValueMust(types.StringType, regionsValues)
}

func (m *CapabilityCustomS3RegionsModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_CUSTOM_S3_REGIONS.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *CapabilityCustomS3RegionsModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
