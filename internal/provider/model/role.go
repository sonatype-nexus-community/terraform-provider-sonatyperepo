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

type RolesModel struct {
	Roles []RoleModelIncludingReadOnly `tfsdk:"roles"`
}

type RoleModel struct {
	Id          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Privileges  []types.String `tfsdk:"privileges"`
	Roles       []types.String `tfsdk:"roles"`
}

type RoleModelResource struct {
	RoleModel
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (m *RoleModel) MapFromApi(api *sonatyperepo.RoleXOResponse) {
	m.Id = types.StringPointerValue(api.Id)
	m.Name = types.StringPointerValue(api.Name)
	m.Description = types.StringPointerValue(api.Description)
	m.Privileges = make([]types.String, 0)
	for _, p := range api.Privileges {
		m.Privileges = append(m.Privileges, types.StringValue(p))
	}
	m.Roles = make([]types.String, 0)
	for _, r := range api.Roles {
		m.Roles = append(m.Roles, types.StringValue(r))
	}
}

func (m *RoleModel) MapToApi(api *sonatyperepo.RoleXORequest) {
	if api == nil {
		api = sonatyperepo.NewRoleXORequestWithDefaults()
		api.Privileges = make([]string, 0)
		api.Roles = make([]string, 0)
	}
	api.Id = m.Id.ValueStringPointer()
	api.Name = m.Name.ValueStringPointer()
	api.Description = m.Description.ValueStringPointer()
	for _, p := range m.Privileges {
		api.Privileges = append(api.Privileges, p.ValueString())
	}
	for _, r := range m.Roles {
		api.Roles = append(api.Roles, r.ValueString())
	}
}

type RoleModelIncludingReadOnly struct {
	RoleModel
	ReadOnly types.Bool   `tfsdk:"read_only"`
	Source   types.String `tfsdk:"source"`
}

func (m *RoleModelIncludingReadOnly) MapFromApi(api *sonatyperepo.RoleXOResponse) {
	m.Id = types.StringPointerValue(api.Id)
	m.Name = types.StringPointerValue(api.Name)
	m.Description = types.StringPointerValue(api.Description)
	for _, p := range api.Privileges {
		m.Privileges = append(m.Privileges, types.StringValue(p))
	}
	for _, r := range api.Roles {
		m.Roles = append(m.Roles, types.StringValue(r))
	}
	m.ReadOnly = types.BoolPointerValue(api.ReadOnly)
	m.Source = types.StringPointerValue(api.Source)
}
