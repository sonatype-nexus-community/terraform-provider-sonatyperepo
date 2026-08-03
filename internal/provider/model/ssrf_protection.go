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

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// SsrfProtectionModel represents the Terraform model for the SSRF protection configuration
type SsrfProtectionModel struct {
	Enabled        types.Bool     `tfsdk:"enabled"`
	AllowedDomains []types.String `tfsdk:"allowed_domains"`
	AllowedIPs     []types.String `tfsdk:"allowed_ips"`
	LastUpdated    types.String   `tfsdk:"last_updated"`
}

// MapFromApi maps API response to model
func (m *SsrfProtectionModel) MapFromApi(api *sonatyperepo.SsrfProtectionConfigurationXO) {
	m.Enabled = types.BoolValue(api.Enabled)

	m.AllowedDomains = make([]types.String, 0)
	for _, d := range api.AllowedDomains {
		m.AllowedDomains = append(m.AllowedDomains, types.StringValue(d))
	}

	m.AllowedIPs = make([]types.String, 0)
	for _, ip := range api.AllowedIPs {
		m.AllowedIPs = append(m.AllowedIPs, types.StringValue(ip))
	}
}

// MapToApi maps model to API request payload
func (m *SsrfProtectionModel) MapToApi(api *sonatyperepo.SsrfProtectionConfigurationXO) {
	api.Enabled = m.Enabled.ValueBool()

	api.AllowedDomains = make([]string, 0)
	for _, d := range m.AllowedDomains {
		api.AllowedDomains = append(api.AllowedDomains, d.ValueString())
	}

	api.AllowedIPs = make([]string, 0)
	for _, ip := range m.AllowedIPs {
		api.AllowedIPs = append(api.AllowedIPs, ip.ValueString())
	}
}
