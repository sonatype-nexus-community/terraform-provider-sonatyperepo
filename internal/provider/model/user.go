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

// BaseUserModel
// ------------------------------------
type BaseUserModel struct {
	UserId       types.String   `tfsdk:"user_id"`
	FirstName    types.String   `tfsdk:"first_name"`
	LastName     types.String   `tfsdk:"last_name"`
	EmailAddress types.String   `tfsdk:"email_address"`
	Status       types.String   `tfsdk:"status"`
	Roles        []types.String `tfsdk:"roles"`
	ReadOnly     types.Bool     `tfsdk:"read_only"`
	Source       types.String   `tfsdk:"source"`
}

// BaseUserModel
// ------------------------------------
type UserModelResource struct {
	BaseUserModel
	Password    types.String `tfsdk:"password"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (m *UserModelResource) MapFromApi(api *sonatyperepo.ApiUser) {
	m.UserId = types.StringPointerValue(api.UserId)
	m.FirstName = types.StringPointerValue(api.FirstName)
	m.LastName = types.StringPointerValue(api.LastName)
	m.EmailAddress = types.StringPointerValue(api.EmailAddress)
	m.Status = types.StringValue(api.Status)
	m.Roles = make([]types.String, 0)
	for _, r := range api.GetRoles() {
		m.Roles = append(m.Roles, types.StringValue(r))
	}
	// Ignore Password
	m.Source = types.StringPointerValue(api.Source)
	m.ReadOnly = types.BoolPointerValue(api.ReadOnly)
}

func (m *UserModelResource) MapToApi(api *sonatyperepo.ApiUser) {
	api.UserId = m.UserId.ValueStringPointer()
	api.FirstName = m.FirstName.ValueStringPointer()
	api.LastName = m.LastName.ValueStringPointer()
	api.EmailAddress = m.EmailAddress.ValueStringPointer()
	api.Status = m.Status.ValueString()
	api.Roles = make([]string, 0)
	for _, r := range m.Roles {
		api.Roles = append(api.Roles, r.ValueString())
	}
	api.ReadOnly = m.ReadOnly.ValueBoolPointer()
	// Source should be set by the caller for updates since it's a computed field
}

func (m *UserModelResource) MapToCreateApi(api *sonatyperepo.ApiCreateUser) {
	api.UserId = m.UserId.ValueStringPointer()
	api.FirstName = m.FirstName.ValueStringPointer()
	api.LastName = m.LastName.ValueStringPointer()
	api.EmailAddress = m.EmailAddress.ValueStringPointer()
	api.Status = m.Status.ValueString()
	api.Roles = make([]string, 0)
	for _, r := range m.Roles {
		api.Roles = append(api.Roles, r.ValueString())
	}
	api.Password = m.Password.ValueStringPointer()
}

// UserModel (used by DataSource)
// ------------------------------------
type UserModel struct {
	BaseUserModel
	ExternalRoles []types.String `tfsdk:"external_roles"`
}

func (m *UserModel) MapFromApi(api *sonatyperepo.ApiUser) {
	m.UserId = types.StringPointerValue(api.UserId)
	m.FirstName = types.StringPointerValue(api.FirstName)
	m.LastName = types.StringPointerValue(api.LastName)
	m.EmailAddress = types.StringPointerValue(api.EmailAddress)
	m.Status = types.StringValue(api.Status)
	m.Roles = make([]types.String, 0)
	for _, r := range api.GetRoles() {
		m.Roles = append(m.Roles, types.StringValue(r))
	}
	m.ExternalRoles = make([]types.String, 0)
	for _, r := range api.GetExternalRoles() {
		m.ExternalRoles = append(m.ExternalRoles, types.StringValue(r))
	}
	// Ignore Password
	m.Source = types.StringPointerValue(api.Source)
	m.ReadOnly = types.BoolPointerValue(api.ReadOnly)
}

// UserModels (used by DataSource)
// ------------------------------------
type UsersModel struct {
	Users []UserModel `tfsdk:"users"`
}
