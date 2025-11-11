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

// Properties for Core: Base URL
// ----------------------------------------
type CapabilityPropertiesUiBranding struct {
	FooterEnabled types.Bool   `tfsdk:"footer_enabled" nxrm:"footerEnabled"`
	FooterHtml    types.String `tfsdk:"footer_html" nxrm:"footerHtml"`
	HeaderEnabled types.Bool   `tfsdk:"header_enabled" nxrm:"headerEnabled"`
	HeaderHtml    types.String `tfsdk:"header_html" nxrm:"headerHtml"`
}

func (p *CapabilityPropertiesUiBranding) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: UI: Branding
// ----------------------------------------
type UiBrandingCapabilityModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesUiBranding `tfsdk:"properties"`
}

func (m *UiBrandingCapabilityModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesUiBranding{}
	m.Properties.FooterEnabled = types.BoolValue(ParseBool(
		(*api.Properties)["footerEnabled"],
		common.CAPABILITY_UI_BRANDING_DEFAULT_FOOTER_ENABLED,
	))
	m.Properties.FooterHtml = types.StringValue((*api.Properties)["footerHtml"])
	m.Properties.HeaderEnabled = types.BoolValue(ParseBool(
		(*api.Properties)["headerEnabled"],
		common.CAPABILITY_UI_BRANDING_DEFAULT_HEADER_ENABLED,
	))
	m.Properties.HeaderHtml = types.StringValue((*api.Properties)["headerHtml"])
}

func (m *UiBrandingCapabilityModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_UI_BRANDING.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *UiBrandingCapabilityModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

// Properties for UI: Settings
// ----------------------------------------
type CapabilityPropertiesUiSettings struct {
	DebugAllowed                types.Bool   `tfsdk:"debug_allowed" nxrm:"debugAllowed"`
	LongRequestTimeout          types.Int32  `tfsdk:"long_request_timeout" nxrm:"longRequestTimeout"`
	RequestTimeout              types.Int32  `tfsdk:"request_timeout" nxrm:"requestTimeout"`
	SessionTimeout              types.Int32  `tfsdk:"session_timeout" nxrm:"sessionTimeout"`
	StatusIntervalAnonymous     types.Int32  `tfsdk:"status_interval_anonymous" nxrm:"statusIntervalAnonymous"`
	StatusIntervalAuthenticated types.Int32  `tfsdk:"status_interval_authenticated" nxrm:"statusIntervalAuthenticated"`
	Title                       types.String `tfsdk:"title" nxrm:"title"`
	// SearchRequestTimeout        types.Int32  `tfsdk:"search_request_timeout" nxrm:"searchRequestTimeout"`
}

func (p *CapabilityPropertiesUiSettings) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Capability: UI: Settings
// ----------------------------------------
type UiSettingsCapabilityModel struct {
	BaseCapabilityModel
	Properties *CapabilityPropertiesUiSettings `tfsdk:"properties"`
}

func (m *UiSettingsCapabilityModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringValue(*api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolValue(*api.Enabled)
	m.Properties = &CapabilityPropertiesUiSettings{}
	m.Properties.DebugAllowed = types.BoolValue(ParseBool(
		(*api.Properties)["debugAllowed"],
		common.CAPABILITY_UI_SETTINGS_DEFAULT_DEBUG_ALLOWED,
	))
	m.Properties.LongRequestTimeout = types.Int32Value(ParseInt32(
		(*api.Properties)["longRequestTimeout"],
		common.CAPABILITY_UI_SETTINGS_DEFAULT_LONG_REQUEST_TIMEOUT,
	))
	m.Properties.RequestTimeout = types.Int32Value(ParseInt32(
		(*api.Properties)["requestTimeout"],
		common.CAPABILITY_UI_SETTINGS_DEFAULT_REQUEST_TIMEOUT,
	))
	// m.Properties.SearchRequestTimeout = types.Int32Value(ParseInt32(
	// 	(*api.Properties)["searchRequestTimeout"],
	// 	common.CAPABILITY_UI_SETTINGS_DEFAULT_SEARCH_REQUEST_TIMEOUT,
	// ))
	m.Properties.SessionTimeout = types.Int32Value(ParseInt32(
		(*api.Properties)["sessionTimeout"],
		common.CAPABILITY_UI_SETTINGS_DEFAULT_SESSION_TIMEOUT,
	))
	m.Properties.StatusIntervalAnonymous = types.Int32Value(ParseInt32(
		(*api.Properties)["statusIntervalAnonymous"],
		common.CAPABILITY_UI_SETTINGS_DEFAULT_STATUS_INTERVAL_ANONYMOUS,
	))
	m.Properties.StatusIntervalAuthenticated = types.Int32Value(ParseInt32(
		(*api.Properties)["statusIntervalAuthenticated"],
		common.CAPABILITY_UI_SETTINGS_DEFAULT_STATUS_INTERVAL_AUTHENTICATED,
	))
	m.Properties.Title = types.StringValue((*api.Properties)["title"])
}

func (m *UiSettingsCapabilityModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_UI_SETTINGS.StringPointer()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *UiSettingsCapabilityModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
