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

// BasePrivilegeModel
// ------------------------------------
type BasePrivilegeModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	Type        types.String `tfsdk:"type"`
}

func (m *BasePrivilegeModel) MapFromApi(api *sonatyperepo.ApiPrivilegeRequest) {
	m.Name = types.StringValue(api.Name)
	m.Description = types.StringPointerValue(api.Description)
	m.ReadOnly = types.BoolPointerValue(api.ReadOnly)
	m.Type = types.StringValue(api.Type)
}

// PrivilegesModel (used by DataSource)
// ------------------------------------
type PrivilegesModel struct {
	Privileges []BasePrivilegeModel `tfsdk:"privileges"`
}

// PrivilegeModelResource
// ------------------------------------
type PrivilegeModelResource struct {
	BasePrivilegeModel
	LastUpdated types.String `tfsdk:"last_updated"`
}
