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

// Properties for Webhook: Global
// ----------------------------------------
type CapabilityPropertiesWebhookGlobal struct {
	Names  types.Set    `tfsdk:"names" nxrm:"names"`
	Secret types.String `tfsdk:"secret" nxrm:"secret"`
	Url    types.String `tfsdk:"url" nxrm:"url"`
}

func (p *CapabilityPropertiesWebhookGlobal) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	result := make(map[string]string)

	// names
	if !p.Names.IsNull() && !p.Names.IsUnknown() {
		elements := p.Names.Elements()
		var strs []string
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				strs = append(strs, strVal.ValueString())
			}
		}
		result["names"] = strings.Join(strs, ",")
	}

	// secret
	if !p.Secret.IsNull() && !p.Secret.IsUnknown() {
		result["secret"] = p.Secret.ValueString()
	}

	// url
	if !p.Url.IsNull() && !p.Url.IsUnknown() {
		result["url"] = p.Url.ValueString()
	}

	return &result
}

// Capability: Webhook: Global
// ----------------------------------------
type WebhookGlobalCapabilityModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesWebhookGlobal `tfsdk:"properties"`
}

func (m *WebhookGlobalCapabilityModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Notes = types.StringPointerValue(api.Notes)
	m.Enabled = types.BoolPointerValue(api.Enabled)
	m.Properties = &CapabilityPropertiesWebhookGlobal{}
	namesStr := (*api.Properties)["names"]
	namesParts := strings.Split(namesStr, ",")
	var namesValues []attr.Value
	for _, part := range namesParts {
		namesValues = append(namesValues, types.StringValue(strings.TrimSpace(part)))
	}
	m.Properties.Names = types.SetValueMust(types.StringType, namesValues)
	if secret, ok := (*api.Properties)["secret"]; ok && secret != "" {
		m.Properties.Secret = types.StringValue(secret)
	} else {
		m.Properties.Secret = types.StringNull()
	}
	m.Properties.Url = types.StringValue((*api.Properties)["url"])
}

func (m *WebhookGlobalCapabilityModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_WEBHOOK_GLOBAL.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *WebhookGlobalCapabilityModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

// Properties for Webhook: Repository
// ----------------------------------------
type CapabilityPropertiesWebhookRepository struct {
	CapabilityPropertiesWebhookGlobal
	Repository types.String `tfsdk:"repository" nxrm:"repository"`
}

func (p *CapabilityPropertiesWebhookRepository) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	result := make(map[string]string)

	// names
	if !p.Names.IsNull() && !p.Names.IsUnknown() {
		elements := p.Names.Elements()
		var strs []string
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				strs = append(strs, strVal.ValueString())
			}
		}
		result["names"] = strings.Join(strs, ",")
	}

	// secret
	if !p.Secret.IsNull() && !p.Secret.IsUnknown() {
		result["secret"] = p.Secret.ValueString()
	}

	// url
	if !p.Url.IsNull() && !p.Url.IsUnknown() {
		result["url"] = p.Url.ValueString()
	}

	// repository
	if !p.Repository.IsNull() && !p.Repository.IsUnknown() {
		result["repository"] = p.Repository.ValueString()
	}

	return &result
}

// Capability: Webhook: Repository
// ----------------------------------------
type WebhookRepositoryCapabilityModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesWebhookRepository `tfsdk:"properties"`
}

func (m *WebhookRepositoryCapabilityModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesWebhookRepository{}
	namesStr := (*api.Properties)["names"]
	namesParts := strings.Split(namesStr, ",")
	var namesValues []attr.Value
	for _, part := range namesParts {
		namesValues = append(namesValues, types.StringValue(strings.TrimSpace(part)))
	}
	m.Properties.Names = types.SetValueMust(types.StringType, namesValues)
	if secret, ok := (*api.Properties)["secret"]; ok && secret != "" {
		m.Properties.Secret = types.StringValue(secret)
	} else {
		m.Properties.Secret = types.StringNull()
	}
	m.Properties.Url = types.StringValue((*api.Properties)["url"])
	m.Properties.Repository = types.StringValue((*api.Properties)["repository"])
}

func (m *WebhookRepositoryCapabilityModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_WEBHOOK_REPOSITORY.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *WebhookRepositoryCapabilityModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
