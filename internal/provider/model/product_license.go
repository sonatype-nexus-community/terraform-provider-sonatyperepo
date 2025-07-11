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
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// ProductLicenseModelResource
// -----------------------------------
type ProductLicenseModelResource struct {
	LicenseData       types.String `tfsdk:"license_data"`
	ContactCompany    types.String `tfsdk:"contact_company"`
	ContactEmail      types.String `tfsdk:"contact_email"`
	ContactName       types.String `tfsdk:"contact_name"`
	EffectiveDate     types.String `tfsdk:"effective_date"`
	ExpirationDate    types.String `tfsdk:"expiration_date"`
	Features          types.String `tfsdk:"features"`
	Fingerprint       types.String `tfsdk:"fingerprint"`
	LicenseType       types.String `tfsdk:"license_type"`
	LicensedUsers     types.String `tfsdk:"licensed_users"`
	MaxRepoComponents types.Int64  `tfsdk:"max_repo_components"`
	MaxRepoRequests   types.Int64  `tfsdk:"max_repo_requests"`
	LastUpdated       types.String `tfsdk:"last_updated"`
}

type ProductLicenseCreateModel struct {
	LicenseData types.String `tfsdk:"license_data"`
}

func (m *ProductLicenseModelResource) MapFromApi(api *sonatyperepo.ApiLicenseDetailsXO) {
	m.ContactCompany = types.StringPointerValue(api.ContactCompany)
	m.ContactEmail = types.StringPointerValue(api.ContactEmail)
	m.ContactName = types.StringPointerValue(api.ContactName)
	m.EffectiveDate = types.StringValue(api.EffectiveDate.Format(time.RFC850))
	m.ExpirationDate = types.StringValue(api.ExpirationDate.Format(time.RFC850))
	m.Features = types.StringPointerValue(api.Features)
	m.Fingerprint = types.StringPointerValue(api.Fingerprint)
	m.LicenseType = types.StringPointerValue(api.LicenseType)
	m.LicensedUsers = types.StringPointerValue(api.LicensedUsers)
	m.MaxRepoComponents = types.Int64PointerValue(api.MaxRepoComponents)
	m.MaxRepoRequests = types.Int64PointerValue(api.MaxRepoRequests)
}

// func (m *RoleModel) MapToApi(api *sonatyperepo.RoleXORequest) {
// 	if api == nil {
// 		api = sonatyperepo.NewRoleXORequestWithDefaults()
// 		api.Privileges = make([]string, 0)
// 		api.Roles = make([]string, 0)
// 	}
// 	api.Id = m.Id.ValueStringPointer()
// 	api.Name = m.Name.ValueStringPointer()
// 	api.Description = m.Description.ValueStringPointer()
// 	for _, p := range m.Privileges {
// 		api.Privileges = append(api.Privileges, p.ValueString())
// 	}
// 	for _, r := range m.Roles {
// 		api.Roles = append(api.Roles, r.ValueString())
// 	}
// }
